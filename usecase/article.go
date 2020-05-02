package usecase

// ArticleUseCase はスクレイピングプログラムに必要なメソッドを定義するインターフェース
type ArticleUseCase interface {
	CollectArticle() error
}
