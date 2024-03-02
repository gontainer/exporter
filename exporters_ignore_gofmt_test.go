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

package exporter_test

import (
	"testing"

	"github.com/gontainer/exporter"
	"github.com/stretchr/testify/assert"
)

//nolint:testifylint
func TestExport_multidimensionalArray(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name   string
		input  any
		output string
	}{
		{
			name:   "[][2][][]int{{{{1, 2}}, nil}}",
			input:  [][2][][]int{[2][][]int{[][]int{[]int{int(1), int(2)}}, ([][]int)(nil)}},
			output: `[][2][][]int{[2][][]int{[][]int{[]int{int(1), int(2)}}, ([][]int)(nil)}}`,
		},
		{
			name:   "[]any{[][]int{{1, 2}, {3, 4}}, ([][][]any)(nil)}",
			input:  []interface{}{[][]int{[]int{int(1), int(2)}, []int{int(3), int(4)}}, ([][][]interface{})(nil)},
			output: `[]interface{}{[][]int{[]int{int(1), int(2)}, []int{int(3), int(4)}}, ([][][]interface{})(nil)}`,
		},
	}

	for _, s := range scenarios {
		s := s

		t.Run(s.name, func(t *testing.T) {
			t.Parallel()

			output, err := exporter.Export(s.input)
			assert.NoError(t, err)
			assert.Equal(t, s.output, output)
		})
	}
}
