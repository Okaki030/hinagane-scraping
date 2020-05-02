package repository

// MemberCountRepository はメンバーの出現回数を取得するのに必要なメソッドを定義するインターフェース
type MemberCountS3Repository interface {
	InsertMemberCountInThreeDays() error
	DownloadCSV() error
	SelectS3CSV(string, string) (string, error)
	UploadCSV() error
}
