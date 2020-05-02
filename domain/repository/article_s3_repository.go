package repository

import (
	"github.com/Okaki030/hinagane-scraping/domain/model"
)

// ArticleRepository はまとめ記事の処理に必要なメソッドを定義するインターフェース
type ArticleS3Repository interface {
	InsertArticle(model.Article, []string) (bool, error)
	ConfirmExistenceArticle(string) (bool, error)
	UploadArticlePic(string) (string, error)
	DownloadArticle() error
	UploadArticle() error
}
