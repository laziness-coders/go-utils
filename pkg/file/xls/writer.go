package xls

import (
	"context"
	"fmt"
	"github.com/laziness-coders/structs"
	"github.com/xuri/excelize/v2"
	"reflect"
	"strings"
)

type Writer interface {
	Render(ctx context.Context, data interface{}) ([]byte, error)
}

var _ Writer = (*writerImpl)(nil)

func NewWriter() Writer {
	return &writerImpl{}
}

type writerImpl struct {
}

func (x *writerImpl) Render(ctx context.Context, data interface{}) ([]byte, error) {
	if !isKindOf(data, reflect.Struct) {
		return nil, ErrorNotStruct
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheets := structs.Fields(data)
	for index, sheet := range sheets {
		if sheet.IsZero() {
			continue
		}

		if err := x.writeDataToSheet(index, f, sheet); err != nil {
			return nil, err
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (x writerImpl) writeDataToSheet(index int, f *excelize.File, sheet *structs.Field) error {
	sheetName := sheet.Tag(tagXlsSheet)
	if sheetName == "" {
		return nil
	}

	if err := x.upsertSheetName(f, index, sheetName); err != nil {
		return err
	}

	rows, err := x.validateSheetData(sheet.Elem())
	if err != nil {
		return err
	}

	for index := range rows {
		row := rows[index]
		axis := fmt.Sprintf("%s%d", firstColumnName, index+1) // A1, A2, A3, ...
		err := f.SetSheetRow(sheetName, axis, &row)
		if err != nil {
			return err
		}
	}

	return nil
}

func (x *writerImpl) upsertSheetName(f *excelize.File, index int, sheetName string) error {
	// if is first sheet then change first sheet name to current sheet name
	isFirstSheet := index == 0
	if isFirstSheet {
		if err := f.SetSheetName(defaultFirstSheetName, sheetName); err != nil {
			return err
		}
		return nil
	}

	if _, err := f.NewSheet(sheetName); err != nil {
		return err
	}

	return nil
}

func (x *writerImpl) validateSheetData(v interface{}) ([][]interface{}, error) {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Slice {
		return nil, ErrSheetDataNotSlice
	}

	sliceLen := rv.Len()
	onlyAddHeader := false
	if sliceLen == 0 {
		rv = addOneItemToSlice(rv)
		onlyAddHeader = true
		sliceLen = rv.Len()
	}

	resultLen := sliceLen + 1
	if onlyAddHeader {
		resultLen--
	}
	result := make([][]interface{}, resultLen)
	for rowIndex := 0; rowIndex < sliceLen; rowIndex++ {
		item := rv.Index(rowIndex)
		itemV := reflect.Indirect(item)
		rt := itemV.Type()
		numField := rt.NumField()
		nextRowIndex := rowIndex + 1

		for columnIndex := 0; columnIndex < numField; columnIndex++ {
			fieldValue := itemV.Field(columnIndex)
			fieldStruct := rt.Field(columnIndex)
			fieldKind := fieldValue.Kind()
			tag, ok := fieldStruct.Tag.Lookup(tagsXlsColumn)
			if !ok {
				return nil, ErrPleaseAddTagXlsToField
			}

			// Add header row
			if rowIndex == 0 {
				if columnIndex == 0 {
					result[rowIndex] = make([]interface{}, numField)
				}
				result[rowIndex][columnIndex] = tag
			}

			// Skip further processing if only adding header
			if onlyAddHeader {
				continue
			}

			if columnIndex == 0 {
				result[nextRowIndex] = make([]interface{}, numField)
			}
			if fieldKind != reflect.Slice {
				result[nextRowIndex][columnIndex] = fieldValue.Interface()
				continue
			}

			// special case for slice of string
			arrayValue, canParseArray := fieldValue.Interface().([]string)
			if !canParseArray {
				result[nextRowIndex][columnIndex] = fieldValue.Interface()
			} else {
				result[nextRowIndex][columnIndex] = strings.Join(arrayValue, sliceDelimiter)
			}
		}
	}

	return result, nil
}

func addOneItemToSlice(rv reflect.Value) reflect.Value {
	sliceType := reflect.SliceOf(rv.Type().Elem())
	elementType := sliceType.Elem()
	newElement := reflect.New(elementType.Elem())
	rv = reflect.Append(rv, newElement)
	return rv
}
