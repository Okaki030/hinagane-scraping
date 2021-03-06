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

// memberCountPresistence はまとめ記事のメンバーの出現回数をカウントするための構造体
type memberCountS3Presistence struct {
	Sess *session.Session
}

// NewMemberCountPersistence はmemberCountPresistence型のインスタンスを生成するための関数
func NewMemberCountS3Persistence(sess *session.Session) repository.MemberCountS3Repository {
	return &memberCountS3Presistence{
		Sess: sess,
	}
}

// InsertMemberCountInThreeDays は直近3日間のまとめ記事へのメンバーの出現回数をカウントするためのメソッド
func (mcp memberCountS3Presistence) InsertMemberCountInThreeDays(now string) error {

	fmt.Println("insert member count start")
	var err error

	// メンバー名を取得
	var memberStr string
	objectKey := "./data/member/member.csv"
	sql := "SELECT name FROM S3Object"
	memberStr, err = mcp.SelectS3CSV(objectKey, sql)
	if err != nil {
		return err
	}
	memberSlice := strings.Split(memberStr, "\n")
	memberSlice = memberSlice[:len(memberSlice)-1]
	pretty.Println(memberSlice)

	// ファイルオープン
	file, err := os.OpenFile("./member-count.csv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	// メンバーごとの3日寛の記事数を集計
	var memberCount string
	for _, member := range memberSlice {

		// 記事の存在を確認
		exist, err := mcp.ConfirmExistenceMemberCount(member, now)
		if err != nil {
			return err
		}
		// 記事が存在したら終わり
		if exist == true {
			continue
		}

		// メンバーの3日間の記事数を取得
		objectKey = "./data/article/articles.csv"
		sql := "SELECT count(*) FROM S3Object WHERE DATE_DIFF(hour, TO_TIMESTAMP(dateTime), DATE_ADD(hour,9,UTCNOW())) <= 72 AND memberNames like '%" + member + "%'"
		fmt.Println(sql)
		memberCount, err = mcp.SelectS3CSV(objectKey, sql)
		if err != nil {
			return err
		}
		memberCount = strings.ReplaceAll(memberCount, "\n", "")
		fmt.Println(member, memberCount)

		memberCountContent := []string{
			member,
			now,
			memberCount,
		}
		fmt.Println(memberCountContent)

		writer.Write(memberCountContent)
	}
	writer.Flush()

	return nil
}

func (mcp memberCountS3Presistence) ConfirmExistenceMemberCount(memberName string, now string) (bool, error) {

	fmt.Println("confirm membercount start")
	sql := "SELECT name FROM S3Object WHERE name='" + memberName + "' AND dateTime='" + now + "' LIMIT 1"
	fmt.Println(sql)
	svc := s3.New(mcp.Sess)

	params := &s3.SelectObjectContentInput{
		Bucket:          aws.String("hinagane"),
		Key:             aws.String("./data/member_count/member-count.csv"),
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

func (mcp memberCountS3Presistence) SelectS3CSV(objectKey string, sql string) (string, error) {

	var str string

	// メンバーの名前を取得
	svc := s3.New(mcp.Sess)

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

	fmt.Println("str", str)

	return str, nil
}

func (mcp memberCountS3Presistence) DownloadCSV() error {

	fmt.Println("memver count s3 download start")

	// S3オブジェクトを書き込むファイルの作成
	file, err := os.Create("./member-count.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	bucketName := "hinagane"
	objectKey := "./data/member_count/member-count.csv"

	// Downloaderを作成し、S3オブジェクトをダウンロード
	downloader := s3manager.NewDownloader(mcp.Sess)
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
func (mcp memberCountS3Presistence) UploadCSV() error {

	fmt.Println("csv upload membervcount start")

	file, err := os.Open("member-count.csv")
	if err != nil {
		return nil
	}
	defer file.Close()

	bucketName := "hinagane"
	objectKey := "./data/member_count/member-count.csv"

	// Uploaderを作成し、ローカルファイルをアップロード
	uploader := s3manager.NewUploader(mcp.Sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		return err
	}

	err = os.Remove("member-count.csv")

	return nil
}
