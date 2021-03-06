package repository

// WordCountRepository はワードの出現回数を取得するのに必要なメソッドを定義するインターフェース
type WordCountS3Repository interface {
	InsertWordCountInThreeDays(string) error
	ConfirmExistenceWordCount(string, string) (bool, error)
	DownloadCSV() error
	SelectS3CSV(string, string) (string, error)
	UploadCSV() error
}
