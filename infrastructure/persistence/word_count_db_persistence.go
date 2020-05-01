package persistence

import (
	"database/sql"

	"github.com/Okaki030/hinagane-scraping/domain/repository"
)

// wordCountPresistence はまとめ記事のワードの出現回数をカウントするための構造体
type wordCountDBPresistence struct {
	DB *sql.DB
}

// NewWordCountPersistence はwordCountPresistence型のインスタンスを生成するための関数
func NewWordCountDBPersistence(db *sql.DB) repository.WordCountRepository {
	return &wordCountDBPresistence{
		DB: db,
	}
}

// InsertWordCountInThreeDays は直近3日間のまとめ記事へのワードの出現回数をカウントするためのメソッド
func (wcp wordCountDBPresistence) InsertWordCountInThreeDays() error {

	var wordCnt, appearCnt int
	var row *sql.Row
	var err error

	// 単語の総数を取得
	row = wcp.DB.QueryRow(`SELECT count(*) FROM word`)
	err = row.Scan(&(wordCnt))
	if err != nil {
		return err
	}

	// 直近3日間のメンバーの記事数を取得する
	for wordId := 1; wordId <= wordCnt; wordId++ {

		row = wcp.DB.QueryRow(`
			SELECT count(*) FROM article_word_link 
				INNER JOIN article on article_word_link.article_id=article.id 
				WHERE (NOW( ) - INTERVAL 3 DAY)<article.date_time and article_word_link.word_id=?`, wordId)
		err = row.Scan(&(appearCnt))
		if err != nil {
			return err
		}

		_, err = wcp.DB.Exec(`
			INSERT INTO 
				word_counter (word_id, counter, date_time) 
				VALUES (?,?,now())`, wordId, appearCnt)
		if err != nil {
			return err
		}
	}

	return nil
}
