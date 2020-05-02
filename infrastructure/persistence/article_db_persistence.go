package persistence

import (
	"database/sql"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/Okaki030/hinagane-scraping/domain/model"
	"github.com/Okaki030/hinagane-scraping/domain/repository"
)

// articlePersistence はまとめ記事の処理を扱うための構造体
type articleDBPersistence struct {
	DB   *sql.DB
	Sess *session.Session
}

// NewArticlePersistence はarticlePersistenceのインスタンスを作成するための関数
func NewArticleDBPersistence(db *sql.DB, sess *session.Session) repository.ArticleRepository {
	return &articleDBPersistence{
		DB:   db,
		Sess: sess,
	}
}

// InsertArticle は1つのまとめ記事を保存するためのメソッド
func (ap articleDBPersistence) InsertArticle(article model.Article) (int, error) {

	var err error
	var articleId int

	// 記事がすでに登録されていないかチェックする
	row := ap.DB.QueryRow(`SELECT id FROM article WHERE name=?`, article.Name)
	err = row.Scan(&(articleId))
	if articleId != 0 {
		err = os.Remove(article.LocalPicPath)
		return 0, nil
	}

	// 画像をs3に保存
	article.S3PicUrl, err = ap.UploadArticlePic(article.LocalPicPath)
	if err != nil {
		return 0, err
	}

	// 記事をdbに追加
	res, err := ap.DB.Exec(`
		INSERT INTO 
		article (name,url,date_time,site_id,pic_url) 
		VALUES (?,?,now(),?,?)`, article.Name, article.Url, article.SiteId, article.S3PicUrl)
	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastId), nil
}

// InsertMemberLinkToArticle は記事ごとのメンバーカテゴリを保存するためのメソッド
func (ap articleDBPersistence) InsertMemberLinkToArticle(memberName string, articleId int) error {

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
func (ap articleDBPersistence) InsertWord(word string) (int, error) {

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

	return int(lastId), nil
}

// InsertWordLinkToArticle は記事ごとのワードを保存するためのメソッド
func (ap articleDBPersistence) InsertWordLinkToArticle(word string, articleId int) error {

	var wordId int

	// メンバー名からメンバーidを取得
	row := ap.DB.QueryRow(`SELECT id FROM word WHERE name=?`, word)
	_ = row.Scan(&(wordId))

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

// UploadArticlePic は記事の画像をs3にアップロードするメソッド
func (ap articleDBPersistence) UploadArticlePic(picName string) (string, error) {

	file, err := os.Open(picName)
	if err != nil {
		return "", nil
	}
	defer file.Close()

	bucketName := "hinagane"
	objectKey := "./pic/article/" + picName

	// Uploaderを作成し、ローカルファイルをアップロード
	uploader := s3manager.NewUploader(ap.Sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		return "", err
	}

	err = os.Remove(picName)

	return result.Location, nil
}
