package persistence

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/Okaki030/hinagane-scraping/domain/model"
	"github.com/Okaki030/hinagane-scraping/domain/repository"
)

// articlePersistence はまとめ記事の処理を扱うための構造体
type articleS3Persistence struct {
	Sess *session.Session
}

// NewArticlePersistence はarticlePersistenceのインスタンスを作成するための関数
func NewArticleS3Persistence(sess *session.Session) repository.ArticleS3Repository {
	return &articleS3Persistence{
		Sess: sess,
	}
}

// InsertArticle は1つのまとめ記事を保存するためのメソッド
func (ap articleS3Persistence) InsertArticle(article model.Article, words []string) (string, error) {

	fmt.Println("Insert article start")
	var err error

	// 記事の存在を確認

	// メンバースライスを","で結合
	memberNamesStr := strings.Join(article.MemberNames, ",")
	wordsStr := strings.Join(words, ",")

	// 画像をs3に保存
	article.S3PicUrl, err = ap.UploadArticlePic(article.LocalPicPath)
	if err != nil {
		return "", err
	}

	// csvに書き込むデータをsliceにまとめる
	var articleContent = []string{
		article.Name,
		article.Url,
		strconv.Itoa(article.Year),
		strconv.Itoa(article.Month),
		strconv.Itoa(article.Day),
		strconv.Itoa(article.Hour),
		memberNamesStr,
		wordsStr,
		strconv.Itoa(article.SiteId),
		article.S3PicUrl,
	}

	// csvに書き込み
	csvName := strconv.Itoa(article.Year) + strconv.Itoa(article.Month) + strconv.Itoa(article.Day) + ".csv"

	// 存在しなかったら作成、存在したら追記
	file, err := os.OpenFile("./"+csvName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Write(articleContent)
	writer.Flush()

	return csvName, nil
}

// InsertMemberLinkToArticle は記事ごとのメンバーカテゴリを保存するためのメソッド
func (ap articleS3Persistence) InsertMemberLinkToArticle(memberName string, articleId int) error {

	// var memberId int

	// // メンバー名からメンバーidを取得
	// row := ap.DB.QueryRow(`SELECT id FROM member WHERE name=?`, memberName)
	// _ = row.Scan(&(memberId))

	// // 記事ごとにカテゴリ(メンバー名)を格納
	// // TODO:1カテゴリ目でメンバーを取得し、2カテゴリ名で違った場合重複で登録使用しエラーを吐く
	// if memberId != 0 {
	// 	_, err := ap.DB.Exec(`
	// 		INSERT INTO
	// 		article_member_link (article_id, member_id)
	// 		VALUES (?,?)`, articleId, memberId)
	// 	if err != nil {
	// 		// 何もしない
	// 	}
	// }

	// return nil
	return nil
}

// InsertWord は単語をdbに保存するメソッド
func (ap articleS3Persistence) InsertWord(word string) (int, error) {

	// // 固有名詞をwordテーブルに登録
	// res, err := ap.DB.Exec(`
	// 	INSERT INTO word (name)
	// 		SELECT ? FROM dual
	// 		WHERE NOT EXISTS(SELECT * FROM word WHERE name=?);`, word, word)
	// if err != nil {
	// 	return 0, err
	// }

	// lastId, err := res.LastInsertId()
	// if err != nil {
	// 	return 0, err
	// }

	// return int(lastId), nil
	return 10, nil
}

// InsertWordLinkToArticle は記事ごとのワードを保存するためのメソッド
func (ap articleS3Persistence) InsertWordLinkToArticle(word string, articleId int) error {

	// var wordId int

	// // メンバー名からメンバーidを取得
	// row := ap.DB.QueryRow(`SELECT id FROM word WHERE name=?`, word)
	// _ = row.Scan(&(wordId))

	// // 記事ごとにカテゴリ(メンバー名)を格納
	// // TODO:1カテゴリ目でメンバーを取得し、2カテゴリ名で違った場合重複で登録使用しエラーを吐く
	// if wordId != 0 {
	// 	_, err := ap.DB.Exec(`
	// 		INSERT INTO
	// 		article_word_link (article_id, word_id)
	// 		VALUES (?,?)`, articleId, wordId)
	// 	if err != nil {
	// 		// 何もしない
	// 	}
	// }

	return nil
}

// UploadArticlePic は記事の画像をs3にアップロードするメソッド
func (ap articleS3Persistence) UploadArticlePic(picName string) (string, error) {

	fmt.Println("uploadarticle pic start")

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

// UploadArticlePic は記事の画像をs3にアップロードするメソッド
func (ap articleS3Persistence) UploadArticle(csvName string) error {

	fmt.Println("csv upload article start")

	file, err := os.Open(csvName)
	if err != nil {
		return nil
	}
	defer file.Close()

	bucketName := "hinagane"
	objectKey := "./data/article/" + csvName

	// Uploaderを作成し、ローカルファイルをアップロード
	uploader := s3manager.NewUploader(ap.Sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		return err
	}

	err = os.Remove(csvName)

	return nil
}
