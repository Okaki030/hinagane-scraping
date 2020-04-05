package persistence

import (
	"fmt"

	"github.com/Okaki030/hinagane-scraping/domain/repository"
)

// memberCountPresistence はまとめ記事のメンバーの出現回数をカウントするための構造体
type memberCountPresistence struct{}

// NewMemberCountPersistence はmemberCountPresistence型のインスタンスを生成するための関数
func NewMemberCountPersistence() repository.MemberCountRepository {
	return &memberCountPresistence{}
}

// InsertMemberCountInThreeDays は直近3日間のまとめ記事へのメンバーの出現回数をカウントするためのメソッド
func (mcp memberCountPresistence) InsertMemberCountInThreeDays() error {
	fmt.Println("Insertmembercountinthreedays start")

	return nil
}
