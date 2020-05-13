package usecase

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Okaki030/hinagane-scraping/domain/model"
	"github.com/PuerkitoBio/goquery"
	"github.com/shogo82148/go-mecab"
)

// Scraping はまとめ記事のスクレイピング関数をまとめた関数
func Scraping() ([]model.Article, error) {

	var err error
	var ars, articles []model.Article

	// 日向坂まとめ速報からスクレイピング
	ars, err = ScrapingMatomesokuhou()
	if err != nil {
		return nil, err
	}
	articles = append(articles, ars...)

	// まとめキングダムからスクレイピング
	ars, err = ScrapingMatomekingdom()
	if err != nil {
		return nil, err
	}
	articles = append(articles, ars...)

	// 日向速報からスクレイピング
	ars, err = ScrapingHinatasokuhou()
	if err != nil {
		return nil, err
	}
	articles = append(articles, ars...)

	return articles, nil
}

// ScrapingMatomesokuhou は日向坂まとめ速報の記事をスクレイピングする関数
func ScrapingMatomesokuhou() ([]model.Article, error) {

	targetUrl := "http://hiraganakeyaki.blog.jp/"
	doc, err := goquery.NewDocument(targetUrl)
	if err != nil {
		return nil, err
	}

	var articles []model.Article
	var ok bool

	// 1記事ずつまとめ記事を取得
	articleList := doc.Find("article.article")
	articleList.Each(func(i int, articleBox *goquery.Selection) {

		var article model.Article

		titleBox := articleBox.Find("h1.article-title").Find("a")

		// サイトIDを取得
		article.SiteId = 1

		// タイトル取得
		article.Name = titleBox.Text()

		// url取得
		article.Url, ok = titleBox.Attr("href")
		if ok == false {
			err = errors.New("Did not get Article URL")
		}

		// カテゴリー取得
		categorySet := articleBox.Find("ul.article-header-category")
		article.MemberNames = append(article.MemberNames, categorySet.Find("dd.article-category1").Text())
		article.MemberNames = append(article.MemberNames, categorySet.Find("dd.article-category2").Text())

		// 時間を取得
		article.DateTime = GetNow()

		// 画像取得
		article.PicUrl, ok = articleBox.Find("img.pict").Attr("src")
		if ok == false {
			err = errors.New("Did not get Picture URL")
		}
		// article.LocalPicPath, err = ScrapingPic(picUrl)

		articles = append(articles, article)
	})
	// if err != nil {
	// 	return nil, err
	// }

	return articles, nil
}

// ScrapingMatomesokuhou は日向坂まとめキングダムの記事をスクレイピングする関数
func ScrapingMatomekingdom() ([]model.Article, error) {

	targetUrl := "http://hiragana46matome.com/"
	doc, err := goquery.NewDocument(targetUrl)
	if err != nil {
		return nil, err
	}

	var articles []model.Article
	var ok bool

	// 1記事ずつまとめ記事を取得
	articleList := doc.Find("div.article")
	articleList.Each(func(i int, articleBox *goquery.Selection) {

		var article model.Article

		// サイトIDを取得
		article.SiteId = 2

		// タイトル取得
		article.Name = articleBox.Find("h3.article-title").Text()

		// url取得
		article.Url, ok = articleBox.Find("a").Attr("href")
		if ok == false {
			err = errors.New("Did not get Article URL")
		}

		// カテゴリー取得
		categorySet := articleBox.Find("li.article-category-item")
		categorySet.Each(func(i int, category *goquery.Selection) {
			article.MemberNames = append(article.MemberNames, category.Text())
		})

		// 時間を取得
		article.DateTime = GetNow()

		// 画像取得
		article.PicUrl, ok = articleBox.Find("img").Attr("src")
		if ok == false {
			err = errors.New("Did not get Picture URL")
		}
		// article.LocalPicPath, err = ScrapingPic(picUrl)

		articles = append(articles, article)
	})
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// ScrapingMatomesokuhou は日向速報の記事をスクレイピングする関数
func ScrapingHinatasokuhou() ([]model.Article, error) {

	targetUrl := "http://hinatasoku.blog.jp/"
	doc, err := goquery.NewDocument(targetUrl)
	if err != nil {
		return nil, err
	}

	var articles []model.Article
	var ok bool

	// 1記事ずつまとめ記事を取得
	articleList := doc.Find("article.article")
	articleList.Each(func(i int, articleBox *goquery.Selection) {

		var article model.Article

		// サイトID取得
		article.SiteId = 3

		// タイトル取得
		article.Name = articleBox.Find("h1.article-title").Text()

		// url取得
		article.Url, ok = articleBox.Find("h1.article-title").Find("a").Attr("href")
		if ok == false {
			err = errors.New("Did not get Article URL")
		}

		// カテゴリー取得
		article.MemberNames = append(article.MemberNames, articleBox.Find("dd.article-category1").Text())
		article.MemberNames = append(article.MemberNames, articleBox.Find("dd.article-category2").Text())

		// 時間を取得
		article.DateTime = GetNow()

		// 画像取得
		article.PicUrl, ok = articleBox.Find("img").Attr("src")
		if ok == false {
			err = errors.New("Did not get Picture URL")
		}
		// article.LocalPicPath, err = ScrapingPic(picUrl)

		articles = append(articles, article)
	})

	return articles, nil
}

// GetNow は現在時刻を取得する関数
func GetNow() string {
	t := time.Now()

	month := formatDateTime(int(t.Month()))
	day := formatDateTime(t.Day())
	hour := formatDateTime(t.Hour())

	dateTime := strconv.Itoa(t.Year()) + "-" + month + "-" + day + "T" + hour + ":00:00Z"

	return dateTime
}

func formatDateTime(dateTimeOne int) string {
	dateTimeOneStr := ""
	if dateTimeOne < 10 {
		dateTimeOneStr = "0" + strconv.Itoa(dateTimeOne)
	} else {
		dateTimeOneStr = strconv.Itoa(dateTimeOne)
	}

	return dateTimeOneStr
}

// ExtractingWords は固有名詞を抽出する関数
func ExtractingWords(title string) ([]string, error) {

	dic := make(map[string]string)
	dic["dicdir"] = "/usr/local/lib/mecab/dic/mecab-ipadic-neologd"

	mecab, err := mecab.New(dic)
	if err != nil {
		return nil, err
	}
	defer mecab.Destroy()

	node, err := mecab.ParseToNode(title)
	if err != nil {
		return nil, err
	}

	stopWords := []string{"小坂菜緒", "日向坂46", "日向坂", "", "www", "wwww", "wwwww", "wwwwww", "wwwwwww", "wwwwwwww", "wwwwwwwww", "wwwwwwwwww", "ｗｗｗｗｗｗｗｗｗ", "ｗｗｗｗｗｗ", "丹生"}

	var words []string
	for ; !node.IsZero(); node = node.Next() {
		stopFlag := false
		slice := strings.Split(node.Feature(), ",")
		if slice[1] == "固有名詞" {
			for _, s := range stopWords {
				if s == node.Surface() {
					stopFlag = true
				}
			}
			if !stopFlag {
				words = append(words, node.Surface())
			}
		}
	}

	return words, nil
}
