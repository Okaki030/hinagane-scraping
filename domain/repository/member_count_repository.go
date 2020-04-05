package repository

// MemberCountRepository はメンバーの出現回数を取得するのに必要なメソッドを定義するインターフェース
type MemberCountRepository interface {
	InsertMemberCountInThreeDays() error
}
