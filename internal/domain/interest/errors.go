package interest

import "errors"

var (
	ErrInterestAlreadyExist = errors.New("interest already exists")
	ErrInterestNotFound     = errors.New("interest not found")
)
