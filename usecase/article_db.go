package usecase

// import (
// 	"github.com/Okaki030/hinagane-scraping/domain/repository"
// )

// // articleUseCase はスクレイピングプログラムに必要な構造体をまとめた構造体
// type articleUseCase struct {
// 	articleRepository     repository.ArticleRepository
// 	memberCountRepository repository.MemberCountRepository
// 	wordCountRepository   repository.WordCountRepository
// }

// // NewArticleUseCase はarticleUseCase型のインスタンスを生成するための関数
// func NewArticleUseCase(ar repository.ArticleRepository, mcr repository.MemberCountRepository, wcr repository.WordCountRepository) ArticleUseCase {
// 	return &articleUseCase{
// 		articleRepository:     ar,
// 		memberCountRepository: mcr,
// 		wordCountRepository:   wcr,
// 	}
// }

// // まとめ記事を登録
// func (au articleUseCase) CollectArticle() error {

// 	var err error

// 	articles, err := Scraping()
// 	if err != nil {
// 		return nil
// 	}

// 	for i, _ := range articles {

// 		// まとめ記事をdbに登録
// 		lastArticleId, err := au.articleRepository.InsertArticle(articles[i])
// 		if err != nil {
// 			return err
// 		}
// 		// すでに登録されている時は次の記事に飛ばす
// 		if lastArticleId == 0 {
// 			continue
// 		}

// 		// 記事にメンバーカテゴリを紐付け
// 		for _, memberName := range articles[i].MemberNames {
// 			err = au.articleRepository.InsertMemberLinkToArticle(memberName, lastArticleId)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		// タイトルから単語を取得
// 		words, err := ExtractingWords(articles[i].Name)
// 		if err != nil {
// 			return err
// 		}
// 		articles[i].Words = append(articles[i].Words, words...)

// 		// 単語をdbに登録し記事に単語を紐付け
// 		for _, word := range articles[i].Words {

// 			// 単語をdbに登録
// 			_, err := au.articleRepository.InsertWord(word)
// 			if err != nil {
// 				return err
// 			}

// 			// 記事に単語を紐付け
// 			err = au.articleRepository.InsertWordLinkToArticle(word, lastArticleId)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	// 直近3日間のまとめ記事へのメンバーの出現回数をカウント
// 	err = au.memberCountRepository.InsertMemberCountInThreeDays()
// 	if err != nil {
// 		return nil
// 	}

// 	err = au.wordCountRepository.InsertWordCountInThreeDays()
// 	if err != nil {
// 		return nil
// 	}

// 	return nil
// }
