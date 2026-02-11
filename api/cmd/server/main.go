package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"os/signal"
	"strconv"
	"syscall"
	"time"

	"article-manager/internal/infrastructure/ai"
	"article-manager/internal/infrastructure/database"
	"article-manager/internal/infrastructure/external"
	applogger "article-manager/internal/infrastructure/logger"
	"article-manager/internal/infrastructure/repository"
	infraservice "article-manager/internal/infrastructure/service"
	"article-manager/internal/interface/handler"
	"article-manager/internal/usecase"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	if err := applogger.InitLogger(env); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer applogger.Sync()

	logger := log.New(os.Stdout, "[article-manager] ", log.LstdFlags|log.Lshortfile)

	// データベース接続
	config := loadConfig()
	db, err := database.NewMySQL(config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
	if err != nil {
		logger.Fatalf("データベース接続に失敗: %v", err)
	}
	defer func() {
		if err := database.CloseMySQLX(db); err != nil {
			logger.Printf("データベースクローズに失敗; %v", err)
		}
	}()

	// ヘルスチェック
	if err := database.HealthCheckX(db); err != nil {
		logger.Fatalf("データベースヘルスチェックに失敗: %v", err)
	}
	logger.Println("データベース接続成功")

	// マイグレーション実行
	migrationManager := database.NewMigrationManager(db, logger)
	if err := migrationManager.RunMigrations(); err != nil {
		logger.Fatalf("マイグレーション実行に失敗: %v", err)
	}

	// 依存性注入(article)
	articleRepo := repository.NewMySQLArticleRepository(db)
	articleUsecase := usecase.NewArticleUsecase(articleRepo)
	articleHandler := handler.NewArticleHandler(articleUsecase)

	// 依存性注入(tag)
	tagRepo := repository.NewMySQLTagRepository(db)
	tagUsecase := usecase.NewTagUsecase(tagRepo)
	tagHandler := handler.NewTagHandler(tagUsecase)

	// 依存性注入(ai generator)
	geminiConfig := ai.DefaultGeminiConfig(config.GeminiAPIKey)
	geminiClient := ai.NewGeminiClient(geminiConfig)
	articleGeneratorUsecase := usecase.NewArticleGeneratorUsecase(geminiClient, articleRepo, tagRepo)
	articleGeneratorHandler := handler.NewArticleGeneratorHandler(articleGeneratorUsecase)

	// 依存性注入(book recommendation)
	googleBooksConfig := external.DefaultGoogleBooksConfig(config.GoogleBooksAPIKey)
	googleBooksClient := external.NewGoogleBooksClient(googleBooksConfig)
	bookRecommendationService := infraservice.NewBookRecommendationService(geminiClient, googleBooksClient)
	bookRecommendationRepo := repository.NewMySQLBookRecommendationRepository(db)
	bookRecommendationUsecase := usecase.NewBookRecommendationUsecase(articleRepo, bookRecommendationRepo, bookRecommendationService)
	bookRecommendationHandler := handler.NewBookRecommendationHandler(bookRecommendationUsecase)

	// HTTPルーターの設定
	mux := http.NewServeMux()

	// ルートエンドポイント
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message":"Article Manager API","version":"1.0.0"}`)
	})

	// 記事一覧取得
	mux.HandleFunc("GET /api/articles", articleHandler.GetAllArticles)

	// 記事作成
	mux.HandleFunc("POST /api/articles", articleHandler.CreateArticle)

	// 記事自動作成
	mux.HandleFunc("POST /api/articles/generate", articleGeneratorHandler.GenerateArticle)

	// 記事検索（{id}より先に定義）
	mux.HandleFunc("GET /api/articles/search", articleHandler.SearchArticles)

	// 記事詳細取得
	mux.HandleFunc("GET /api/articles/{id}", extractArticleID(articleHandler.GetArticleByID))

	// 記事更新
	mux.HandleFunc("PUT /api/articles/{id}", extractArticleID(articleHandler.UpdateArticle))

	// 記事削除
	mux.HandleFunc("DELETE /api/articles/{id}", extractArticleID(articleHandler.DeleteArticle))

	// タグ一覧取得
	mux.HandleFunc("GET /api/tags", tagHandler.GetAllTags)

	// タグ作成
	mux.HandleFunc("POST /api/tags", tagHandler.CreateTag)

	// タグ詳細取得
	mux.HandleFunc("GET /api/tags/{id}", extractTagID(tagHandler.GetTagByID))

	// タグ更新
	mux.HandleFunc("PUT /api/tags/{id}", extractTagID(tagHandler.UpdateTag))

	// タグ削除
	mux.HandleFunc("DELETE /api/tags/{id}", extractTagID(tagHandler.DeleteTag))

	// 書籍推薦取得
	mux.HandleFunc("GET /api/book-recommendations", bookRecommendationHandler.GetBookRecommendations)

	// CORSミドルウェアの設定
	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// ログミドルウェア
	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			logger.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)
			next.ServeHTTP(w, r)
			logger.Printf("[%s] %s 処理時間: %v", r.Method, r.URL.Path, time.Since(start))
		})
	}

	// ミドルウェア適用
	handler := corsMiddleware(loggingMiddleware(mux))

	// HTTPサーバー設定
	server := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// サーバーを別ゴルーチンで起動
	serverErrors := make(chan error, 1)
	go func() {
		logger.Printf("HTTPサーバーを起動します: http://localhost:%s", config.Port)
		serverErrors <- server.ListenAndServe()
	}()

	// グレースフルシャットダウンのためのシグナルハンドリング
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// サーバー停止を待機
	select {
	case err := <-serverErrors:
		logger.Fatalf("サーバー起動に失敗: %v", err)
	case sig := <-shutdown:
		logger.Printf("シャットダウンシグナルを受信しました: %v", sig)

		// タイムアウト付きコンテキストを作成
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// グレースフルシャットダウンを実行
		if err := server.Shutdown(ctx); err != nil {
			logger.Printf("グレースフルシャットダウンに失敗、強制終了します: %v", err)
			if err := server.Close(); err != nil {
				logger.Printf("サーバークローズに失敗: %v", err)
			}
		}
		logger.Println("サーバーを正常にシャットダウンしました")
	}
}

type Config struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	Port              string
	GeminiAPIKey      string
	GoogleBooksAPIKey string
}

func loadConfig() Config {
	config := Config{
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "3306"),
		DBUser:            getEnv("DB_USER", ""),
		DBPassword:        getEnv("DB_PASSWORD", ""),
		DBName:            getEnv("DB_NAME", ""),
		Port:              getEnv("PORT", "8080"),
		GeminiAPIKey:      getEnv("GEMINI_API_KEY", ""),
		GoogleBooksAPIKey: getEnv("GOOGLE_BOOKS_API_KEY", ""),
	}

	// ユーザー名が設定されていない場合はエラー
	if config.DBUser == "" {
		log.Fatal("DBUser environment variable is required")
	}

	// パスワードが設定されていない場合はエラー
	if config.DBPassword == "" {
		log.Fatal("DBPassword environment variable is required")
	}

	// データベース名が設定されていない場合はエラー
	if config.DBName == "" {
		log.Fatal("DBName environment variable is required")
	}

	if config.GeminiAPIKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required")
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func extractArticleID(next func(http.ResponseWriter, *http.Request, int64)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid article ID", http.StatusBadRequest)
			return
		}
		next(w, r, id)
	}
}

func extractTagID(next func(http.ResponseWriter, *http.Request, int64)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid tag ID", http.StatusBadRequest)
			return
		}
		next(w, r, id)
	}
}
