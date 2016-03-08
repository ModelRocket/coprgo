package coprhd

import (
	"strings"
)

const (
	ErrCodeOK               = 0
	ErrCodeInvalidParam     = 1008
	ErrCodeCreateNotAllowed = 1054
)

func (err ApiError) IsOK() bool {
	return err.Code == ErrCodeOK
}

func (err ApiError) IsCreateVolDup() bool {
	if err.Code != ErrCodeInvalidParam {
		return false
	}

	return strings.Contains(err.Details, "already exists")
}

func (err ApiError) IsExportVolDup() bool {
	if err.Code != ErrCodeCreateNotAllowed {
		return false
	}

	return strings.Contains(err.Details, "already exists")
}

func (err ApiError) IsCreateHostDup() bool {
	if err.Code != ErrCodeInvalidParam {
		return false
	}

	return strings.Contains(err.Details, "already exists")
}
