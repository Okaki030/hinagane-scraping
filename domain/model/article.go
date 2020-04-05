package model

import "time"

// Article はまとめ記事を扱うための構造体
type Article struct {
	Name        string
	Url         string
	DateTime    time.Time
	MemberNames []string
	Words       []string
	SiteId      int
}
