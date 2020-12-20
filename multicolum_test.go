// -*- coding: utf-8 -*-
// multicolum_test.go
// -----------------------------------------------------------------------------
//
// Started on <dom 20-12-2020 00:55:55.426747268 (1608422155)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

//
// Description
//
package table

import (
	"reflect"
	"testing"
)

func Test_multicolumn_Process(t *testing.T) {
	type args struct {
		colspec   string
		nbcolumns int
		spec      string
		args      []interface{}
	}
	tests := []struct {
		name string
		args args
		want []formatter
	}{

		// Suppressing the first vertical separator of the table
		{args: args{colspec: "|l",
			nbcolumns: 1,
			spec:      "l",
			args:      []interface{}{"Black lives matter"}},
			want: []formatter{multicolumn{output: "Black lives matter"}}},

		{args: args{colspec: "|l",
			nbcolumns: 1,
			spec:      "c",
			args:      []interface{}{"Black lives matter"}},
			want: []formatter{multicolumn{output: "Black lives matter"}}},

		{args: args{colspec: "|l",
			nbcolumns: 1,
			spec:      "r",
			args:      []interface{}{"Black lives matter"}},
			want: []formatter{multicolumn{output: "Black lives matter"}}},

		// Preserving the first vertical separator of the table but changing the
		// horizontal alignment
		{args: args{colspec: "|l",
			nbcolumns: 1,
			spec:      "|l",
			args:      []interface{}{"Black lives matter"}},
			want: []formatter{multicolumn{output: "│Black lives matter"}}},

		{args: args{colspec: "|l",
			nbcolumns: 1,
			spec:      "|c",
			args:      []interface{}{"Black lives matter"}},
			want: []formatter{multicolumn{output: "│Black lives matter"}}},

		{args: args{colspec: "|l",
			nbcolumns: 1,
			spec:      "|r",
			args:      []interface{}{"Black lives matter"}},
			want: []formatter{multicolumn{output: "│Black lives matter"}}},

		// Splitting one column in two subcolumns with different horizontal alignments
		{args: args{colspec: "|l",
			nbcolumns: 2,
			spec:      "l l",
			args:      []interface{}{"Black lives matter", "They do!"}},
			want: []formatter{multicolumn{output: "Black lives matter They do!"}}},

		{args: args{colspec: "|l",
			nbcolumns: 2,
			spec:      "l c",
			args:      []interface{}{"Black lives matter", "They do!"}},
			want: []formatter{multicolumn{output: "Black lives matter They do!"}}},

		{args: args{colspec: "|l",
			nbcolumns: 2,
			spec:      "l r",
			args:      []interface{}{"Black lives matter", "They do!"}},
			want: []formatter{multicolumn{output: "Black lives matter They do!"}}},

		// Splitting one column in two subcolumns with two different lines
		{args: args{colspec: "|l",
			nbcolumns: 2,
			spec:      "l l",
			args:      []interface{}{"Black lives\nmatter", "They do!"}},
			want: []formatter{multicolumn{output: "Black lives They do!"},
				multicolumn{output: "matter              "}}},

		{args: args{colspec: "|l",
			nbcolumns: 2,
			spec:      "c l",
			args:      []interface{}{"Black lives\nmatter", "They do!"}},
			want: []formatter{multicolumn{output: "Black lives They do!"},
				multicolumn{output: "  matter            "}}},

		{args: args{colspec: "|l",
			nbcolumns: 2,
			spec:      "r l",
			args:      []interface{}{"Black lives\nmatter", "They do!"}},
			want: []formatter{multicolumn{output: "Black lives They do!"},
				multicolumn{output: "     matter         "}}},

		{args: args{colspec: "|l",
			nbcolumns: 2,
			spec:      "l | l",
			args:      []interface{}{"Black lives\nmatter", "They do!"}},
			want: []formatter{multicolumn{output: "Black lives │ They do!"},
				multicolumn{output: "matter      │         "}}},

		{args: args{colspec: "|l",
			nbcolumns: 2,
			spec:      "c | l",
			args:      []interface{}{"Black lives\nmatter", "They do!"}},
			want: []formatter{multicolumn{output: "Black lives │ They do!"},
				multicolumn{output: "  matter    │         "}}},

		{args: args{colspec: "|l",
			nbcolumns: 2,
			spec:      "r | l",
			args:      []interface{}{"Black lives\nmatter", "They do!"}},
			want: []formatter{multicolumn{output: "Black lives │ They do!"},
				multicolumn{output: "     matter │         "}}},
	}

	for _, tt := range tests {

		// Create a table with the given column specification
		if tab, ok := NewTable(tt.args.colspec); ok != nil {
			panic("It was not possible to create a table!")
		} else {

			t.Run(tt.name, func(t *testing.T) {

				// create a multicolumn with the given arguments
				column := Multicolumn(tt.args.nbcolumns, tt.args.spec, tt.args.args...)

				// and now make sure that processing the multicolum produces the
				// expected results
				if got := column.Process(tab, 0, 0); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("content.Process() = '%v', want '%v'", got, tt.want)
				}
			})
		}
	}
}
