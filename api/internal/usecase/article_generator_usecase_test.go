package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"article-manager/internal/domain/entity"
	"article-manager/internal/domain/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// モックAIGeneratorサービス
type mockAIGeneratorService struct {
	generateFunc func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error)
}

func (m *mockAIGeneratorService) GenerateArticleFromURL(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
	return m.generateFunc(ctx, req)
}

// GenerateArticleFromURLのテスト
func TestGenerateArticleFromURL(t *testing.T) {
	tests := []struct {
		name             string
		url              string
		memo             string
		setupMocks       func() (*mockAIGeneratorService, *mockArticleRepository, *mockTagRepository)
		expectedError    bool
		expectedErrorMsg string
		validateResult   func(t *testing.T, result *entity.Article)
	}{
		{
			name: "正常系：URLから記事を生成して保存成功（既存タグあり）",
			url:  "https://example.com/article",
			memo: "テストメモ",
			setupMocks: func() (*mockAIGeneratorService, *mockArticleRepository, *mockTagRepository) {
				aiService := &mockAIGeneratorService{
					generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
						return &service.GeneratedArticle{
							Title:         "AIが生成した記事タイトル",
							Summary:       "AIが生成した記事の要約です。",
							SuggestedTags: []string{"Go", "AI"},
							SourceURL:     "https://example.com/article",
							TokenUsed:     100,
							GeneratedAt:   time.Now(),
						}, nil
					},
				}
				tagRepo := &mockTagRepository{
					findByNameFunc: func(ctx context.Context, name string) (*entity.Tag, error) {
						switch name {
						case "Go":
							return &entity.Tag{ID: 1, Name: "Go"}, nil
						case "AI":
							return &entity.Tag{ID: 2, Name: "AI"}, nil
						default:
							return nil, errors.New("tag not found")
						}
					},
				}
				articleRepo := &mockArticleRepository{
					createFunc: func(ctx context.Context, article *entity.Article) (*entity.Article, error) {
						article.ID = 1
						return article, nil
					},
				}
				return aiService, articleRepo, tagRepo
			},
			expectedError: false,
			validateResult: func(t *testing.T, result *entity.Article) {
				assert.Equal(t, int64(1), result.ID)
				assert.Equal(t, "AIが生成した記事タイトル", result.Title)
				assert.Equal(t, "https://example.com/article", result.URL)
				assert.Equal(t, "AIが生成した記事の要約です。", result.Summary)
				assert.ElementsMatch(t, []string{"Go", "AI"}, result.Tags)
				assert.Equal(t, "テストメモ", result.Memo)
			},
		},
		{
			name: "正常系：URLから記事を生成して保存成功（新規タグ作成）",
			url:  "https://example.com/article",
			memo: "",
			setupMocks: func() (*mockAIGeneratorService, *mockArticleRepository, *mockTagRepository) {
				aiService := &mockAIGeneratorService{
					generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
						return &service.GeneratedArticle{
							Title:         "新しい技術記事",
							Summary:       "新しい技術の解説記事です。",
							SuggestedTags: []string{"NewTech", "Tutorial"},
							SourceURL:     "https://example.com/article",
							TokenUsed:     150,
							GeneratedAt:   time.Now(),
						}, nil
					},
				}
				tagRepo := &mockTagRepository{
					findByNameFunc: func(ctx context.Context, name string) (*entity.Tag, error) {
						return nil, errors.New("tag not found")
					},
					createFunc: func(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
						tag.ID = int64(len(tag.Name))
						return tag, nil
					},
				}
				articleRepo := &mockArticleRepository{
					createFunc: func(ctx context.Context, article *entity.Article) (*entity.Article, error) {
						article.ID = 2
						return article, nil
					},
				}
				return aiService, articleRepo, tagRepo
			},
			expectedError: false,
			validateResult: func(t *testing.T, result *entity.Article) {
				assert.Equal(t, int64(2), result.ID)
				assert.Equal(t, "新しい技術記事", result.Title)
				assert.Equal(t, "https://example.com/article", result.URL)
				assert.Equal(t, "新しい技術の解説記事です。", result.Summary)
				assert.ElementsMatch(t, []string{"NewTech", "Tutorial"}, result.Tags)
				assert.Equal(t, "", result.Memo)
			},
		},
		{
			name: "正常系：タグなしで記事生成",
			url:  "https://example.com/article",
			memo: "",
			setupMocks: func() (*mockAIGeneratorService, *mockArticleRepository, *mockTagRepository) {
				aiService := &mockAIGeneratorService{
					generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
						return &service.GeneratedArticle{
							Title:         "タグなし記事",
							Summary:       "タグのない記事です。",
							SuggestedTags: []string{},
							SourceURL:     "https://example.com/article",
							TokenUsed:     80,
							GeneratedAt:   time.Now(),
						}, nil
					},
				}
				tagRepo := &mockTagRepository{}
				articleRepo := &mockArticleRepository{
					createFunc: func(ctx context.Context, article *entity.Article) (*entity.Article, error) {
						article.ID = 3
						return article, nil
					},
				}
				return aiService, articleRepo, tagRepo
			},
			expectedError: false,
			validateResult: func(t *testing.T, result *entity.Article) {
				assert.Equal(t, int64(3), result.ID)
				assert.Equal(t, "タグなし記事", result.Title)
				assert.Equal(t, []string{}, result.Tags)
			},
		},
		{
			name: "正常系：既存タグと新規タグの混在",
			url:  "https://example.com/article",
			memo: "混在タグテスト",
			setupMocks: func() (*mockAIGeneratorService, *mockArticleRepository, *mockTagRepository) {
				aiService := &mockAIGeneratorService{
					generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
						return &service.GeneratedArticle{
							Title:         "混在タグ記事",
							Summary:       "既存と新規タグが混在する記事です。",
							SuggestedTags: []string{"Go", "NewFramework"},
							SourceURL:     "https://example.com/article",
							TokenUsed:     120,
							GeneratedAt:   time.Now(),
						}, nil
					},
				}
				tagRepo := &mockTagRepository{
					findByNameFunc: func(ctx context.Context, name string) (*entity.Tag, error) {
						if name == "Go" {
							return &entity.Tag{ID: 1, Name: "Go"}, nil
						}
						return nil, errors.New("tag not found")
					},
					createFunc: func(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
						tag.ID = 10
						return tag, nil
					},
				}
				articleRepo := &mockArticleRepository{
					createFunc: func(ctx context.Context, article *entity.Article) (*entity.Article, error) {
						article.ID = 4
						return article, nil
					},
				}
				return aiService, articleRepo, tagRepo
			},
			expectedError: false,
			validateResult: func(t *testing.T, result *entity.Article) {
				assert.Equal(t, int64(4), result.ID)
				assert.ElementsMatch(t, []string{"Go", "NewFramework"}, result.Tags)
				assert.Equal(t, "混在タグテスト", result.Memo)
			},
		},
		{
			name: "異常系：URLが空の場合エラー",
			url:  "",
			memo: "",
			setupMocks: func() (*mockAIGeneratorService, *mockArticleRepository, *mockTagRepository) {
				return &mockAIGeneratorService{}, &mockArticleRepository{}, &mockTagRepository{}
			},
			expectedError:    true,
			expectedErrorMsg: "url is required",
		},
		{
			name: "異常系：URLが不正な形式の場合エラー",
			url:  "invalid-url",
			memo: "",
			setupMocks: func() (*mockAIGeneratorService, *mockArticleRepository, *mockTagRepository) {
				return &mockAIGeneratorService{}, &mockArticleRepository{}, &mockTagRepository{}
			},
			expectedError:    true,
			expectedErrorMsg: "invalid url format",
		},
		{
			name: "異常系：AI生成サービスがエラーを返す（API制限）",
			url:  "https://example.com/article",
			memo: "",
			setupMocks: func() (*mockAIGeneratorService, *mockArticleRepository, *mockTagRepository) {
				aiService := &mockAIGeneratorService{
					generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
						return nil, &service.AIGeneratorError{
							Code:    service.ErrCodeAPILimit,
							Message: "API rate limit exceeded",
							Err:     errors.New("too many requests"),
						}
					},
				}
				return aiService, &mockArticleRepository{}, &mockTagRepository{}
			},
			expectedError:    true,
			expectedErrorMsg: "API rate limit exceeded",
		},
		{
			name: "異常系：AI生成サービスがエラーを返す（タイムアウト）",
			url:  "https://example.com/article",
			memo: "",
			setupMocks: func() (*mockAIGeneratorService, *mockArticleRepository, *mockTagRepository) {
				aiService := &mockAIGeneratorService{
					generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
						return nil, &service.AIGeneratorError{
							Code:    service.ErrCodeTimeout,
							Message: "request timeout",
							Err:     errors.New("context deadline exceeded"),
						}
					},
				}
				return aiService, &mockArticleRepository{}, &mockTagRepository{}
			},
			expectedError:    true,
			expectedErrorMsg: "request timeout",
		},
		{
			name: "異常系：AI生成サービスがエラーを返す（不正なレスポンス）",
			url:  "https://example.com/article",
			memo: "",
			setupMocks: func() (*mockAIGeneratorService, *mockArticleRepository, *mockTagRepository) {
				aiService := &mockAIGeneratorService{
					generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
						return nil, &service.AIGeneratorError{
							Code:    service.ErrCodeInvalidResponse,
							Message: "invalid response format",
							Err:     errors.New("failed to parse response"),
						}
					},
				}
				return aiService, &mockArticleRepository{}, &mockTagRepository{}
			},
			expectedError:    true,
			expectedErrorMsg: "invalid response format",
		},
		{
			name: "異常系：生成されたタイトルが不正（空）",
			url:  "https://example.com/article",
			memo: "",
			setupMocks: func() (*mockAIGeneratorService, *mockArticleRepository, *mockTagRepository) {
				aiService := &mockAIGeneratorService{
					generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
						return &service.GeneratedArticle{
							Title:         "",
							Summary:       "要約です",
							SuggestedTags: []string{"Go"},
							SourceURL:     "https://example.com/article",
							TokenUsed:     100,
							GeneratedAt:   time.Now(),
						}, nil
					},
				}
				return aiService, &mockArticleRepository{}, &mockTagRepository{}
			},
			expectedError:    true,
			expectedErrorMsg: "title is required",
		},
		{
			name: "異常系：生成された要約が不正（空）",
			url:  "https://example.com/article",
			memo: "",
			setupMocks: func() (*mockAIGeneratorService, *mockArticleRepository, *mockTagRepository) {
				aiService := &mockAIGeneratorService{
					generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
						return &service.GeneratedArticle{
							Title:         "タイトル",
							Summary:       "",
							SuggestedTags: []string{"Go"},
							SourceURL:     "https://example.com/article",
							TokenUsed:     100,
							GeneratedAt:   time.Now(),
						}, nil
					},
				}
				return aiService, &mockArticleRepository{}, &mockTagRepository{}
			},
			expectedError:    true,
			expectedErrorMsg: "summary is required",
		},
		{
			name: "異常系：ArticleRepository保存時にエラー",
			url:  "https://example.com/article",
			memo: "",
			setupMocks: func() (*mockAIGeneratorService, *mockArticleRepository, *mockTagRepository) {
				aiService := &mockAIGeneratorService{
					generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
						return &service.GeneratedArticle{
							Title:         "記事タイトル",
							Summary:       "記事の要約",
							SuggestedTags: []string{"Go"},
							SourceURL:     "https://example.com/article",
							TokenUsed:     100,
							GeneratedAt:   time.Now(),
						}, nil
					},
				}
				tagRepo := &mockTagRepository{
					findByNameFunc: func(ctx context.Context, name string) (*entity.Tag, error) {
						return &entity.Tag{ID: 1, Name: "Go"}, nil
					},
				}
				articleRepo := &mockArticleRepository{
					createFunc: func(ctx context.Context, article *entity.Article) (*entity.Article, error) {
						return nil, errors.New("database connection error")
					},
				}
				return aiService, articleRepo, tagRepo
			},
			expectedError:    true,
			expectedErrorMsg: "database connection error",
		},
		{
			name: "異常系：TagRepository検索時にエラー（DB接続エラー）",
			url:  "https://example.com/article",
			memo: "",
			setupMocks: func() (*mockAIGeneratorService, *mockArticleRepository, *mockTagRepository) {
				aiService := &mockAIGeneratorService{
					generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
						return &service.GeneratedArticle{
							Title:         "記事タイトル",
							Summary:       "記事の要約",
							SuggestedTags: []string{"Go"},
							SourceURL:     "https://example.com/article",
							TokenUsed:     100,
							GeneratedAt:   time.Now(),
						}, nil
					},
				}
				tagRepo := &mockTagRepository{
					findByNameFunc: func(ctx context.Context, name string) (*entity.Tag, error) {
						return nil, errors.New("database connection failed")
					},
					createFunc: func(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
						return nil, errors.New("database connection failed")
					},
				}
				return aiService, &mockArticleRepository{}, tagRepo
			},
			expectedError:    true,
			expectedErrorMsg: "database connection failed",
		},
		{
			name: "異常系：TagRepository新規作成時にエラー",
			url:  "https://example.com/article",
			memo: "",
			setupMocks: func() (*mockAIGeneratorService, *mockArticleRepository, *mockTagRepository) {
				aiService := &mockAIGeneratorService{
					generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
						return &service.GeneratedArticle{
							Title:         "記事タイトル",
							Summary:       "記事の要約",
							SuggestedTags: []string{"NewTag"},
							SourceURL:     "https://example.com/article",
							TokenUsed:     100,
							GeneratedAt:   time.Now(),
						}, nil
					},
				}
				tagRepo := &mockTagRepository{
					findByNameFunc: func(ctx context.Context, name string) (*entity.Tag, error) {
						return nil, errors.New("tag not found")
					},
					createFunc: func(ctx context.Context, tag *entity.Tag) (*entity.Tag, error) {
						return nil, errors.New("failed to create tag: unique constraint violation")
					},
				}
				return aiService, &mockArticleRepository{}, tagRepo
			},
			expectedError:    true,
			expectedErrorMsg: "failed to create tag",
		},
		{
			name: "正常系：生成されたタグに空文字が含まれる場合、スキップして保存",
			url:  "https://example.com/article",
			memo: "",
			setupMocks: func() (*mockAIGeneratorService, *mockArticleRepository, *mockTagRepository) {
				aiService := &mockAIGeneratorService{
					generateFunc: func(ctx context.Context, req service.ArticleGenerationRequest) (*service.GeneratedArticle, error) {
						return &service.GeneratedArticle{
							Title:         "記事タイトル",
							Summary:       "記事の要約",
							SuggestedTags: []string{"Go", ""},
							SourceURL:     "https://example.com/article",
							TokenUsed:     100,
							GeneratedAt:   time.Now(),
						}, nil
					},
				}
				tagRepo := &mockTagRepository{
					findByNameFunc: func(ctx context.Context, name string) (*entity.Tag, error) {
						if name == "Go" {
							return &entity.Tag{ID: 1, Name: "Go"}, nil
						}
						return nil, errors.New("tag not found")
					},
				}
				articleRepo := &mockArticleRepository{
					createFunc: func(ctx context.Context, article *entity.Article) (*entity.Article, error) {
						article.ID = 1
						article.CreatedAt = time.Now()
						article.UpdatedAt = time.Now()
						return article, nil
					},
				}
				return aiService, articleRepo, tagRepo
			},
			expectedError: false,
			validateResult: func(t *testing.T, article *entity.Article) {
				assert.Equal(t, "記事タイトル", article.Title)
				assert.Equal(t, "https://example.com/article", article.URL)
				assert.Equal(t, "記事の要約", article.Summary)
				assert.Equal(t, []string{"Go"}, article.Tags)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aiService, articleRepo, tagRepo := tt.setupMocks()
			usecase := NewArticleGeneratorUsecase(aiService, articleRepo, tagRepo)

			result, err := usecase.GenerateArticleFromURL(context.Background(), tt.url, tt.memo)

			if tt.expectedError {
				require.Error(t, err)
				assert.Nil(t, result)
				if tt.expectedErrorMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrorMsg)
				}
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				if tt.validateResult != nil {
					tt.validateResult(t, result)
				}
			}
		})
	}
}
