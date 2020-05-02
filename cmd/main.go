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
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// s3 connect
	sess, err := config.ConnectS3()
	if err != nil {
		log.Fatal(err)
	}

	// 依存性注入(MySQL,S3)
	articlePersistence := persistence.NewArticleDBPersistence(db, sess)
	memberCountPersistence := persistence.NewMemberCountDBPersistence(db)
	wordCountPersistence := persistence.NewWordCountDBPersistence(db)
	articleUseCase := usecase.NewArticleUseCase(articlePersistence, memberCountPersistence, wordCountPersistence)

	// // 依存性注入(S3)
	// articlePersistence := persistence.NewArticleS3Persistence(sess)
	// memberCountPersistence := persistence.NewMemberCountS3Persistence(sess)
	// wordCountPersistence := persistence.NewWordCountS3Persistence(sess)
	// articleUseCase := usecase.NewArticleUseCase(articlePersistence, memberCountPersistence, wordCountPersistence)

	err = articleUseCase.CollectArticle()
	if err != nil {
		log.Fatal(err)
	}
}
