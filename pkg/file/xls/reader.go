package xls

import (
	"context"
	"fmt"
	"github.com/laziness-coders/structs"
	"github.com/xuri/excelize/v2"
	"io"
	"reflect"
	"strings"
)

type Reader interface {
	Read(ctx context.Context, r io.Reader, data interface{}) error
}

var _ Reader = (*readerImpl)(nil)

func NewReader() Reader {
	return &readerImpl{}
}

type readerImpl struct {
}

func (x *readerImpl) Read(ctx context.Context, r io.Reader, data interface{}) error {
	file, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	if !isKindOf(data, reflect.Struct) {
		return ErrorNotStruct
	}

	fields := structs.Fields(data)
	for _, field := range fields {
		kind := field.Kind()
		if kind != reflect.Slice {
			return ErrSheetDataNotSlice
		}

		sheetName := field.Tag(tagXlsSheet)
		rows, err := file.Rows(sheetName)
		if err != nil {
			return err
		}

		rowsData := make([][]string, 0)
		for rows.Next() {
			row, err := rows.Columns()
			if err != nil {
				return err
			}

			if row == nil {
				continue
			}

			rowsData = append(rowsData, row)
		}

		if err := x.copyRowsToSlice(rowsData, field); err != nil {
			return err
		}

		if err := rows.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (x *readerImpl) copyRowsToSlice(rows [][]string, toSlice *structs.Field) error {
	onlyHaveHeader := len(rows) <= 1
	if onlyHaveHeader {
		return nil
	}

	// rows[0] is columns name of sheet
	columnsName := rows[0]
	mapColumnNameWithIndex := make(map[interface{}]int, 0)
	for i, columnName := range columnsName {
		mapColumnNameWithIndex[columnName] = i
	}

	value := toSlice.Value()
	sliceElemType := reflect.TypeOf(value).Elem()
	newValue := reflect.MakeSlice(reflect.SliceOf(sliceElemType), 0, len(rows)-1)
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		toStruct := reflect.New(sliceElemType.Elem())
		err := x.copyRowToStruct(row, mapColumnNameWithIndex, toStruct.Interface())
		if err != nil {
			return err
		}
		newValue = reflect.Append(newValue, toStruct)
	}

	iNewValue := newValue.Interface()
	err := toSlice.Set(iNewValue)
	if err != nil {
		return err
	}

	return nil
}

func (x *readerImpl) copyRowToStruct(row []string, mapColumnNameWithIndex map[interface{}]int, toStruct interface{}) error {
	lengthRow := len(row)
	fields := structs.Fields(toStruct)
	for _, field := range fields {
		tagXlsName := field.Tag(tagsXlsColumn)
		columnIndex, ok := mapColumnNameWithIndex[tagXlsName]
		if !ok {
			continue
		}

		isOutOfRange := lengthRow-columnIndex-1 < 0
		if isOutOfRange {
			return ErrorOutOfRangeIndex
		}

		stringValue := row[columnIndex]
		to, err := x.getToValue(field, stringValue)
		if err != nil {
			return err
		}

		if err := field.Set(to); err != nil {
			return err
		}
	}
	return nil
}

func (x *readerImpl) getToValue(field *structs.Field, stringValue string) (interface{}, error) {
	to := field.Value()
	fieldKind := field.Kind()
	switch fieldKind {
	case reflect.Slice:
		// remove [ and ] from string
		stringValue = strings.Trim(stringValue, sliceBrackets)
		arr := strings.Split(stringValue, sliceDelimiter)
		_ = arr
		// if err := utils.Copy(&to, arr); err != nil {
		// 	return nil, err
		// }
		return to, nil
	case reflect.Int, reflect.Float32, reflect.Float64, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// remove "," from string
		stringValue = strings.ReplaceAll(stringValue, stringNumberDelimiter, emptyString)
	}

	// if err := utils.Copy(&to, stringValue); err != nil {
	// 	return nil, err
	// }
	return to, nil
}
