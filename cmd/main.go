package main

import (
	"log"

	"github.com/Okaki030/hinagane-scraping/config"
	"github.com/Okaki030/hinagane-scraping/infrastructure/persistence"
	"github.com/Okaki030/hinagane-scraping/usecase"
)

// main は最初に実行される関数
func main() {

	// db connect
	db, err := config.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 依存性注入
	articlePersistence := persistence.NewArticlePersistence(db)
	memberCountPersistence := persistence.NewMemberCountPersistence(db)
	wordCountPersistence := persistence.NewWordCountPersistence(db)
	articleUseCase := usecase.NewArticleUseCase(articlePersistence, memberCountPersistence, wordCountPersistence)

	err = articleUseCase.CollectArticle()
	if err != nil {
		log.Fatal(err)
	}
}
