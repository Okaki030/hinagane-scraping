package usecase

// import (
// 	"reflect"
// 	"strconv"
// 	"testing"

// 	"github.com/Okaki030/hinagane-scraping/domain/model"
// )

// // TestScrapingMatomesokuhou はScrapingMatomesokuhou関数のテスト
// func TestScrapingMatomesokuhou(t *testing.T) {

// 	// テスト対象関数の実行
// 	articles, err := ScrapingMatomesokuhou()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// 取得できる記事数のテスト
// 	if len(articles) != 8 {
// 		t.Errorf("The number of retrieved articles is abnormal")
// 	}

// 	CommonScrapingTest(t, articles)
// }

// // TestScrapingMatomekingdom はScrapingMatomekingdom関数のテスト
// func TestScrapingMatomekingdom(t *testing.T) {

// 	// テスト対象関数の実行
// 	articles, err := ScrapingMatomekingdom()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// 取得できる記事数のテスト
// 	if len(articles) != 25 {
// 		t.Errorf("The number of retrieved articles is abnormal")
// 	}

// 	CommonScrapingTest(t, articles)
// }

// // TestScrapingHinatasokuhou はScrapingHinatasokuhou関数のテスト
// func TestScrapingHinatasokuhou(t *testing.T) {

// 	// テスト対象関数の実行
// 	articles, err := ScrapingHinatasokuhou()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// 取得できる記事数のテスト
// 	if len(articles) != 17 {
// 		t.Errorf("The number of retrieved articles is abnormal")
// 	}

// 	CommonScrapingTest(t, articles)
// }

// // CommonScrapingTest はScraping関数の共通テスト
// func CommonScrapingTest(t *testing.T, articles []model.Article) {

// 	for i, article := range articles {

// 		// タイトルが取得できているか
// 		if article.Name == "" {
// 			t.Errorf("Number %s article was not retrieved", strconv.Itoa(i))
// 		}

// 		// URLが取得できているか
// 		if article.Name == "" {
// 			t.Errorf("Number %s article was not retrieved", strconv.Itoa(i))
// 		}
// 	}
// }

// // TestExtractingWords はExtractingWords関数のテスト
// func TestExtractingWords(t *testing.T) {

// 	cases := []struct {
// 		name  string
// 		input string
// 		want  []string
// 	}{
// 		{
// 			name:  "normal case",
// 			input: "3月27日放送の「しくじり先生　春の特別授業SP」(テレビ朝日系)に、日向坂46・小坂菜緒が生徒役で登場。",
// 			want:  []string{"3月27日", "しくじり先生", "SP", "テレビ朝日", "日向坂46", "小坂菜緒"},
// 		},
// 	}

// 	for _, c := range cases {
// 		words, err := ExtractingWords(c.input)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		if r := reflect.DeepEqual(words, c.want); r == false {
// 			t.Errorf("%s want %s, got %s", c.name, c.want, words)
// 		}
// 	}
// }
