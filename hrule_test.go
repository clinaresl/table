package table

import (
	"reflect"
	"testing"
)

func Test_hrule_Process(t *testing.T) {
	type args struct {
		colspec string
		rule    rune
		irow    int
		jcol    int
	}
	tests := []struct {
		name string
		args args
		want []formatter
	}{

		// Only the fully supported cases with full rules are tested. Other
		// weird combinations that include the double rule (either horizontal or
		// vertical) or the usage of clines is not tested

		// horizontal single / vertical single
		{args: args{colspec: "|c|c|", rule: horizontal_single, irow: 0, jcol: 0},
			want: []formatter{hrule("┌")}},

		// upper-mid edge
		{args: args{colspec: "|c|c|", rule: horizontal_single, irow: 0, jcol: 1},
			want: []formatter{hrule("┬")}},

		// upper-right corner
		{args: args{colspec: "|c|c|", rule: horizontal_single, irow: 0, jcol: 2},
			want: []formatter{hrule("┐")}},

		// right-mid edge
		{args: args{colspec: "|c|c|", rule: horizontal_single, irow: 2, jcol: 0},
			want: []formatter{hrule("├")}},

		// center
		{args: args{colspec: "|c|c|", rule: horizontal_single, irow: 2, jcol: 1},
			want: []formatter{hrule("┼")}},

		// left-mid edge
		{args: args{colspec: "|c|c|", rule: horizontal_single, irow: 2, jcol: 2},
			want: []formatter{hrule("┤")}},

		// bottom-left corner
		{args: args{colspec: "|c|c|", rule: horizontal_single, irow: 4, jcol: 0},
			want: []formatter{hrule("└")}},

		// bottom-mid edge
		{args: args{colspec: "|c|c|", rule: horizontal_single, irow: 4, jcol: 1},
			want: []formatter{hrule("┴")}},

		// bottom-right corner
		{args: args{colspec: "|c|c|", rule: horizontal_single, irow: 4, jcol: 2},
			want: []formatter{hrule("┘")}},

		// horizontal double / vertical single

		// upper-left corner
		{args: args{colspec: "|c|c|", rule: horizontal_double, irow: 0, jcol: 0},
			want: []formatter{hrule("╒")}},

		// upper-mid edge
		{args: args{colspec: "|c|c|", rule: horizontal_double, irow: 0, jcol: 1},
			want: []formatter{hrule("╤")}},

		// upper-right corner
		{args: args{colspec: "|c|c|", rule: horizontal_double, irow: 0, jcol: 2},
			want: []formatter{hrule("╕")}},

		// right-mid edge
		{args: args{colspec: "|c|c|", rule: horizontal_double, irow: 2, jcol: 0},
			want: []formatter{hrule("╞")}},

		// center
		{args: args{colspec: "|c|c|", rule: horizontal_double, irow: 2, jcol: 1},
			want: []formatter{hrule("╪")}},

		// left-mid edge
		{args: args{colspec: "|c|c|", rule: horizontal_double, irow: 2, jcol: 2},
			want: []formatter{hrule("╡")}},

		// bottom-left corner
		{args: args{colspec: "|c|c|", rule: horizontal_double, irow: 4, jcol: 0},
			want: []formatter{hrule("╘")}},

		// bottom-mid edge
		{args: args{colspec: "|c|c|", rule: horizontal_double, irow: 4, jcol: 1},
			want: []formatter{hrule("\u2567")}},

		// bottom-right corner
		{args: args{colspec: "|c|c|", rule: horizontal_double, irow: 4, jcol: 2},
			want: []formatter{hrule("╛")}},

		// horizontal thick / vertical single

		// upper-left corner
		{args: args{colspec: "|c|c|", rule: horizontal_thick, irow: 0, jcol: 0},
			want: []formatter{hrule("┍")}},

		// upper-mid edge
		{args: args{colspec: "|c|c|", rule: horizontal_thick, irow: 0, jcol: 1},
			want: []formatter{hrule("┯")}},

		// upper-right corner
		{args: args{colspec: "|c|c|", rule: horizontal_thick, irow: 0, jcol: 2},
			want: []formatter{hrule("┑")}},

		// right-mid edge
		{args: args{colspec: "|c|c|", rule: horizontal_thick, irow: 2, jcol: 0},
			want: []formatter{hrule("┝")}},

		// center
		{args: args{colspec: "|c|c|", rule: horizontal_thick, irow: 2, jcol: 1},
			want: []formatter{hrule("┿")}},

		// left-mid edge
		{args: args{colspec: "|c|c|", rule: horizontal_thick, irow: 2, jcol: 2},
			want: []formatter{hrule("┥")}},

		// bottom-left corner
		{args: args{colspec: "|c|c|", rule: horizontal_thick, irow: 4, jcol: 0},
			want: []formatter{hrule("┕")}},

		// bottom-mid edge
		{args: args{colspec: "|c|c|", rule: horizontal_thick, irow: 4, jcol: 1},
			want: []formatter{hrule("┷")}},

		// bottom-right corner
		{args: args{colspec: "|c|c|", rule: horizontal_thick, irow: 4, jcol: 2},
			want: []formatter{hrule("┙")}},

		// horizontal double / vertical single

		// upper-left corner
		{args: args{colspec: "|c|c|", rule: horizontal_double, irow: 0, jcol: 0},
			want: []formatter{hrule("╒")}},

		// upper-mid edge
		{args: args{colspec: "|c|c|", rule: horizontal_double, irow: 0, jcol: 1},
			want: []formatter{hrule("╤")}},

		// upper-right corner
		{args: args{colspec: "|c|c|", rule: horizontal_double, irow: 0, jcol: 2},
			want: []formatter{hrule("╕")}},

		// right-mid edge
		{args: args{colspec: "|c|c|", rule: horizontal_double, irow: 2, jcol: 0},
			want: []formatter{hrule("╞")}},

		// center
		{args: args{colspec: "|c|c|", rule: horizontal_double, irow: 2, jcol: 1},
			want: []formatter{hrule("╪")}},

		// left-mid edge
		{args: args{colspec: "|c|c|", rule: horizontal_double, irow: 2, jcol: 2},
			want: []formatter{hrule("╡")}},

		// bottom-left corner
		{args: args{colspec: "|c|c|", rule: horizontal_double, irow: 4, jcol: 0},
			want: []formatter{hrule("╘")}},

		// bottom-mid edge
		{args: args{colspec: "|c|c|", rule: horizontal_double, irow: 4, jcol: 1},
			want: []formatter{hrule("\u2567")}},

		// bottom-right corner
		{args: args{colspec: "|c|c|", rule: horizontal_double, irow: 4, jcol: 2},
			want: []formatter{hrule("╛")}},

		// horizontal double / vertical double

		// upper-left corner
		{args: args{colspec: "||c||c||", rule: horizontal_double, irow: 0, jcol: 0},
			want: []formatter{hrule("╔")}},

		// upper-mid edge
		{args: args{colspec: "||c||c||", rule: horizontal_double, irow: 0, jcol: 1},
			want: []formatter{hrule("╦")}},

		// upper-right corner
		{args: args{colspec: "||c||c||", rule: horizontal_double, irow: 0, jcol: 2},
			want: []formatter{hrule("╗")}},

		// right-mid edge
		{args: args{colspec: "||c||c||", rule: horizontal_double, irow: 2, jcol: 0},
			want: []formatter{hrule("╠")}},

		// center
		{args: args{colspec: "||c||c||", rule: horizontal_double, irow: 2, jcol: 1},
			want: []formatter{hrule("╬")}},

		// left-mid edge
		{args: args{colspec: "||c||c||", rule: horizontal_double, irow: 2, jcol: 2},
			want: []formatter{hrule("╣")}},

		// bottom-left corner
		{args: args{colspec: "||c||c||", rule: horizontal_double, irow: 4, jcol: 0},
			want: []formatter{hrule("╚")}},

		// bottom-mid edge
		{args: args{colspec: "||c||c||", rule: horizontal_double, irow: 4, jcol: 1},
			want: []formatter{hrule("╩")}},

		// bottom-right corner
		{args: args{colspec: "||c||c||", rule: horizontal_double, irow: 4, jcol: 2},
			want: []formatter{hrule("╝")}},

		// horizontal thick / vertical single

		// upper-left corner
		{args: args{colspec: "|c|c|", rule: horizontal_thick, irow: 0, jcol: 0},
			want: []formatter{hrule("┍")}},

		// upper-mid edge
		{args: args{colspec: "|c|c|", rule: horizontal_thick, irow: 0, jcol: 1},
			want: []formatter{hrule("┯")}},

		// upper-right corner
		{args: args{colspec: "|c|c|", rule: horizontal_thick, irow: 0, jcol: 2},
			want: []formatter{hrule("┑")}},

		// right-mid edge
		{args: args{colspec: "|c|c|", rule: horizontal_thick, irow: 2, jcol: 0},
			want: []formatter{hrule("┝")}},

		// center
		{args: args{colspec: "|c|c|", rule: horizontal_thick, irow: 2, jcol: 1},
			want: []formatter{hrule("┿")}},

		// left-mid edge
		{args: args{colspec: "|c|c|", rule: horizontal_thick, irow: 2, jcol: 2},
			want: []formatter{hrule("┥")}},

		// bottom-left corner
		{args: args{colspec: "|c|c|", rule: horizontal_thick, irow: 4, jcol: 0},
			want: []formatter{hrule("┕")}},

		// bottom-mid edge
		{args: args{colspec: "|c|c|", rule: horizontal_thick, irow: 4, jcol: 1},
			want: []formatter{hrule("┷")}},

		// bottom-right corner
		{args: args{colspec: "|c|c|", rule: horizontal_thick, irow: 4, jcol: 2},
			want: []formatter{hrule("┙")}},

		// horizontal thick / vertical thick

		// upper-left corner
		{args: args{colspec: "|||c|||c|||", rule: horizontal_thick, irow: 0, jcol: 0},
			want: []formatter{hrule("┏")}},

		// upper-mid edge
		{args: args{colspec: "|||c|||c|||", rule: horizontal_thick, irow: 0, jcol: 1},
			want: []formatter{hrule("┳")}},

		// upper-right corner
		{args: args{colspec: "|||c|||c|||", rule: horizontal_thick, irow: 0, jcol: 2},
			want: []formatter{hrule("┓")}},

		// right-mid edge
		{args: args{colspec: "|||c|||c|||", rule: horizontal_thick, irow: 2, jcol: 0},
			want: []formatter{hrule("┣")}},

		// center
		{args: args{colspec: "|||c|||c|||", rule: horizontal_thick, irow: 2, jcol: 1},
			want: []formatter{hrule("╋")}},

		// left-mid edge
		{args: args{colspec: "|||c|||c|||", rule: horizontal_thick, irow: 2, jcol: 2},
			want: []formatter{hrule("┫")}},

		// bottom-left corner
		{args: args{colspec: "|||c|||c|||", rule: horizontal_thick, irow: 4, jcol: 0},
			want: []formatter{hrule("┗")}},

		// bottom-mid edge
		{args: args{colspec: "|||c|||c|||", rule: horizontal_thick, irow: 4, jcol: 1},
			want: []formatter{hrule("┻")}},

		// bottom-right corner
		{args: args{colspec: "|||c|||c|||", rule: horizontal_thick, irow: 4, jcol: 2},
			want: []formatter{hrule("┛")}},
	}

	for _, tt := range tests {

		// Create a table with the given column specification, and add the given
		// text in a single row
		if tab, ok := NewTable(tt.args.colspec); ok != nil {
			panic("It was not possible to create a table!")
		} else {

			switch tt.args.rule {
			case horizontal_single:
				tab.AddSingleRule()
			case horizontal_double:
				tab.AddDoubleRule()
			case horizontal_thick:
				tab.AddThickRule()
			}
			tab.AddRow("cell (1,1)", "cell(1, 2)")
			switch tt.args.rule {
			case horizontal_single:
				tab.AddSingleRule()
			case horizontal_double:
				tab.AddDoubleRule()
			case horizontal_thick:
				tab.AddThickRule()
			}
			tab.AddRow("cell (2,1)", "cell(2, 2)")
			switch tt.args.rule {
			case horizontal_single:
				tab.AddSingleRule()
			case horizontal_double:
				tab.AddDoubleRule()
			case horizontal_thick:
				tab.AddThickRule()
			}

			t.Run(tt.name, func(t *testing.T) {
				if got := tab.cells[0][0].Process(tab, tt.args.irow, tt.args.jcol); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("hrule.Process() = %v, want %v", got, tt.want)
				}
			})
		}
	}
}
