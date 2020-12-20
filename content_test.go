// -*- coding: utf-8 -*-
// content_test.go
// -----------------------------------------------------------------------------
//
// Started on <sáb 19-12-2020 22:45:26.735542876 (1608414326)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

package table

import (
	"reflect"
	"testing"
)

func Test_content_Process(t *testing.T) {
	type args struct {
		colspec string
		text    string
	}
	tests := []struct {
		name string
		args args
		want []formatter
	}{

		// All tests are performed over a table which consists of only one row
		// and one column.
		{args: args{colspec: "|l",
			text: ""},
			want: []formatter{content("")}},

		{args: args{colspec: "|c",
			text: ""},
			want: []formatter{content("")}},

		{args: args{colspec: "|r",
			text: ""},
			want: []formatter{content("")}},

		{args: args{colspec: "|l",
			text: "Black lives matter"},
			want: []formatter{content("Black lives matter")}},

		{args: args{colspec: "|c",
			text: "Black lives matter"},
			want: []formatter{content("Black lives matter")}},

		{args: args{colspec: "|r",
			text: "Black lives matter"},
			want: []formatter{content("Black lives matter")}},

		{args: args{colspec: "|l",
			text: "Black\nlives\nmatter"},
			want: []formatter{content("Black"), content("lives"), content("matter")}},

		{args: args{colspec: "|c",
			text: "Black\nlives\nmatter"},
			want: []formatter{content("Black"), content("lives"), content("matter")}},

		{args: args{colspec: "|r",
			text: "Black\nlives\nmatter"},
			want: []formatter{content("Black"), content("lives"), content("matter")}},

		{args: args{colspec: "|l",
			text: "Black\nlives\nmatter\n"},
			want: []formatter{content("Black"), content("lives"), content("matter"), content("")}},

		{args: args{colspec: "|c",
			text: "Black\nlives\nmatter\n"},
			want: []formatter{content("Black"), content("lives"), content("matter"), content("")}},

		{args: args{colspec: "|r",
			text: "Black\nlives\nmatter\n"},
			want: []formatter{content("Black"), content("lives"), content("matter"), content("")}},

		{args: args{colspec: "|l",
			text: "Black\n\nlives\nmatter"},
			want: []formatter{content("Black"), content(""), content("lives"), content("matter")}},

		{args: args{colspec: "|c",
			text: "Black\nlives\n\nmatter"},
			want: []formatter{content("Black"), content("lives"), content(""), content("matter")}},

		{args: args{colspec: "|r",
			text: "Black\nlives\nmatter\n\nAlways"},
			want: []formatter{content("Black"), content("lives"), content("matter"), content(""), content("Always")}},
	}

	for _, tt := range tests {

		// Create a table with the given column specification, and add the given
		// text in a single row
		if tab, ok := NewTable(tt.args.colspec); ok != nil {
			panic("It was not possible to create a table!")
		} else {
			tab.AddRow(tt.args.text)
			t.Run(tt.name, func(t *testing.T) {
				if got := tab.cells[0][0].Process(tab, 0, 0); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("content.Process() = '%v', want '%v'", got, tt.want)
				}
			})
		}
	}
}
