package persistence

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
func (ap articleS3Persistence) InsertArticle(article model.Article, words []string, csvName string) error {

	fmt.Println("Insert article start")
	var err error

	// 記事の存在を確認

	// メンバースライスを","で結合
	memberNamesStr := strings.Join(article.MemberNames, ",")
	wordsStr := strings.Join(words, ",")

	// 画像をs3に保存
	article.S3PicUrl, err = ap.UploadArticlePic(article.LocalPicPath)
	if err != nil {
		return err
	}

	// 存在しなかったら作成、存在したら追記
	file, err := os.OpenFile("./"+csvName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	// ファイル情報取得
	fileinfo, err := file.Stat()
	if err != nil {
		return err
	}

	// ファイルサイズを表示
	fmt.Println(fileinfo.Size())

	// ファイルに何も書き込まれていないときにヘッダー名をcsvに書き込む
	if fileinfo.Size() == 0 {
		var header = []string{
			"name",
			"url",
			"year",
			"month",
			"day",
			"hour",
			"memberNames",
			"words",
			"siteId",
			"s3PicUrl",
		}
		writer.Write(header)
	}

	// 記事データをcsvに書き込む
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
	writer.Write(articleContent)
	writer.Flush()

	return nil
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

func (ap articleS3Persistence) DownloadArticle() (string, error) {

	t := time.Now()

	csvName := strconv.Itoa(t.Year()) + strconv.Itoa(int(t.Month())) + strconv.Itoa(t.Day()) + ".csv"

	// S3オブジェクトを書き込むファイルの作成
	file, err := os.Create("./" + csvName)
	if err != nil {
		return "", err
	}

	bucketName := "hinagane"
	objectKey := "./data/article/" + csvName

	// Downloaderを作成し、S3オブジェクトをダウンロード
	downloader := s3manager.NewDownloader(ap.Sess)
	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return csvName, nil
	}

	return csvName, nil
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
