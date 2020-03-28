package jsonstreamsplit

import (
	"reflect"
	"strings"
	"testing"
)

func TestSplitString(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    []string
		wantErr bool
	}{
		{"simple", "{a}{b}", []string{"{a}", "{b}"}, false},
		{"with a delimiter", "{a}notjson{b}", []string{"{a}", "{b}"}, false},
		{"nesting", "{a:{c}}{b}", []string{"{a:{c}}", "{b}"}, false},
		{"quoted", `{a:"{c"}{b}`, []string{`{a:"{c"}`, "{b}"}, false},
		{"escaped quote", `{a:"{c\""}{b}`, []string{`{a:"{c\""}`, "{b}"}, false},
		{"just a slash", `{a:"{c\something"}{b}`, []string{`{a:"{c\something"}`, "{b}"}, false},
		{"double escape quote", `{a:"\\"}{b}`, []string{`{a:"\\"}`, "{b}"}, false},
		{"non-ascii", `{a:"世界"}{b}`, []string{`{a:"世界"}`, "{b}"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SplitString(strings.NewReader(tt.arg))
			if (err != nil) != tt.wantErr {
				t.Errorf("Split() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Split() got = %v, want %v", got, tt.want)
			}
		})
	}
}
