package persistence

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/oklog/ulid"

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
func (ap articleS3Persistence) InsertArticle(article model.Article, words []string) (bool, error) {

	fmt.Println("Insert article start")
	var err error

	// 記事の存在を確認
	exist, err := ap.ConfirmExistenceArticle(article.Name)
	if err != nil {
		os.Remove(article.LocalPicPath)
		return false, err
	}
	// 記事が存在したら終わり
	if exist == true {
		os.Remove(article.LocalPicPath)
		return true, nil
	}

	// メンバースライスを"+"で結合
	memberNamesStr := strings.Join(article.MemberNames, "+")
	wordsStr := strings.Join(words, "+")

	// 画像を取得
	article.LocalPicPath, err = DownloadPic(article.PicUrl)

	// 画像をs3に保存
	article.S3PicUrl, err = ap.UploadArticlePic(article.LocalPicPath)
	if err != nil {
		return false, err
	}

	// 存在しなかったら作成、存在したら追記
	file, err := os.OpenFile("./articles.csv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return false, err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	// ファイル情報取得
	fileinfo, err := file.Stat()
	if err != nil {
		return false, err
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

	return false, nil
}

func (ap articleS3Persistence) ConfirmExistenceArticle(articleName string) (bool, error) {

	fmt.Println("confirm start")
	sql := "SELECT name FROM S3Object where name='" + articleName + "' LIMIT 1"
	fmt.Println(sql)
	svc := s3.New(ap.Sess)

	params := &s3.SelectObjectContentInput{
		Bucket:          aws.String("hinagane"),
		Key:             aws.String("./data/article/articles.csv"),
		ExpressionType:  aws.String(s3.ExpressionTypeSql),
		Expression:      aws.String("SELECT name FROM S3Object where name='" + articleName + "' LIMIT 1"),
		RequestProgress: &s3.RequestProgress{},
		InputSerialization: &s3.InputSerialization{
			CompressionType: aws.String("NONE"),
			CSV: &s3.CSVInput{
				FileHeaderInfo: aws.String(s3.FileHeaderInfoUse),
			},
		},
		OutputSerialization: &s3.OutputSerialization{
			CSV: &s3.CSVOutput{},
		},
	}

	resp, err := svc.SelectObjectContent(params)
	if err != nil {
		return false, err
	}
	defer resp.EventStream.Close()

	for event := range resp.EventStream.Events() {
		// 取得できたか判定する
		s, ok := event.(*s3.StatsEvent)
		if ok {
			if int(*s.Details.BytesReturned) > 0 {
				return true, nil
			}
		}
	}

	return false, nil
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

func (ap articleS3Persistence) DownloadArticle() error {

	fmt.Println("s3 download start")

	// S3オブジェクトを書き込むファイルの作成
	file, err := os.Create("./articles.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	bucketName := "hinagane"
	objectKey := "./data/article/articles.csv"

	// Downloaderを作成し、S3オブジェクトをダウンロード
	downloader := s3manager.NewDownloader(ap.Sess)
	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return err
	}

	return nil
}

// UploadArticlePic は記事の画像をs3にアップロードするメソッド
func (ap articleS3Persistence) UploadArticle() error {

	fmt.Println("csv upload article start")

	file, err := os.Open("articles.csv")
	if err != nil {
		return nil
	}
	defer file.Close()

	bucketName := "hinagane"
	objectKey := "./data/article/articles.csv"

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

	err = os.Remove("articles.csv")

	return nil
}

// DownloadPic は画像をスクレイピングするための関数
func DownloadPic(picUrl string) (string, error) {
	response, err := http.Get(picUrl)
	if err != nil {
		return "", err
	}

	// 画像名(UUID)を生成
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	uuid := ulid.MustNew(ulid.Timestamp(t), entropy)

	picName := uuid.String() + ".jpg"

	file, err := os.Create(picName)
	if err != nil {
		return "", err
	}

	io.Copy(file, response.Body)
	response.Body.Close()
	file.Close()

	return picName, nil
}
