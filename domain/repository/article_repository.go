package repository

import (
	"github.com/Okaki030/hinagane-scraping/domain/model"
)

// ArticleRepository はまとめ記事の処理に必要なメソッドを定義するインターフェース
type ArticleRepository interface {
	InsertArticle(model.Article) (int, error)
	InsertMemberLinkToArticle(string, int) error
	InsertWord(string) (int, error)
	InsertWordLinkToArticle(int, int) error
}
