package snapper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gregoryjjb/snapper"
	"github.com/gregoryjjb/snapper/otherpkg"
)

type SampleStruct struct {
	foo float64
	Bar string
	Baz int
}

type NestedStruct struct {
	Value float64
	Inner SampleStruct
}

func TestSnap(t *testing.T) {

	cases := []struct {
		name   string
		input  any
		expect string
	}{
		{name: "nil", input: nil, expect: "nil"},
		{name: "bool true", input: true, expect: "true"},
		{name: "bool false", input: false, expect: "false"},
		{
			name:   "int",
			input:  420,
			expect: "420",
		},
		{
			name:   "int8",
			input:  int8(32),
			expect: "32",
		},
		{
			name:   "int64",
			input:  int64(1825),
			expect: "1825",
		},
		{
			name:   "float",
			input:  12.34,
			expect: "12.34",
		},
		{
			name:   "string",
			input:  "foo\nbar",
			expect: `"foo\nbar"`,
		},
		{
			name:  "struct",
			input: SampleStruct{Bar: "asdf", Baz: 420},
			expect: `SampleStruct{
	Bar: "asdf",
	Baz: 420,
}`,
		},
		{
			name: "nested structs",
			input: NestedStruct{
				Value: 56.78,
				Inner: SampleStruct{
					Bar: "aaa",
					Baz: 69,
				},
			},
			expect: `NestedStruct{
	Value: 56.78,
	Inner: SampleStruct{
		Bar: "aaa",
		Baz: 69,
	},
}`,
		},
		{
			name:  "pointer to struct",
			input: &SampleStruct{Bar: "jkl", Baz: 456},
			expect: `&SampleStruct{
	Bar: "jkl",
	Baz: 456,
}`,
		},
		{
			name:   "empty array",
			input:  []any{},
			expect: `[]any{}`,
		},
		{
			name:  "int array",
			input: []int{1, 2, 3},
			expect: `[]int{
	1,
	2,
	3,
}`,
		},
		{
			name: "array struct",
			input: []SampleStruct{
				{
					Bar: "aaa",
					Baz: 13,
				},
			},
			expect: `[]SampleStruct{
	{
		Bar: "aaa",
		Baz: 13,
	},
}`,
		},
		{
			name: "array struct nested",
			input: []NestedStruct{
				{
					Value: 34.56,
					Inner: SampleStruct{
						Bar: "aa",
						Baz: 16,
					},
				},
			},
			expect: `[]NestedStruct{
	{
		Value: 34.56,
		Inner: SampleStruct{
			Bar: "aa",
			Baz: 16,
		},
	},
}`,
		},
		{
			name: "map",
			input: map[string]SampleStruct{
				"foo": {
					Bar: "a",
					Baz: 2,
				},
			},
			expect: `map[string]SampleStruct{
	"foo": {
		Bar: "a",
		Baz: 2,
	},
}`,
		},
		{
			name: "struct external package",
			input: otherpkg.User{
				Name: "bob",
				Orders: []otherpkg.Order{
					{
						Id: 420,
					},
				},
			},
			expect: `otherpkg.User{
	Name: "bob",
	Orders: []otherpkg.Order{
		{
			Id: 420,
		},
	},
}`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actual := snapper.Ssnap(tt.input, map[string]string{"snapper_test": ""})
			assert.Equal(t, tt.expect, actual)
		})
	}
}
