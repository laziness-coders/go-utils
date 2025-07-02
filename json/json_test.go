package json

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type CompatTestStruct struct {
	ID         int            `json:"id"`
	Name       string         `json:"name"`
	Slice      []string       `json:"slice"`
	Map        map[string]int `json:"map"`
	Skip       string         `json:"-"`
	Omit       string         `json:"omit,omitempty"`
	Custom     string         `json:"custom_name"`
	OmitFilled string         `json:"omit_filled,omitempty"`
	Default    string         `json:"default_field"`
}

type ComplexStruct struct {
	ID        int                          `json:"id"`
	Name      string                       `json:"name"`
	Numbers   []int                        `json:"numbers"`
	Nested    *CompatTestStruct            `json:"nested"`
	NestedArr []CompatTestStruct           `json:"nested_arr"`
	MapData   map[string]*CompatTestStruct `json:"map_data"`
	CreatedAt time.Time                    `json:"created_at"`
	Active    bool                         `json:"active"`
	FloatVal  float64                      `json:"float_val"`
	PtrVal    *int                         `json:"ptr_val"`
	Omit      string                       `json:"omit,omitempty"`
	Custom    string                       `json:"custom_name"`
	Embedded  CompatTestStruct             `json:",inline"`
}

type CustomString string

type Embedded struct {
	Flag bool   `json:"flag"`
	Note string `json:"note"`
}

type VeryComplexStruct struct {
	ComplexStruct `json:",inline"`
	CustomField   CustomString                            `json:"custom_field"`
	DeepNested    *ComplexStruct                          `json:"deep_nested"`
	PtrSlice      []*CompatTestStruct                     `json:"ptr_slice"`
	MapOfSlices   map[string][]*CompatTestStruct          `json:"map_of_slices"`
	MapOfMaps     map[string]map[string]*CompatTestStruct `json:"map_of_maps"`
	Interface     interface{}                             `json:"interface"`
	Embedded      `json:",inline"`
	TimePtr       *time.Time `json:"time_ptr"`
	IntArr        [3]int     `json:"int_arr"`
	Skip          string     `json:"-"`
	Omit          string     `json:"omit,omitempty"`
	Custom        string     `json:"custom_name"`
}

var engineNames = []string{"std", "jsoniter", "gojson", "sonic"}
var marshalEngines = []Engine{
	NewStdEngine(),
	NewJsoniterEngine(),
	NewGoJsonEngine(),
	NewSonicEngine(),
}
var unmarshalEngines = []Engine{
	NewStdEngine(),
	NewJsoniterEngine(),
	NewGoJsonEngine(),
	NewSonicEngine(),
}

func TestEngineCompatibility(t *testing.T) {
	input := CompatTestStruct{
		ID:         42,
		Name:       "compat-test",
		Slice:      []string{"a", "b", "c"},
		Map:        map[string]int{"x": 1, "y": 2},
		Skip:       "should not appear",
		Omit:       "",
		Custom:     "custom value",
		OmitFilled: "not empty",
		Default:    "default",
	}
	expected := CompatTestStruct{
		ID:         42,
		Name:       "compat-test",
		Slice:      []string{"a", "b", "c"},
		Map:        map[string]int{"x": 1, "y": 2},
		Skip:       "",
		Omit:       "",
		Custom:     "custom value",
		OmitFilled: "not empty",
		Default:    "default",
	}

	for mi, mEngine := range marshalEngines {
		for ui, uEngine := range unmarshalEngines {
			SetMarshalEngine(mEngine)
			SetUnmarshalEngine(uEngine)

			data, err := Marshal(input)
			if err != nil {
				t.Fatalf("Marshal failed with marshal engine %s: %v", engineNames[mi], err)
			}

			var out CompatTestStruct
			err = Unmarshal(data, &out)
			if err != nil {
				t.Errorf("Unmarshal failed with marshal %s, unmarshal %s: %v", engineNames[mi], engineNames[ui], err)
				continue
			}

			if !cmp.Equal(expected, out) {
				t.Errorf("Mismatch with marshal %s, unmarshal %s:\n%s", engineNames[mi], engineNames[ui], cmp.Diff(expected, out))
			}
		}
	}
}

