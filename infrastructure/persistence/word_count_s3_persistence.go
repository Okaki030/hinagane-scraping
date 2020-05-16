package persistence

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/kr/pretty"

	"github.com/Okaki030/hinagane-scraping/domain/repository"
)

// wordCountPresistence はまとめ記事のワードの出現回数をカウントするための構造体
type wordCountS3Presistence struct {
	Sess *session.Session
}

// NewWordCountS3Persistence はwordCountPresistence型のインスタンスを生成するための関数
func NewWordCountS3Persistence(sess *session.Session) repository.WordCountS3Repository {
	return &wordCountS3Presistence{
		Sess: sess,
	}
}

// InsertWordCountInThreeDays は直近3日間のまとめ記事へのワードの出現回数をカウントするためのメソッド
func (wcp wordCountS3Presistence) InsertWordCountInThreeDays(now string) error {

	fmt.Println("insert word count start")
	var err error

	// 単語を取得
	var wordStr string
	var wordSlice []string
	objectKey := "./data/article/articles.csv"
	sql := "SELECT words FROM s3object"
	wordStr, err = wcp.SelectS3CSV(objectKey, sql)
	if err != nil {
		return err
	}
	wordStr = strings.ReplaceAll(wordStr, "\"", "")
	wordsSlice := strings.Split(wordStr, "\n")
	wordsSlice = wordsSlice[:len(wordsSlice)-1]
	pretty.Println("wordsSlice", wordsSlice)
	for _, words := range wordsSlice {
		wordSlice = append(wordSlice, strings.Split(words, "+")...)
	}
	pretty.Println(wordSlice)

	m := make(map[string]bool)
	uniqWordSlice := []string{}

	for _, ele := range wordSlice {
		if !m[ele] {
			m[ele] = true
			uniqWordSlice = append(uniqWordSlice, ele)
		}
	}

	pretty.Println(uniqWordSlice)

	// ファイルオープン
	file, err := os.OpenFile("./word-count.csv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	// 単語ごとの3日寛の記事数を集計
	var wordCount string
	for _, word := range uniqWordSlice {

		// 記事の存在を確認
		exist, err := wcp.ConfirmExistenceWordCount(word, now)
		if err != nil {
			return err
		}
		// 記事が存在したら終わり
		if exist == true {
			continue
		}

		// 単語ごとの3日寛の記事数を取得
		objectKey = "./data/article/articles.csv"
		sql := "SELECT count(*) FROM S3Object WHERE DATE_DIFF(hour, TO_TIMESTAMP(dateTime), DATE_ADD(hour,9,UTCNOW())) <= 72 AND  words like '%" + word + "%'"
		fmt.Println(sql)
		wordCount, err = wcp.SelectS3CSV(objectKey, sql)
		if err != nil {
			return err
		}
		wordCount = strings.ReplaceAll(wordCount, "\n", "")

		wordCountContent := []string{
			word,
			now,
			wordCount,
		}
		fmt.Println(wordCountContent)
		writer.Write(wordCountContent)
	}
	writer.Flush()

	return nil
}

func (wcp wordCountS3Presistence) ConfirmExistenceWordCount(wordName string, now string) (bool, error) {

	fmt.Println("confirm membercount start")
	sql := "SELECT name FROM S3Object WHERE name='" + wordName + "' AND dateTime='" + now + "' LIMIT 1"
	fmt.Println(sql)
	svc := s3.New(wcp.Sess)

	params := &s3.SelectObjectContentInput{
		Bucket:          aws.String("hinagane"),
		Key:             aws.String("./data/word_count/word-count.csv"),
		ExpressionType:  aws.String(s3.ExpressionTypeSql),
		Expression:      aws.String(sql),
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

func (wcp wordCountS3Presistence) SelectS3CSV(objectKey string, sql string) (string, error) {

	var str string

	// メンバーの名前を取得
	svc := s3.New(wcp.Sess)

	params := &s3.SelectObjectContentInput{
		Bucket:          aws.String("hinagane"),
		Key:             aws.String(objectKey),
		ExpressionType:  aws.String(s3.ExpressionTypeSql),
		Expression:      aws.String(sql),
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
		return "", err
	}
	defer resp.EventStream.Close()

	for event := range resp.EventStream.Events() {
		v, ok := event.(*s3.RecordsEvent)
		if ok {
			str = string(v.Payload)
		}
	}

	return str, nil
}

func (wcp wordCountS3Presistence) DownloadCSV() error {

	fmt.Println("word count s3 download start")

	// S3オブジェクトを書き込むファイルの作成
	file, err := os.Create("./word-count.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	bucketName := "hinagane"
	objectKey := "./data/word_count/word-count.csv"

	// Downloaderを作成し、S3オブジェクトをダウンロード
	downloader := s3manager.NewDownloader(wcp.Sess)
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
func (wcp wordCountS3Presistence) UploadCSV() error {

	fmt.Println("csv upload wordcount start")

	file, err := os.Open("word-count.csv")
	if err != nil {
		return nil
	}
	defer file.Close()

	bucketName := "hinagane"
	objectKey := "./data/word_count/word-count.csv"

	// Uploaderを作成し、ローカルファイルをアップロード
	uploader := s3manager.NewUploader(wcp.Sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		return err
	}

	err = os.Remove("word-count.csv")

	return nil
}
