package gotestrail

import (
	"maps"
	"testing"
)

var parseCaseFieldItemsMapTests = []struct {
	v    string
	want map[uint]string
}{
	{"", map[uint]string{}},
	{"     ", map[uint]string{}},
	{"1, FooBar\n2, Foo\n3, Bar\n4, Baz\n5, Qux\n6, Quux_Corge", map[uint]string{
		1: "FooBar",
		2: "Foo",
		3: "Bar",
		4: "Baz",
		5: "Qux",
		6: "Quux_Corge"}},
}

func TestParseCaseFieldItemsMap(t *testing.T) {
	for _, tt := range parseCaseFieldItemsMapTests {
		try, err := parseCaseFieldItemsMap(tt.v)
		if err != nil {
			t.Errorf("gotestrail.parseCaseFieldItemsMap(%s) Error: (%s)", tt.v, err.Error())
		} else if !maps.Equal(tt.want, try) {
			t.Errorf("gotestrail.parseCaseFieldItemsMap(\"%s\") Mismatch: want (%v), got (%v)",
				tt.v, tt.want, try)
		}
	}
}
