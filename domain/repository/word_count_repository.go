package repository

// WordCountRepository はワードの出現回数を取得するのに必要なメソッドを定義するインターフェース
type WordCountRepository interface {
	InsertWordCountInThreeDays() error
}
