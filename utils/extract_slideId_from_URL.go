package utils

import (
	"fmt"
	"regexp"
)

func ExtractSlideIDFromURL(url string) (string, error) {
	// Google SlidesのURL形式に一致するより汎用的な正規表現パターン
	// 末尾のスラッシュがオプショナルで、追加のパスやクエリパラメータがあっても対応
	pattern := `https://docs\.google\.com/presentation/d/([a-zA-Z0-9_-]+)`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err // 正規表現のコンパイルに失敗した場合のエラー
	}

	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("no slide ID found in URL")
	}

	// スライドIDを返します
	return matches[1], nil
}
