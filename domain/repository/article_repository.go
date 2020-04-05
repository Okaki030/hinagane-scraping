package repository

import (
	"database/sql"

	"github.com/Okaki030/hinagane-scraping/domain/model"
)

// ArticleRepository はまとめ記事の処理に必要なメソッドを定義するインターフェース
type ArticleRepository interface {
	InsertArticle(*sql.DB, model.Article) (int, error)
	InsertMemberLinkToArticle(*sql.DB, string, int) error
	InsertWord(*sql.DB, string) (int, error)
	InsertWordLinkToArticle(*sql.DB, int, int) error
}
