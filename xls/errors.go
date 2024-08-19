package xls

import "errors"

var (
	ErrPleaseAddTagXlsToField = errors.New("please add tag xls to field")
	ErrSheetDataNotSlice      = errors.New("sheet data not slice")
	ErrorNotStruct            = errors.New("not struct")
	ErrorOutOfRangeIndex      = errors.New("out of range index")
)
