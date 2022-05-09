package internal

// FIXME: ここにあるstructは役割が曖昧なので消したい

type Problem struct {
	// URLに使われる、コンテスト内の問題を一意に識別するID
	ID string
	// 問題ページで表示される一意なID
	DisplayID string
}
