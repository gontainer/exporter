// Copyright (c) 2023–present Bartłomiej Krukowski
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package exporter //nolint:testpackage

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type (
	myString    string
	aliasString = string
	myInt       int
	aliasInt    = int
	myBool      bool
	aliasBool   = bool
)

//nolint:testifylint
func TestChainExporter_Export(t *testing.T) {
	t.Parallel()

	t.Run("Given scenarios", func(t *testing.T) {
		t.Parallel()

		//nolint:exhaustruct
		scenarios := map[string]struct {
			input  any
			output string
			error  string
		}{
			"nil": {
				input:  nil,
				output: "nil",
			},
			"false": {
				input:  false,
				output: "false",
			},
			"true": {
				input:  true,
				output: "true",
			},
			"123": {
				input:  int(123),
				output: "int(123)",
			},
			"`hello world`": {
				input:  "hello world",
				output: `"hello world"`,
			},
			"[]byte": {
				input:  []byte("hello world 你好，世界"),
				output: `[]byte("hello world \u4f60\u597d\uff0c\u4e16\u754c")`,
			},
			"struct {}": {
				input: struct{}{},
				error: "type struct {} is not supported",
			},
			"*testing.T": {
				input: t,
				error: "type *testing.T is not supported",
			},
			`myString("foo")`: {
				input: myString("foo"),
				error: "type exporter.myString is not supported",
			},
			`aliasString("foo")`: {
				input:  aliasString("foo"),
				output: `"foo"`,
			},
			`myInt(5)`: {
				input: myInt(5),
				error: "type exporter.myInt is not supported",
			},
			`aliasInt(5)`: {
				input:  aliasInt(5),
				output: "int(5)",
			},
			`myBool(true)`: {
				input: myBool(true),
				error: "type exporter.myBool is not supported",
			},
			`aliasBool(true)`: {
				input:  aliasBool(true),
				output: "true",
			},
			`([][2][][]interface{})(nil)`: {
				input:  ([][2][][]interface{})(nil),
				output: `([][2][][]interface{})(nil)`,
			},
			`[][2][][]int{{{{1, 2}}, nil}}`: {
				input:  [][2][][]int{{{{1, 2}}, nil}},
				output: `[][2][][]int{[2][][]int{[][]int{[]int{int(1), int(2)}}, ([][]int)(nil)}}`,
			},
			`[][]any{nil, nil, {(*int)(nil)}}`: {
				input: [][]any{nil, nil, {(*int)(nil)}},
				error: `cannot export ([][]interface{})[2]: cannot export ([]interface{})[0]: type *int is not supported`,
			},
			`[]any{(*int)(nil)}`: {
				input: []any{(*int)(nil)},
				error: `cannot export ([]interface{})[0]: type *int is not supported`,
			},
			`[0][][]any{}`: {
				input:  [0][][]any{},
				output: `[0][][]interface{}{}`,
			},
			`[0][][]any{} #2`: {
				input:  [0][][]interface{}{},
				output: `[0][][]interface{}{}`,
			},
			`[]any{[][]int{{1, 2}, {3, 4}}, ([][][]any)(nil)}`: {
				input:  []any{[][]int{{1, 2}, {3, 4}}, ([][][]any)(nil)},
				output: `[]interface{}{[][]int{[]int{int(1), int(2)}, []int{int(3), int(4)}}, ([][][]interface{})(nil)}`,
			},
		}

		for k, s := range scenarios {
			s := s
			t.Run(k, func(t *testing.T) {
				t.Parallel()

				output, err := Export(s.input)
				if s.error != "" {
					assert.EqualError(t, err, s.error)
					assert.Equal(t, "", output)

					return
				}
				assert.NoError(t, err)
				assert.Equal(t, s.output, output)
			})
		}
	})
}

