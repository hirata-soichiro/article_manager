package entity

import (
    "errors"
    "strings"
    "time"
    "unicode/utf8"
)

// 購入リンク
type PurchaseLinks struct {
    Amazon  string
    Rakuten string
}

// 書籍
type Book struct {
    Title         string
    Author        string
    ISBN          string
    PurchaseLinks PurchaseLinks
}

// 書籍推薦キャッシュ
type BookRecommendationCache struct {
    ID          int64
    Books       []Book
    GeneratedAt time.Time
    ExpiresAt   time.Time
}

// 新しい書籍の作成
func NewBook(title, author, isbn string, purchaseLinks PurchaseLinks) (*Book, error) {
    if err := validateBookTitle(title); err != nil {
        return nil, err
    }
    if err := validateAuthor(author); err != nil {
        return nil, err
    }
    if err := validateISBN(isbn); err != nil {
        return nil, err
    }
    if err := validatePurchaseLinks(purchaseLinks); err != nil {
        return nil, err
    }

    book := &Book{
        Title:         title,
        Author:        author,
        ISBN:          isbn,
        PurchaseLinks: purchaseLinks,
    }

    return book, nil
}

// 新しい書籍推薦キャッシュの作成
func NewBookRecommendationCache(books []Book) (*BookRecommendationCache, error) {
    if err := validateBooks(books); err != nil {
        return nil, err
    }

    now := time.Now()
    expiresAt := now.Add(24 * time.Hour)

    cache := &BookRecommendationCache{
        ID:          0,
        Books:       books,
        GeneratedAt: now,
        ExpiresAt:   expiresAt,
    }

    return cache, nil
}

// キャッシュが有効化どうかを確認
func (c *BookRecommendationCache) IsValid() bool {
    return time.Now().Before(c.ExpiresAt)
}

// バリデーション関数
func validateBookTitle(title string) error {
    if title == "" {
        return errors.New("book title is required")
    }
    if utf8.RuneCountInString(title) > 500 {
        return errors.New("book title must be 500 characters or less")
    }
    return nil
}

func validateAuthor(author string) error {
    if author == "" {
        return errors.New("author is required")
    }
    if utf8.RuneCountInString(author) > 255 {
        return errors.New("author must be 255 characters or less")
    }
    return nil
}

func validateISBN(isbn string) error {
    if author == "" {
        return nil
    }
    if len(isbn) != 10 && len(isbn) != 13 {
        return errors.New("isbn must be 10 or 13 characters")
    }
    return nil
}

func validatePurchaseLinks(links PurchaseLinks) error {
    if links.Amazon != "" {
        if !strings.HasPrefix(links.Amazon, "http://") && !strings.HasPrefix(links.Amazon, "https://") {
            return errors.New("amazon link must start with http:// or https://")
        }
    }
    if links.Rakuten != "" {
        if !strings.HasPrefix(links.Rakuten, "http://") && !strings.HasPrefix(links.Rakuten, "https://") {
            return errors.New("rakuten link must start with http:// or https://")
        }
    }
    return nil
}

func validateBooks(books []Book) error {
    if books == nil {
        return errors.New("books cannot be nil")
    }

    for i, book := range books {
        if book.Title == "" {
            return errors.New("book title is required at index " + string(rune(i)))
        }
        if book.Author == "" {
            return errors.New("author is required at index " + string(rune(i)))
        }
    }

    return nil
}