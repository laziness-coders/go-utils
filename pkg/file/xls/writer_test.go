package xls

import (
	"context"
	"testing"
)

func TestNewWriter(t *testing.T) {
	type sheet1 struct {
		ColumnString      string   `xls:"String"`
		ColumnInt         int      `xls:"Int"`
		ColumnArrayString []string `xls:"Array String"`
	}

	type sheet2 struct {
		ColumnFloat64 float64 `xls:"Float64"`
		ColumnInt64   float64 `xls:"Int64"`
	}

	type testWriter struct {
		Sheet1Rows []*sheet1 `xls-sheet:"sheet dau tien"`
		Sheet2Rows []*sheet2 `xls-sheet:"sheet thu hai"`
	}

	writer := NewWriter()
	value := &testWriter{
		Sheet1Rows: []*sheet1{
			{
				ColumnString:      "ahihi",
				ColumnInt:         1,
				ColumnArrayString: []string{"string1", "string2"},
			},
			{
				ColumnString:      "ahhihi row 2",
				ColumnInt:         2,
				ColumnArrayString: []string{"string3"},
			},
		},
		Sheet2Rows: []*sheet2{},
	}

	bytes, err := writer.Render(context.Background(), value)
	if err != nil {
		t.Error(err)
	}

	t.Log("Bytes: ", string(bytes))
}