func TestEngineCompatibility_ComplexStruct(t *testing.T) {
	ptrVal := 99
	input := ComplexStruct{
		ID:      1001,
		Name:    "complex-test",
		Numbers: []int{1, 2, 3, 4, 5},
		Nested: &CompatTestStruct{
			ID:         7,
			Name:       "nested",
			Slice:      []string{"x", "y"},
			Map:        map[string]int{"a": 10, "b": 20},
			Skip:       "should not appear",
			Omit:       "",
			Custom:     "custom nested",
			OmitFilled: "filled",
			Default:    "default-nested",
		},
		NestedArr: []CompatTestStruct{
			{ID: 1, Name: "arr1", Slice: []string{"a"}, Map: map[string]int{"k": 1}, Skip: "should not appear", Omit: "", Custom: "arr1", OmitFilled: "", Default: "d1"},
			{ID: 2, Name: "arr2", Slice: []string{"b"}, Map: map[string]int{"l": 2}, Skip: "should not appear", Omit: "", Custom: "arr2", OmitFilled: "filled", Default: "d2"},
		},
		MapData: map[string]*CompatTestStruct{
			"first":  {ID: 3, Name: "map1", Slice: []string{"c"}, Map: map[string]int{"m": 3}, Skip: "should not appear", Omit: "", Custom: "map1", OmitFilled: "", Default: "dm1"},
			"second": {ID: 4, Name: "map2", Slice: []string{"d"}, Map: map[string]int{"n": 4}, Skip: "should not appear", Omit: "", Custom: "map2", OmitFilled: "filled", Default: "dm2"},
		},
		CreatedAt: time.Date(2023, 5, 1, 12, 0, 0, 0, time.UTC),
		Active:    true,
		FloatVal:  3.14159,
		PtrVal:    &ptrVal,
		Omit:      "",
		Custom:    "custom-complex",
		Embedded: CompatTestStruct{
			ID:         99,
			Name:       "embedded",
			Slice:      []string{"em"},
			Map:        map[string]int{"em": 99},
			Skip:       "should not appear",
			Omit:       "",
			Custom:     "em-custom",
			OmitFilled: "em-filled",
			Default:    "em-default",
		},
	}
	expected := ComplexStruct{
		ID:      1001,
		Name:    "complex-test",
		Numbers: []int{1, 2, 3, 4, 5},
		Nested: &CompatTestStruct{
			ID:         7,
			Name:       "nested",
			Slice:      []string{"x", "y"},
			Map:        map[string]int{"a": 10, "b": 20},
			Skip:       "",
			Omit:       "",
			Custom:     "custom nested",
			OmitFilled: "filled",
			Default:    "default-nested",
		},
		NestedArr: []CompatTestStruct{
			{ID: 1, Name: "arr1", Slice: []string{"a"}, Map: map[string]int{"k": 1}, Skip: "", Omit: "", Custom: "arr1", OmitFilled: "", Default: "d1"},
			{ID: 2, Name: "arr2", Slice: []string{"b"}, Map: map[string]int{"l": 2}, Skip: "", Omit: "", Custom: "arr2", OmitFilled: "filled", Default: "d2"},
		},
		MapData: map[string]*CompatTestStruct{
			"first":  {ID: 3, Name: "map1", Slice: []string{"c"}, Map: map[string]int{"m": 3}, Skip: "", Omit: "", Custom: "map1", OmitFilled: "", Default: "dm1"},
			"second": {ID: 4, Name: "map2", Slice: []string{"d"}, Map: map[string]int{"n": 4}, Skip: "", Omit: "", Custom: "map2", OmitFilled: "filled", Default: "dm2"},
		},
		CreatedAt: time.Date(2023, 5, 1, 12, 0, 0, 0, time.UTC),
		Active:    true,
		FloatVal:  3.14159,
		PtrVal:    &ptrVal,
		Omit:      "",
		Custom:    "custom-complex",
		Embedded: CompatTestStruct{
			ID:         99,
			Name:       "embedded",
			Slice:      []string{"em"},
			Map:        map[string]int{"em": 99},
			Skip:       "",
			Omit:       "",
			Custom:     "em-custom",
			OmitFilled: "em-filled",
			Default:    "em-default",
		},
	}

	for mi, mEngine := range marshalEngines {
		for ui, uEngine := range unmarshalEngines {
			SetMarshalEngine(mEngine)
			SetUnmarshalEngine(uEngine)

			data, err := Marshal(input)
			if err != nil {
				t.Fatalf("Marshal failed with marshal engine %s: %v", engineNames[mi], err)
			}

			var out ComplexStruct
			err = Unmarshal(data, &out)
			if err != nil {
				t.Errorf("Unmarshal failed with marshal %s, unmarshal %s: %v", engineNames[mi], engineNames[ui], err)
				continue
			}

			if !cmp.Equal(expected, out) {
				t.Errorf("Mismatch with marshal %s, unmarshal %s:\n%s", engineNames[mi], engineNames[ui], cmp.Diff(expected, out))
			}
		}
	}
}

