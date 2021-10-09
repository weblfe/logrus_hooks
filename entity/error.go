package entity

import "errors"

var (
	ErrNotExists = errors.New("hook not exist")
)

func IsNotExistsErr(err error) bool {
	return err == ErrNotExists
}
