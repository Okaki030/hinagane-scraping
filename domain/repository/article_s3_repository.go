package repository

import (
	"github.com/Okaki030/hinagane-scraping/domain/model"
)

// ArticleRepository はまとめ記事の処理に必要なメソッドを定義するインターフェース
type ArticleS3Repository interface {
	InsertArticle(model.Article, []string) (string, error)
	InsertMemberLinkToArticle(string, int) error
	InsertWord(string) (int, error)
	InsertWordLinkToArticle(string, int) error
	UploadArticlePic(string) (string, error)
	UploadArticle(string) error
}
