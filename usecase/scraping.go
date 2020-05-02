package usecase

import (
	"errors"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Okaki030/hinagane-scraping/domain/model"
	"github.com/PuerkitoBio/goquery"
	"github.com/oklog/ulid"
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
		article.Year, article.Month, article.Day, article.Hour = GetNow()

		// 画像取得
		picUrl, ok := articleBox.Find("img.pict").Attr("src")
		if ok == false {
			err = errors.New("Did not get Picture URL")
		}
		article.LocalPicPath, err = ScrapingPic(picUrl)

		articles = append(articles, article)
	})
	if err != nil {
		return nil, err
	}

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
		article.Year, article.Month, article.Day, article.Hour = GetNow()

		// 画像取得
		picUrl, ok := articleBox.Find("img").Attr("src")
		if ok == false {
			err = errors.New("Did not get Picture URL")
		}
		article.LocalPicPath, err = ScrapingPic(picUrl)

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
		article.Year, article.Month, article.Day, article.Hour = GetNow()

		// 画像取得
		picUrl, ok := articleBox.Find("img").Attr("src")
		if ok == false {
			err = errors.New("Did not get Picture URL")
		}
		article.LocalPicPath, err = ScrapingPic(picUrl)

		articles = append(articles, article)
	})

	return articles, nil
}

func GetNow() (int, int, int, int) {
	t := time.Now()

	return t.Year(), int(t.Month()), t.Day(), t.Hour()
}

// ScrapingPic は画像をスクレイピングするための関数
func ScrapingPic(picUrl string) (string, error) {
	response, err := http.Get(picUrl)
	if err != nil {
		return "", err
	}

	// 画像名(UUID)を生成
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	uuid := ulid.MustNew(ulid.Timestamp(t), entropy)

	picName := uuid.String() + ".jpg"

	file, err := os.Create(picName)
	if err != nil {
		return "", err
	}

	io.Copy(file, response.Body)
	response.Body.Close()
	file.Close()

	return picName, nil
}
