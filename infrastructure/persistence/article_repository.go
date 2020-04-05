package persistence

import (
	"database/sql"
	"fmt"

	"github.com/Okaki030/hinagane-scraping/domain/model"
	"github.com/Okaki030/hinagane-scraping/domain/repository"
)

// articlePersistence はまとめ記事の処理を扱うための構造体
type articlePersistence struct {
	DB *sql.DB
}

// NewArticlePersistence はarticlePersistenceのインスタンスを作成するための関数
func NewArticlePersistence(db *sql.DB) repository.ArticleRepository {
	return &articlePersistence{
		DB: db,
	}
}

// InsertArticle は1つのまとめ記事を保存するためのメソッド
func (ap articlePersistence) InsertArticle(article model.Article) (int, error) {

	fmt.Println("---Insertarticle start---")

	var err error

	// 記事をdbに追加
	res, err := ap.DB.Exec(`
		INSERT INTO 
		article (name,url,date_time,site_id) 
		VALUES (?,?,now(),?)`, article.Name, article.Url, article.SiteId)
	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	fmt.Println("挿入したID", lastId)

	return int(lastId), nil
}

// InsertMemberLinkToArticle は記事ごとのメンバーカテゴリを保存するためのメソッド
func (ap articlePersistence) InsertMemberLinkToArticle(memberName string, articleId int) error {

	fmt.Println("Insertmembercategory start")

	var memberId int

	// メンバー名からメンバーidを取得
	row := ap.DB.QueryRow(`SELECT id FROM member WHERE name=?`, memberName)
	_ = row.Scan(&(memberId))

	// 記事ごとにカテゴリ(メンバー名)を格納
	// TODO:1カテゴリ目でメンバーを取得し、2カテゴリ名で違った場合重複で登録使用しエラーを吐く
	if memberId != 0 {
		_, err := ap.DB.Exec(`
			INSERT INTO 
			article_member_link (article_id, member_id) 
			VALUES (?,?)`, articleId, memberId)
		if err != nil {
			// 何もしない
		}
	}

	return nil
}

// InsertWord は単語をdbに保存するメソッド
func (ap articlePersistence) InsertWord(word string) (int, error) {

	// 固有名詞をwordテーブルに登録
	res, err := ap.DB.Exec(`
		INSERT INTO word (name) 
		SELECT ? FROM dual 
		WHERE NOT EXISTS(SELECT * FROM word WHERE name=?);`, word, word)
	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	fmt.Println("挿入したwordID", lastId)

	return int(lastId), nil
}

// InsertWordLinkToArticle は記事ごとのワードを保存するためのメソッド
func (ap articlePersistence) InsertWordLinkToArticle(articleId int, wordId int) error {

	fmt.Println("Insertword start")

	// 記事ごとにカテゴリ(メンバー名)を格納
	// TODO:1カテゴリ目でメンバーを取得し、2カテゴリ名で違った場合重複で登録使用しエラーを吐く
	if wordId != 0 {
		_, err := ap.DB.Exec(`
			INSERT INTO 
			article_word_link (article_id, word_id) 
			VALUES (?,?)`, articleId, wordId)
		if err != nil {
			// 何もしない
		}
	}

	return nil
}
