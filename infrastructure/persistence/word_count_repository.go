package persistence

import (
	"fmt"

	"github.com/Okaki030/hinagane-scraping/domain/repository"
)

// wordCountPresistence はまとめ記事のワードの出現回数をカウントするための構造体
type wordCountPresistence struct{}

// NewWordCountPersistence はwordCountPresistence型のインスタンスを生成するための関数
func NewWordCountPersistence() repository.WordCountRepository {
	return &wordCountPresistence{}
}

// InsertWordCountInThreeDays は直近3日間のまとめ記事へのワードの出現回数をカウントするためのメソッド
func (wcp wordCountPresistence) InsertWordCountInThreeDays() error {
	fmt.Println("Insertwordcountinthreedays start")

	return nil
}