func TestEngineCompatibility_VeryComplexStruct(t *testing.T) {
	ptrVal := 123
	tm := time.Date(2022, 12, 31, 23, 59, 59, 123456789, time.UTC)
	input := VeryComplexStruct{
		ComplexStruct: ComplexStruct{
			ID:      2022,
			Name:    "very-complex",
			Numbers: []int{10, 20, 30},
			Nested: &CompatTestStruct{
				ID:    8,
				Name:  "deep-nested",
				Slice: []string{"deep", "nest"},
				Map:   map[string]int{"deep": 100},
				Skip:  "should not appear",
				Omit:  "",
			},
			NestedArr: []CompatTestStruct{
				{ID: 3, Name: "arr3", Slice: []string{"c"}, Map: map[string]int{"m": 3}, Skip: "should not appear", Omit: ""},
				{ID: 4, Name: "arr4", Slice: []string{"d"}, Map: map[string]int{"n": 4}, Skip: "should not appear", Omit: ""},
			},
			MapData: map[string]*CompatTestStruct{
				"alpha": {ID: 5, Name: "map3", Slice: []string{"e"}, Map: map[string]int{"o": 5}, Skip: "should not appear", Omit: ""},
				"beta":  {ID: 6, Name: "map4", Slice: []string{"f"}, Map: map[string]int{"p": 6}, Skip: "should not appear", Omit: ""},
			},
			CreatedAt: tm,
			Active:    false,
			FloatVal:  2.71828,
			PtrVal:    &ptrVal,
			Omit:      "",
		},
		CustomField: "custom-str",
		DeepNested: &ComplexStruct{
			ID:      3033,
			Name:    "deeper",
			Numbers: []int{100, 200},
			Nested:  nil,
			NestedArr: []CompatTestStruct{
				{ID: 7, Name: "arr7", Slice: []string{"g"}, Map: map[string]int{"q": 7}, Skip: "should not appear", Omit: ""},
			},
			MapData:   map[string]*CompatTestStruct{},
			CreatedAt: tm.Add(-time.Hour),
			Active:    true,
			FloatVal:  1.41421,
			PtrVal:    nil,
			Omit:      "",
		},
		PtrSlice: []*CompatTestStruct{
			{ID: 9, Name: "ptr1", Slice: []string{"h"}, Map: map[string]int{"r": 9}, Skip: "should not appear", Omit: ""},
			nil,
		},
		MapOfSlices: map[string][]*CompatTestStruct{
			"group1": {
				{ID: 10, Name: "g1-1", Slice: []string{"i"}, Map: map[string]int{"s": 10}, Skip: "should not appear", Omit: ""},
				nil,
			},
		},
		MapOfMaps: map[string]map[string]*CompatTestStruct{
			"outer": {
				"inner": {ID: 11, Name: "in", Slice: []string{"j"}, Map: map[string]int{"t": 11}, Skip: "should not appear", Omit: ""},
			},
		},
		Interface: map[string]interface{}{
			"id":            float64(12),
			"name":          "iface",
			"slice":         []interface{}{"k"},
			"map":           map[string]interface{}{"u": float64(12)},
			"custom_name":   "",
			"omit_filled":   "",
			"default_field": "",
		},
		Embedded: Embedded{Flag: true, Note: "embedded"},
		TimePtr:  &tm,
		IntArr:   [3]int{1, 2, 3},
		Skip:     "should not appear",
		Omit:     "",
	}
	expected := VeryComplexStruct{
		ComplexStruct: ComplexStruct{
			ID:      2022,
			Name:    "very-complex",
			Numbers: []int{10, 20, 30},
			Nested: &CompatTestStruct{
				ID:    8,
				Name:  "deep-nested",
				Slice: []string{"deep", "nest"},
				Map:   map[string]int{"deep": 100},
				Skip:  "",
				Omit:  "",
			},
			NestedArr: []CompatTestStruct{
				{ID: 3, Name: "arr3", Slice: []string{"c"}, Map: map[string]int{"m": 3}, Skip: "", Omit: ""},
				{ID: 4, Name: "arr4", Slice: []string{"d"}, Map: map[string]int{"n": 4}, Skip: "", Omit: ""},
			},
			MapData: map[string]*CompatTestStruct{
				"alpha": {ID: 5, Name: "map3", Slice: []string{"e"}, Map: map[string]int{"o": 5}, Skip: "", Omit: ""},
				"beta":  {ID: 6, Name: "map4", Slice: []string{"f"}, Map: map[string]int{"p": 6}, Skip: "", Omit: ""},
			},
			CreatedAt: tm,
			Active:    false,
			FloatVal:  2.71828,
			PtrVal:    &ptrVal,
			Omit:      "",
		},
		CustomField: "custom-str",
		DeepNested: &ComplexStruct{
			ID:      3033,
			Name:    "deeper",
			Numbers: []int{100, 200},
			Nested:  nil,
			NestedArr: []CompatTestStruct{
				{ID: 7, Name: "arr7", Slice: []string{"g"}, Map: map[string]int{"q": 7}, Skip: "", Omit: ""},
			},
			MapData:   map[string]*CompatTestStruct{},
			CreatedAt: tm.Add(-time.Hour),
			Active:    true,
			FloatVal:  1.41421,
			PtrVal:    nil,
			Omit:      "",
		},
		PtrSlice: []*CompatTestStruct{
			{ID: 9, Name: "ptr1", Slice: []string{"h"}, Map: map[string]int{"r": 9}, Skip: "", Omit: ""},
			nil,
		},
		MapOfSlices: map[string][]*CompatTestStruct{
			"group1": {
				{ID: 10, Name: "g1-1", Slice: []string{"i"}, Map: map[string]int{"s": 10}, Skip: "", Omit: ""},
				nil,
			},
		},
		MapOfMaps: map[string]map[string]*CompatTestStruct{
			"outer": {
				"inner": {ID: 11, Name: "in", Slice: []string{"j"}, Map: map[string]int{"t": 11}, Skip: "", Omit: ""},
			},
		},
		Interface: map[string]interface{}{
			"id":            float64(12),
			"name":          "iface",
			"slice":         []interface{}{"k"},
			"map":           map[string]interface{}{"u": float64(12)},
			"custom_name":   "",
			"omit_filled":   "",
			"default_field": "",
		},
		Embedded: Embedded{Flag: true, Note: "embedded"},
		TimePtr:  &tm,
		IntArr:   [3]int{1, 2, 3},
		Skip:     "",
		Omit:     "",
	}

	for mi, mEngine := range marshalEngines {
		for ui, uEngine := range unmarshalEngines {
			SetMarshalEngine(mEngine)
			SetUnmarshalEngine(uEngine)

			data, err := Marshal(input)
			if err != nil {
				t.Fatalf("Marshal failed with marshal engine %s: %v", engineNames[mi], err)
			}

			var out VeryComplexStruct
			err = Unmarshal(data, &out)
			if err != nil {
				t.Errorf("Unmarshal failed with marshal %s, unmarshal %s: %v", engineNames[mi], engineNames[ui], err)
				continue
			}

			if !cmp.Equal(expected, out) {
				t.Errorf("Mismatch with marshal %s, unmarshal %s:\n%s", engineNames[mi], engineNames[ui], cmp.Diff(expected, out))
			}
		}
	}
}
