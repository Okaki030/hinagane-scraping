package model

// Article はまとめ記事を扱うための構造体
type Article struct {
	Name         string
	Url          string
	Year         int
	Month        int
	Day          int
	Hour         int
	MemberNames  []string
	Words        []string
	SiteId       int
	PicUrl       string
	LocalPicPath string
	S3PicUrl     string
}
