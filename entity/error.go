package entity

import "errors"

var (
	// ErrNotExists hook 不存在
	ErrNotExists = errors.New("hook not exist")
	// ErrTagParseFailed 解析 tag
	ErrTagParseFailed = errors.New("tag parse failed")
)

func IsNotExistsErr(err error) bool {
	return err == ErrNotExists
}

func IsTagParseFailed(err error) bool  {
		return err == ErrTagParseFailed
}
