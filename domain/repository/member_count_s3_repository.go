package repository

// MemberCountS3Repository はメンバーの出現回数を取得するのに必要なメソッドを定義するインターフェース
type MemberCountS3Repository interface {
	InsertMemberCountInThreeDays(string) error
	ConfirmExistenceMemberCount(string, string) (bool, error)
	DownloadCSV() error
	SelectS3CSV(string, string) (string, error)
	UploadCSV() error
}