//nolint:testifylint
func TestExport(t *testing.T) {
	t.Parallel()

	//nolint:exhaustruct
	scenarios := []struct {
		input  any
		output string
		error  string
		panic  string
	}{
		{
			input:  123,
			output: "int(123)",
		},
		{
			input:  []any{1, "2", 3.14},
			output: `[]interface{}{int(1), "2", float64(3.14)}`,
		},
		{
			input:  [3]any{1, "2", 3.14},
			output: `[3]interface{}{int(1), "2", float64(3.14)}`,
		},
		{
			input:  []any{},
			output: "make([]interface{}, 0)",
		},
		{
			input:  ([]interface{})(nil),
			output: "([]interface{})(nil)",
		},
		{
			input:  ([]any)(nil),
			output: "([]interface{})(nil)",
		},
		{
			input:  [0]any{},
			output: "[0]interface{}{}",
		},
		{
			input: []any{struct{}{}},
			error: "cannot export ([]interface{})[0]: type struct {} is not supported",
			panic: "cannot export []interface {} to string: cannot export ([]interface{})[0]: type struct {} is not supported",
		},
		{
			input: [1]any{struct{}{}},
			error: "cannot export ([1]interface{})[0]: type struct {} is not supported",
			panic: "cannot export [1]interface {} to string: cannot export ([1]interface{})[0]: type struct {} is not supported",
		},
		{
			input:  []int{1, 2, 3, -1000000},
			output: "[]int{int(1), int(2), int(3), int(-1000000)}",
		},
		{
			input:  [3]int{1, 2, 3},
			output: "[3]int{int(1), int(2), int(3)}",
		},
		{
			input:  [0]int{},
			output: "[0]int{}",
		},
		{
			input:  ([]int)(nil),
			output: "([]int)(nil)",
		},
		{
			input:  []any{([]uint)(nil), []int{1, 2, 3}},
			output: "[]interface{}{([]uint)(nil), []int{int(1), int(2), int(3)}}",
		},
		{
			input:  [][]int{nil, {1, 2, 3}},
			output: `[][]int{([]int)(nil), []int{int(1), int(2), int(3)}}`,
		},
		{
			input:  []float32{},
			output: "make([]float32, 0)",
		},
		{
			input:  [0]float32{},
			output: "[0]float32{}",
		},
		{
			input: struct{}{},
			error: "type struct {} is not supported",
			panic: "cannot export struct {} to string: type struct {} is not supported",
		},
		{
			input: []interface{ Do() }{nil, nil, nil},
			error: "type []interface { Do() } is not supported",
			panic: "cannot export []interface { Do() } to string: type []interface { Do() } is not supported",
		},
		{
			input: [3]interface{ Do() }{},
			error: "type [3]interface { Do() } is not supported",
			panic: "cannot export [3]interface { Do() } to string: type [3]interface { Do() } is not supported",
		},
		{
			input:  []any{nil, nil, nil},
			output: "[]interface{}{nil, nil, nil}",
		},
		{
			input:  [3]any{},
			output: "[3]interface{}{nil, nil, nil}",
		},
		{
			input: []interface{ Do() }{nil},
			error: `type []interface { Do() } is not supported`,
			panic: `cannot export []interface { Do() } to string: type []interface { Do() } is not supported`,
		},
	}

	for i, s := range scenarios {
		s := s

		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			t.Parallel()

			func() {
				defer func() {
					r := recover()
					if s.panic == "" {
						assert.Nil(t, r)

						return
					}
					assert.Equal(t, s.panic, r)
				}()
				assert.Equal(t, s.output, MustExport(s.input))
			}()

			o, err := Export(s.input)
			if s.error == "" {
				require.NoError(t, err)
				assert.Equal(t, s.output, o)

				return
			}

			assert.EqualError(t, err, s.error)
		})
	}

	t.Run("Pointer loop", func(t *testing.T) {
		t.Parallel()

		a := make([]any, 2)
		a[1] = a
		v, err := Export(a)
		assert.EqualError(t, err, "cannot export ([]interface{})[1]: unexpected infinite loop")
		assert.Empty(t, v)
	})

	t.Run("Cannot convert []byte to string", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "[]uint8{uint8(255)}", MustExport([]byte{0xff}))
	})
}

//nolint:testifylint
func TestCastToString(t *testing.T) {
	t.Parallel()

	//nolint:exhaustruct
	scenarios := []struct {
		input  any
		output string
		error  string
	}{
		{
			input:  true,
			output: "true",
		},
		{
			input:  nil,
			output: "nil",
		},
		{
			input: struct{}{},
			error: "type struct {} is not supported",
		},
		{
			input:  "Bernhard Riemann",
			output: "Bernhard Riemann",
		},
		{
			input:  int(5),
			output: "5",
		},
		{
			input:  float64(3.14),
			output: "3.14",
		},
		{
			input:  int(10000000000),
			output: `10000000000`,
		},
		{
			input:  float64(10000000000),
			output: `10000000000`,
		},
		{
			input:  float32(10000000000),
			output: `10000000000`,
		},
	}

	for i, s := range scenarios {
		s := s

		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			t.Parallel()

			t.Run("CastToString", func(t *testing.T) {
				t.Parallel()

				result, err := CastToString(s.input)

				if s.error != "" {
					assert.Empty(t, result)
					assert.EqualError(t, err, s.error)

					return
				}

				assert.NoError(t, err)
				assert.Equal(t, s.output, result)
			})

			t.Run("MustCastToString", func(t *testing.T) {
				t.Parallel()

				defer func() {
					err := recover()
					if s.error == "" {
						assert.Nil(t, err)

						return
					}

					assert.NotNil(t, err)

					assert.Equal(
						t,
						fmt.Sprintf(
							"cannot cast %T to string: %s",
							s.input,
							s.error,
						),
						fmt.Sprintf("%s", err),
					)
				}()

				assert.Equal(t, s.output, MustCastToString(s.input))
			})
		})
	}
}

func TestNumericExporter_Supports(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		input    any
		expected bool
	}{
		{
			input:    nil,
			expected: false,
		},
		{
			input:    math.Pi,
			expected: true,
		},
		{
			input:    "3.14",
			expected: false,
		},
	}

	for i, s := range scenarios {
		s := s

		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			t.Parallel()

			assert.Equal(
				t,
				s.expected,
				numberExporter{}.supports(s.input), //nolint:exhaustruct
			)
		})
	}
}
