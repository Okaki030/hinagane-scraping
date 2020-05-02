package usecase

import (
	"github.com/Okaki030/hinagane-scraping/domain/repository"
)

// articleUseCase はスクレイピングプログラムに必要な構造体をまとめた構造体
type articleS3UseCase struct {
	articleRepository     repository.ArticleS3Repository
	memberCountRepository repository.MemberCountS3Repository
	wordCountRepository   repository.WordCountS3Repository
}

// NewArticleUseCase はarticleUseCase型のインスタンスを生成するための関数
func NewArticleS3UseCase(ar repository.ArticleS3Repository, mcr repository.MemberCountS3Repository, wcr repository.WordCountS3Repository) ArticleUseCase {
	return &articleS3UseCase{
		articleRepository:     ar,
		memberCountRepository: mcr,
		wordCountRepository:   wcr,
	}
}

// まとめ記事を登録
func (au articleS3UseCase) CollectArticle() error {

	var err error

	articles, err := Scraping()
	if err != nil {
		return nil
	}

	// 記事(csv)をs3から取得
	err = au.articleRepository.DownloadArticle()
	if err != nil {
		return err
	}

	for i, _ := range articles {

		var err error

		// タイトルから固有名詞を取得
		words, err := ExtractingWords(articles[i].Name)
		if err != nil {
			return err
		}

		// まとめ記事をdbに登録
		exist, err := au.articleRepository.InsertArticle(articles[i], words)
		if err != nil {
			return err
		}
		if exist == true {
			continue
		}
	}
	// 記事データをs3に保存
	au.articleRepository.UploadArticle()

	// メンバーカウントcsvを取得
	err = au.memberCountRepository.DownloadCSV()
	if err != nil {
		return err
	}

	// 直近3日間のまとめ記事へのメンバーの出現回数をカウント
	err = au.memberCountRepository.InsertMemberCountInThreeDays()
	if err != nil {
		return nil
	}

	// メンバーカウントcsvをアップロード
	err = au.memberCountRepository.UploadCSV()
	if err != nil {
		return nil
	}

	return nil
}
