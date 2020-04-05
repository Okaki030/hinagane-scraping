package main

import (
	"log"

	"github.com/Okaki030/hinagane-scraping/infrastructure/persistence"
	"github.com/Okaki030/hinagane-scraping/usecase"
)

// main は最初に実行される関数
func main() {

	// 依存性注入
	articlePersistence := persistence.NewArticlePersistence()
	memberCountPersistence := persistence.NewMemberCountPersistence()
	wordCountPersistence := persistence.NewWordCountPersistence()
	articleUseCase := usecase.NewArticleUseCase(articlePersistence, memberCountPersistence, wordCountPersistence)

	err := articleUseCase.CollectArticle()
	if err != nil {
		log.Fatal(err)
	}
}
