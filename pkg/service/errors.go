package service

import (
	"errors"
)

var (
	ErrFeatureTagAlreadyExists = errors.New("feature tag already exists")
	ErrBannerNotFound          = errors.New("banner not found")
)
