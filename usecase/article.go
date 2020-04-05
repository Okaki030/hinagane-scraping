package usecase

import (
	"strings"

	"github.com/Okaki030/hinagane-scraping/domain/model"
	"github.com/Okaki030/hinagane-scraping/domain/repository"
	"github.com/PuerkitoBio/goquery"
	"github.com/shogo82148/go-mecab"
)

// ArticleUseCase はスクレイピングプログラムに必要なメソッドを定義するインターフェース
type ArticleUseCase interface {
	CollectArticle() error
}

// articleUseCase はスクレイピングプログラムに必要な構造体をまとめた構造体
type articleUseCase struct {
	articleRepository     repository.ArticleRepository
	memberCountRepository repository.MemberCountRepository
	wordCountRepository   repository.WordCountRepository
}

// NewArticleUseCase はarticleUseCase型のインスタンスを生成するための関数
func NewArticleUseCase(ar repository.ArticleRepository, mcr repository.MemberCountRepository, wcr repository.WordCountRepository) ArticleUseCase {
	return &articleUseCase{
		articleRepository:     ar,
		memberCountRepository: mcr,
		wordCountRepository:   wcr,
	}
}

// まとめ記事を登録
func (au articleUseCase) CollectArticle() error {

	var err error
	var ars, articles []model.Article

	// 日向坂まとめ速報からスクレイピング
	ars, err = ScrapingMatomesokuhou()
	if err != nil {
		return err
	}
	articles = append(articles, ars...)

	// まとめキングダムからスクレイピング
	ars, err = ScrapingMatomekingdom()
	if err != nil {
		return err
	}
	articles = append(articles, ars...)

	// 日向速報からスクレイピング
	ars, err = ScrapingHinatasokuhou()
	if err != nil {
		return err
	}
	articles = append(articles, ars...)

	for i, _ := range articles {

		// まとめ記事をdbに登録
		lastArticleId, err := au.articleRepository.InsertArticle(articles[i])
		if err != nil {
			return err
		}

		// 記事にメンバーカテゴリを紐付け
		for _, memberName := range articles[i].MemberNames {
			err = au.articleRepository.InsertMemberLinkToArticle(memberName, lastArticleId)
			if err != nil {
				return err
			}
		}

		// タイトルから単語を取得
		words, err := ExtractingWords(articles[i].Name)
		if err != nil {
			return err
		}
		articles[i].Words = append(articles[i].Words, words...)

		// 単語をdbに登録し記事に単語を紐付け
		for _, word := range articles[i].Words {

			// 単語をdbに登録
			_, err := au.articleRepository.InsertWord(word)
			if err != nil {
				return err
			}

			// 記事に単語を紐付け
			err = au.articleRepository.InsertWordLinkToArticle(word, lastArticleId)
			if err != nil {
				return err
			}
		}
	}

	// 直近3日間のまとめ記事へのメンバーの出現回数をカウント
	err = au.memberCountRepository.InsertMemberCountInThreeDays()
	if err != nil {
		return nil
	}

	err = au.wordCountRepository.InsertWordCountInThreeDays()
	if err != nil {
		return nil
	}

	return nil
}

// scrapingMatomesokuhou は日向坂まとめ速報の記事をスクレイピングする関数
func ScrapingMatomesokuhou() ([]model.Article, error) {

	targetUrl := "http://hiraganakeyaki.blog.jp/"
	doc, err := goquery.NewDocument(targetUrl)
	if err != nil {
		return nil, err
	}

	var articles []model.Article

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
		article.Url, _ = titleBox.Attr("href")

		// カテゴリー取得
		categorySet := articleBox.Find("ul.article-header-category")
		article.MemberNames = append(article.MemberNames, categorySet.Find("dd.article-category1").Text())
		article.MemberNames = append(article.MemberNames, categorySet.Find("dd.article-category2").Text())

		articles = append(articles, article)
	})

	return articles, nil
}

// scrapingMatomesokuhou は日向坂まとめキングダムの記事をスクレイピングする関数
func ScrapingMatomekingdom() ([]model.Article, error) {

	targetUrl := "http://hiragana46matome.com/"
	doc, err := goquery.NewDocument(targetUrl)
	if err != nil {
		return nil, err
	}

	var articles []model.Article

	// 1記事ずつまとめ記事を取得
	articleList := doc.Find("div.article")
	articleList.Each(func(i int, articleBox *goquery.Selection) {

		var article model.Article

		// サイトIDを取得
		article.SiteId = 2

		// タイトル取得
		article.Name = articleBox.Find("h3.article-title").Text()

		// url取得
		article.Url, _ = articleBox.Find("a").Attr("href")

		// カテゴリー取得
		categorySet := articleBox.Find("li.article-category-item")
		categorySet.Each(func(i int, category *goquery.Selection) {
			article.MemberNames = append(article.MemberNames, category.Text())
		})

		articles = append(articles, article)
	})

	return articles, nil
}

// scrapingMatomesokuhou は日向速報の記事をスクレイピングする関数
func ScrapingHinatasokuhou() ([]model.Article, error) {

	targetUrl := "http://hinatasoku.blog.jp/"
	doc, err := goquery.NewDocument(targetUrl)
	if err != nil {
		return nil, err
	}

	var articles []model.Article

	// 1記事ずつまとめ記事を取得
	articleList := doc.Find("header.article-header")
	articleList.Each(func(i int, articleBox *goquery.Selection) {

		var article model.Article

		// サイトID取得
		article.SiteId = 3

		// タイトル取得
		article.Name = articleBox.Find("h1.article-title").Text()

		// url取得
		article.Url, _ = articleBox.Find("h1.article-title").Find("a").Attr("href")

		// カテゴリー取得
		article.MemberNames = append(article.MemberNames, articleBox.Find("dd.article-category1").Text())
		article.MemberNames = append(article.MemberNames, articleBox.Find("dd.article-category2").Text())

		articles = append(articles, article)
	})

	return articles, nil
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

	var words []string
	for ; !node.IsZero(); node = node.Next() {
		slice := strings.Split(node.Feature(), ",")
		if slice[1] == "固有名詞" {
			words = append(words, node.Surface())
		}
	}

	return words, nil
}
