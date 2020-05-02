package persistence

import (
	"github.com/aws/aws-sdk-go/aws/session"

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
func (mcp memberCountS3Presistence) InsertMemberCountInThreeDays() error {

	// var memberCnt, appearCnt int
	// var row *sql.Row
	// var err error

	// // メンバーの人数を取得
	// row = mcp.DB.QueryRow(`SELECT count(*) FROM member`)
	// err = row.Scan(&(memberCnt))
	// if err != nil {
	// 	return err
	// }

	// // 直近3日間のメンバーの記事数を取得する
	// for memberId := 1; memberId <= memberCnt; memberId++ {

	// 	row = mcp.DB.QueryRow(`
	// 		SELECT count(*) FROM article_member_link
	// 			INNER JOIN article on article_member_link.article_id=article.id
	// 			WHERE (NOW( ) - INTERVAL 3 DAY)<article.date_time and article_member_link.member_id=?`, memberId)
	// 	err = row.Scan(&(appearCnt))
	// 	if err != nil {
	// 		return err
	// 	}

	// 	_, err = mcp.DB.Exec(`
	// 		INSERT INTO
	// 			member_counter (member_id, counter, date_time)
	// 			VALUES (?,?,now())`, memberId, appearCnt)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
