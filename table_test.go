package table

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"
)

// ----------------------------------------------------------------------------
// Tests
// ----------------------------------------------------------------------------

func TestNewTable(t *testing.T) {
	type args struct {
		spec string
	}
	tests := []struct {
		name      string
		args      args
		wantTable *Table
		wantError error
	}{

		// invalid column specifications
		{args: args{""},
			wantTable: &Table{},
			wantError: errors.New("invalid column specification")},

		{args: args{"|"},
			wantTable: &Table{},
			wantError: errors.New("invalid column specification")},

		// correct column specifications
		{args: args{"c"},
			wantTable: &Table{columns: []column{{hformat: style{alignment: 'c'},
				vformat: style{alignment: 't'}}}},
			wantError: nil},

		{args: args{"l"},
			wantTable: &Table{columns: []column{{hformat: style{alignment: 'l'},
				vformat: style{alignment: 't'}}}},
			wantError: nil},

		{args: args{"r"},
			wantTable: &Table{columns: []column{{hformat: style{alignment: 'r'},
				vformat: style{alignment: 't'}}}},
			wantError: nil},

		{args: args{"p{10}"},
			wantTable: &Table{columns: []column{{hformat: style{alignment: 'p', arg: 10},
				vformat: style{alignment: 't'}}}},
			wantError: nil},

		{args: args{"|c"},
			wantTable: &Table{columns: []column{{sep: "│",
				hformat: style{alignment: 'c'},
				vformat: style{alignment: 't'}}}},
			wantError: nil},

		{args: args{"cc"},
			wantTable: &Table{columns: []column{{
				hformat: style{alignment: 'c'},
				vformat: style{alignment: 't'}},
				{hformat: style{alignment: 'c'},
					vformat: style{alignment: 't'}}}},
			wantError: nil},

		{args: args{"l l"},
			wantTable: &Table{columns: []column{{
				hformat: style{alignment: 'l'},
				vformat: style{alignment: 't'}},
				{sep: " ",
					hformat: style{alignment: 'l'},
					vformat: style{alignment: 't'}}}},
			wantError: nil},

		{args: args{"| c ||| c"},
			wantTable: &Table{columns: []column{{sep: "│ ",
				hformat: style{alignment: 'c'},
				vformat: style{alignment: 't'}},
				{sep: " ┃ ",
					hformat: style{alignment: 'c'},
					vformat: style{alignment: 't'}}}},
			wantError: nil},

		{args: args{"| c ||| l || p{4} ||| p{100}  rr|"},
			wantTable: &Table{columns: []column{{sep: "│ ",
				hformat: style{alignment: 'c'},
				vformat: style{alignment: 't'}},
				{sep: " ┃ ",
					hformat: style{alignment: 'l'},
					vformat: style{alignment: 't'}},
				{sep: " ║ ",
					hformat: style{alignment: 'p', arg: 4},
					vformat: style{alignment: 't'}},
				{sep: " ┃ ",
					hformat: style{alignment: 'p', arg: 100},
					vformat: style{alignment: 't'}},
				{sep: "  ",
					hformat: style{alignment: 'r'},
					vformat: style{alignment: 't'}},
				{sep: "",
					hformat: style{alignment: 'r'},
					vformat: style{alignment: 't'}},
				{sep: "│"}}},
			wantError: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// create a new table
			got, err := NewTable(tt.args.spec)

			// verify that the new table is built as expected
			if !reflect.DeepEqual(got, tt.wantTable) {
				t.Errorf("[content] NewTable() = %v (want: %v), %v", got, tt.wantTable, err)
			}

			// verify also that the expected error is properly returned
			if (tt.wantError == nil && err != nil) ||
				(tt.wantError != nil && (err == nil || err.Error() != tt.wantError.Error())) {
				t.Errorf("[error] NewTable() = %v, %v (want: %v)", got, err, tt.wantError)
			}
		})
	}
}

func TestTable_GetNbColumns(t *testing.T) {
	type args struct {
		spec string
	}
	tests := []struct {
		name string
		args args
		want int
	}{

		{args: args{spec: ""},
			want: 0},

		{args: args{spec: "|"},
			want: 0},

		{args: args{spec: "c"},
			want: 1},

		{args: args{spec: "l"},
			want: 1},

		{args: args{spec: "r"},
			want: 1},

		{args: args{spec: "p{5}"},
			want: 1},

		{args: args{spec: "|c|"},
			want: 1},

		{args: args{spec: "|l|"},
			want: 1},

		{args: args{spec: "|r|"},
			want: 1},

		{args: args{spec: "|p{10}|"},
			want: 1},

		{args: args{spec: "cl"},
			want: 2},

		{args: args{spec: "lr"},
			want: 2},

		{args: args{spec: "rc"},
			want: 2},

		{args: args{spec: "p{7}l"},
			want: 2},

		{args: args{spec: "|||cl|"},
			want: 2},

		{args: args{spec: "||lr|"},
			want: 2},

		{args: args{spec: "rc|"},
			want: 2},

		{args: args{spec: "rp{20}|"},
			want: 2},

		{args: args{spec: "c|l"},
			want: 2},

		{args: args{spec: "l||r"},
			want: 2},

		{args: args{spec: "r|||c"},
			want: 2},

		{args: args{spec: "p{100}|||c"},
			want: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, _ := NewTable(tt.args.spec)
			if got := tr.GetNbColumns(); got != tt.want {
				t.Errorf("Table.GetNbColumns() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestTable_getFullSplitter(t *testing.T) {
// 	type args struct {
// 		irow  int
// 		jcol  int
// 		hrule rune
// 		sep   string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want string
// 	}{

// 		// horizontal single / vertical single

// 		// upper-left corner
// 		{args: args{irow: 0, jcol: 0, hrule: '─', sep: ""},
// 			want: ""},
// 		{args: args{irow: 0, jcol: 0, hrule: '─', sep: "│"},
// 			want: "┌"},

// 		// upper-mid edge
// 		{args: args{irow: 0, jcol: 1, hrule: '─', sep: "│"},
// 			want: "┬"},

// 		// upper-right edge
// 		{args: args{irow: 0, jcol: 2, hrule: '─', sep: "│"},
// 			want: "┐"},

// 		// right-mid edge
// 		{args: args{irow: 1, jcol: 0, hrule: '─', sep: "│"},
// 			want: "├"},

// 		// center
// 		{args: args{irow: 1, jcol: 1, hrule: '─', sep: "│"},
// 			want: "┼"},

// 		// left-mid edge
// 		{args: args{irow: 1, jcol: 2, hrule: '─', sep: "│"},
// 			want: "┤"},

// 		// bottom-left corner
// 		{args: args{irow: 2, jcol: 0, hrule: '─', sep: "│"},
// 			want: "└"},

// 		// bottom-mid edge
// 		{args: args{irow: 2, jcol: 1, hrule: '─', sep: "│"},
// 			want: "┴"},

// 		// bottom-right edge
// 		{args: args{irow: 2, jcol: 2, hrule: '─', sep: "│"},
// 			want: "┘"},

// 		// horizontal double / vertical single

// 		// upper-left corner
// 		{args: args{irow: 0, jcol: 0, hrule: '═', sep: ""},
// 			want: ""},
// 		{args: args{irow: 0, jcol: 0, hrule: '═', sep: "│"},
// 			want: "╒"},

// 		// upper-mid edge
// 		{args: args{irow: 0, jcol: 1, hrule: '═', sep: "│"},
// 			want: "╤"},

// 		// upper-right edge
// 		{args: args{irow: 0, jcol: 2, hrule: '═', sep: "│"},
// 			want: "╕"},

// 		// right-mid edge
// 		{args: args{irow: 1, jcol: 0, hrule: '═', sep: "│"},
// 			want: "╞"},

// 		// center
// 		{args: args{irow: 1, jcol: 1, hrule: '═', sep: "│"},
// 			want: "╪"},

// 		// left-mid edge
// 		{args: args{irow: 1, jcol: 2, hrule: '═', sep: "│"},
// 			want: "╡"},

// 		// bottom-left corner
// 		{args: args{irow: 2, jcol: 0, hrule: '═', sep: "│"},
// 			want: "╘"},

// 		// bottom-mid edge
// 		{args: args{irow: 2, jcol: 1, hrule: '═', sep: "│"},
// 			want: "\u2567"},

// 		// bottom-right edge
// 		{args: args{irow: 2, jcol: 2, hrule: '═', sep: "│"},
// 			want: "╛"},

// 		// horizontal single / vertical double

// 		// upper-left corner
// 		{args: args{irow: 0, jcol: 0, hrule: '─', sep: ""},
// 			want: ""},
// 		{args: args{irow: 0, jcol: 0, hrule: '─', sep: "║"},
// 			want: "╓"},

// 		// upper-mid edge
// 		{args: args{irow: 0, jcol: 1, hrule: '─', sep: "║"},
// 			want: "╥"},

// 		// upper-right edge
// 		{args: args{irow: 0, jcol: 2, hrule: '─', sep: "║"},
// 			want: "╖"},

// 		// right-mid edge
// 		{args: args{irow: 1, jcol: 0, hrule: '─', sep: "║"},
// 			want: "╟"},

// 		// center
// 		{args: args{irow: 1, jcol: 1, hrule: '─', sep: "║"},
// 			want: "╫"},

// 		// left-mid edge
// 		{args: args{irow: 1, jcol: 2, hrule: '─', sep: "║"},
// 			want: "╢"},

// 		// bottom-left corner
// 		{args: args{irow: 2, jcol: 0, hrule: '─', sep: "║"},
// 			want: "╙"},

// 		// bottom-mid edge
// 		{args: args{irow: 2, jcol: 1, hrule: '─', sep: "║"},
// 			want: "\u2568"},

// 		// bottom-right edge
// 		{args: args{irow: 2, jcol: 2, hrule: '─', sep: "║"},
// 			want: "╜"},

// 		// horizontal thick / vertical single

// 		// upper-left corner
// 		{args: args{irow: 0, jcol: 0, hrule: '━', sep: ""},
// 			want: ""},
// 		{args: args{irow: 0, jcol: 0, hrule: '━', sep: "│"},
// 			want: "┍"},

// 		// upper-mid edge
// 		{args: args{irow: 0, jcol: 1, hrule: '━', sep: "│"},
// 			want: "┯"},

// 		// upper-right edge
// 		{args: args{irow: 0, jcol: 2, hrule: '━', sep: "│"},
// 			want: "┑"},

// 		// right-mid edge
// 		{args: args{irow: 1, jcol: 0, hrule: '━', sep: "│"},
// 			want: "┝"},

// 		// center
// 		{args: args{irow: 1, jcol: 1, hrule: '━', sep: "│"},
// 			want: "┿"},

// 		// left-mid edge
// 		{args: args{irow: 1, jcol: 2, hrule: '━', sep: "│"},
// 			want: "┥"},

// 		// bottom-left corner
// 		{args: args{irow: 2, jcol: 0, hrule: '━', sep: "│"},
// 			want: "┕"},

// 		// bottom-mid edge
// 		{args: args{irow: 2, jcol: 1, hrule: '━', sep: "│"},
// 			want: "┷"},

// 		// bottom-right edge
// 		{args: args{irow: 2, jcol: 2, hrule: '━', sep: "│"},
// 			want: "┙"},

// 		// horizontal single / vertical thick

// 		// upper-left corner
// 		{args: args{irow: 0, jcol: 0, hrule: '─', sep: ""},
// 			want: ""},
// 		{args: args{irow: 0, jcol: 0, hrule: '─', sep: "┃"},
// 			want: "┎"},

// 		// upper-mid edge
// 		{args: args{irow: 0, jcol: 1, hrule: '─', sep: "┃"},
// 			want: "┰"},

// 		// upper-right edge
// 		{args: args{irow: 0, jcol: 2, hrule: '─', sep: "┃"},
// 			want: "┒"},

// 		// right-mid edge
// 		{args: args{irow: 1, jcol: 0, hrule: '─', sep: "┃"},
// 			want: "┠"},

// 		// center
// 		{args: args{irow: 1, jcol: 1, hrule: '─', sep: "┃"},
// 			want: "╂"},

// 		// left-mid edge
// 		{args: args{irow: 1, jcol: 2, hrule: '─', sep: "┃"},
// 			want: "┨"},

// 		// bottom-left corner
// 		{args: args{irow: 2, jcol: 0, hrule: '─', sep: "┃"},
// 			want: "┖"},

// 		// bottom-mid edge
// 		{args: args{irow: 2, jcol: 1, hrule: '─', sep: "┃"},
// 			want: "┸"},

// 		// bottom-right edge
// 		{args: args{irow: 2, jcol: 2, hrule: '─', sep: "┃"},
// 			want: "┚"},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			// Admittdely, to test these splitters, a 2x2 table is sufficient
// 			// but as we are not adding horizontal rules, a 3x2 table (with
// 			// vertical separators) is used instead. Importantly, the last
// 			// column has no text yet it has a separator
// 			tr, _ := NewTable("l|l|")
// 			tr.AddRow("cell (1,1)", "cell(1, 2)")
// 			tr.AddRow("cell (2,1)", "cell(2, 2)")
// 			tr.AddRow("cell (3,1)", "cell(3, 2)")
// 			if got := tr.getFullSplitter(tt.args.irow, tt.args.jcol, tt.args.hrule, tt.args.sep); got != tt.want {
// 				t.Errorf("Table.getFullSplitter() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// ----------------------------------------------------------------------------
// Examples
// ----------------------------------------------------------------------------

// The following example ilustrates the creation of a simple table with three
// columns and a single row
func ExampleTable_0() {

	t, err := NewTable("l l l")
	if err != nil {
		log.Fatalln(" NewTable: Fatal error!")
	}
	err = t.AddRow("Black", "lives", "matter")
	if err != nil {
		log.Fatalln(" AddRow: Fatal error!")
	}
	fmt.Printf("%v", t)
	// Output: ""
}

// In the next example, some rows expand over various lines. By default, these
// are vertically aligned to the top
func ExampleTable_1() {

	t, err := NewTable("| c || c ||| c |")
	if err != nil {
		log.Fatalln(" NewTable: Fatal error!")
	}
	err = t.AddRow("Year\n1979", "Year\n2013", "Year\n2018")
	if err != nil {
		log.Fatalln(" AddRow: Fatal error!")
	}
	err = t.AddRow("Ariane", "Gaia\nProba Series\nSwarm", "Aeolus\nBepicolombo\nMetop Series")
	if err != nil {
		log.Fatalln(" AddRow: Fatal error!")
	}
	fmt.Printf("Output:\n%v", t)
	// Output: ""
}

// This example shows how to use the package table to show information like in a
// help banner. In this case, the first column contains (some of) the commands
// of the go tool and the right one shows a comment about their usage. Note that
// to make sure that the entire table fits in the terminal, p is used as a
// column specifier
func ExampleTable_2() {

	t, err := NewTable("l   p{25}")
	if err != nil {
		log.Fatalln(" NewTable: Fatal error!")
	}
	t.AddRow("bug", "start a bug report")
	t.AddRow("build", "compile packages and dependencies")
	t.AddRow("clean", "remove object files and cached files")
	t.AddRow("doc", "show documentation for package or symbol")
	t.AddRow("env", "print Go environment information")
	t.AddRow("fix", "update packages to use new APIs")
	t.AddRow("")
	t.AddRow("...", "...")
	t.AddRow("")
	t.AddRow("testflag", "testing flags")
	t.AddRow("testfunc", "testing functions")
	fmt.Printf("Output:\n%v", t)
	// Output: ""
}

// The following example shows how to add single rules to various parts of a
// table
func ExampleTable_3() {

	t, _ := NewTable("> * l | r * <")
	t.AddThickRule()
	t.AddRow("Country", "Population")
	t.AddSingleRule()
	t.AddRow("China", "1,394,015,977")
	t.AddRow("India", "1,326,093,247")
	t.AddRow("United States", "329,877,505")
	t.AddRow("Indonesia", "267,026,366")
	t.AddRow("Pakistan", "233,500,636")
	t.AddRow("Nigeria", "214,028,302")
	t.AddThickRule()
	fmt.Printf("Output:\n%v", t)
	// Output: ""
}

// Horizontal rules can also be drawn from one specific column to another and it
// is possible to draw as many segments in the same line as required.
func ExampleTable_4() {

	t, _ := NewTable("|c|c|c|c|c|")
	t.AddSingleRule(0, 1, 2, 3, 4, 5)
	t.AddRow("(1, 1)", "(1, 2)", "(1, 3)", "(1, 4)", "(1, 5)")
	t.AddSingleRule(1, 2, 3, 4)
	t.AddRow("(2, 1)", "(2, 2)", "(2, 3)", "(2, 4)", "(2, 5)")
	t.AddSingleRule(0, 1, 2, 3, 4, 5)
	t.AddRow("(3, 1)", "(3, 2)", "(3, 3)", "(3, 4)", "(3, 5)")
	t.AddSingleRule(1, 2, 3, 4)
	fmt.Printf("Output:\n%v", t)
	// Output: ""
}

// This example is symmetrical to the previous one
func ExampleTable_5() {

	t, _ := NewTable("|c|c|c|c|c|")
	t.AddSingleRule(1, 2, 3, 4)
	t.AddRow("(1, 1)", "(1, 2)", "(1, 3)", "(1, 4)", "(1, 5)")
	t.AddSingleRule(0, 1, 2, 3, 4, 5)
	t.AddRow("(2, 1)", "(2, 2)", "(2, 3)", "(2, 4)", "(2, 5)")
	t.AddSingleRule(1, 2, 3, 4)
	t.AddRow("(3, 1)", "(3, 2)", "(3, 3)", "(3, 4)", "(3, 5)")
	t.AddSingleRule(0, 1, 2, 3, 4, 5)
	fmt.Printf("Output:\n%v", t)
	// Output: ""
}

// ----------------------------------------------------------------------------
func ExampleTable_6() {

	t, _ := NewTable("||c||c||c||c||c||")
	t.AddSingleRule(0, 1, 2, 3, 4, 5)
	t.AddRow("(1, 1)", "(1, 2)", "(1, 3)", "(1, 4)", "(1, 5)")
	t.AddSingleRule(1, 2, 3, 4)
	t.AddRow("(2, 1)", "(2, 2)", "(2, 3)", "(2, 4)", "(2, 5)")
	t.AddSingleRule(0, 1, 2, 3, 4, 5)
	t.AddRow("(3, 1)", "(3, 2)", "(3, 3)", "(3, 4)", "(3, 5)")
	t.AddSingleRule(1, 2, 3, 4)
	fmt.Printf("Output:\n%v", t)
	// Output: ""
}

func ExampleTable_7() {

	t, _ := NewTable("||c||c||c||c||c||")
	t.AddSingleRule(1, 2, 3, 4)
	t.AddRow("(1, 1)", "(1, 2)", "(1, 3)", "(1, 4)", "(1, 5)")
	t.AddSingleRule(0, 1, 2, 3, 4, 5)
	t.AddRow("(2, 1)", "(2, 2)", "(2, 3)", "(2, 4)", "(2, 5)")
	t.AddSingleRule(1, 2, 3, 4)
	t.AddRow("(3, 1)", "(3, 2)", "(3, 3)", "(3, 4)", "(3, 5)")
	t.AddSingleRule(0, 1, 2, 3, 4, 5)
	fmt.Printf("Output:\n%v", t)
	// Output: ""
}

// ----------------------------------------------------------------------------
func ExampleTable_8() {

	t, _ := NewTable("|||c|||c|||c|||c|||c|||")
	t.AddSingleRule(0, 1, 2, 3, 4, 5)
	t.AddRow("(1, 1)", "(1, 2)", "(1, 3)", "(1, 4)", "(1, 5)")
	t.AddSingleRule(1, 2, 3, 4)
	t.AddRow("(2, 1)", "(2, 2)", "(2, 3)", "(2, 4)", "(2, 5)")
	t.AddSingleRule(0, 1, 2, 3, 4, 5)
	t.AddRow("(3, 1)", "(3, 2)", "(3, 3)", "(3, 4)", "(3, 5)")
	t.AddSingleRule(1, 2, 3, 4)
	fmt.Printf("Output:\n%v", t)
	// Output: ""
}

func ExampleTable_9() {

	t, _ := NewTable("|||c|||c|||c|||c|||c|||")
	t.AddSingleRule(1, 2, 3, 4)
	t.AddRow("(1, 1)", "(1, 2)", "(1, 3)", "(1, 4)", "(1, 5)")
	t.AddSingleRule(0, 1, 2, 3, 4, 5)
	t.AddRow("(2, 1)", "(2, 2)", "(2, 3)", "(2, 4)", "(2, 5)")
	t.AddSingleRule(1, 2, 3, 4)
	t.AddRow("(3, 1)", "(3, 2)", "(3, 3)", "(3, 4)", "(3, 5)")
	t.AddSingleRule(0, 1, 2, 3, 4, 5)
	fmt.Printf("Output:\n%v", t)
	// Output: ""
}

// Horizontal rules can also be drawn from one specific column to another and it
// is possible to draw as many segments in the same line as required.
func ExampleTable_10() {

	t, _ := NewTable("|c|c|c|c|c|")
	t.AddDoubleRule(0, 1, 2, 3, 4, 5)
	t.AddRow("(1, 1)", "(1, 2)", "(1, 3)", "(1, 4)", "(1, 5)")
	t.AddDoubleRule(1, 2, 3, 4)
	t.AddRow("(2, 1)", "(2, 2)", "(2, 3)", "(2, 4)", "(2, 5)")
	t.AddDoubleRule(0, 1, 2, 3, 4, 5)
	t.AddRow("(3, 1)", "(3, 2)", "(3, 3)", "(3, 4)", "(3, 5)")
	t.AddDoubleRule(1, 2, 3, 4)
	fmt.Printf("Output:\n%v", t)
	// Output: ""
}

// This example is symmetrical to the previous one
func ExampleTable_11() {

	t, _ := NewTable("|c|c|c|c|c|")
	t.AddDoubleRule(1, 2, 3, 4)
	t.AddRow("(1, 1)", "(1, 2)", "(1, 3)", "(1, 4)", "(1, 5)")
	t.AddDoubleRule(0, 1, 2, 3, 4, 5)
	t.AddRow("(2, 1)", "(2, 2)", "(2, 3)", "(2, 4)", "(2, 5)")
	t.AddDoubleRule(1, 2, 3, 4)
	t.AddRow("(3, 1)", "(3, 2)", "(3, 3)", "(3, 4)", "(3, 5)")
	t.AddDoubleRule(0, 1, 2, 3, 4, 5)
	fmt.Printf("Output:\n%v", t)
	// Output: ""
}

// ----------------------------------------------------------------------------
func ExampleTable_12() {

	t, _ := NewTable("||c||c||c||c||c||")
	t.AddDoubleRule(0, 1, 2, 3, 4, 5)
	t.AddRow("(1, 1)", "(1, 2)", "(1, 3)", "(1, 4)", "(1, 5)")
	t.AddDoubleRule(1, 2, 3, 4)
	t.AddRow("(2, 1)", "(2, 2)", "(2, 3)", "(2, 4)", "(2, 5)")
	t.AddDoubleRule(0, 1, 2, 3, 4, 5)
	t.AddRow("(3, 1)", "(3, 2)", "(3, 3)", "(3, 4)", "(3, 5)")
	t.AddDoubleRule(1, 2, 3, 4)
	fmt.Printf("Output:\n%v", t)
	// Output: ""
}

func ExampleTable_13() {

	t, _ := NewTable("||c||c||c||c||c||")
	t.AddDoubleRule(1, 2, 3, 4)
	t.AddRow("(1, 1)", "(1, 2)", "(1, 3)", "(1, 4)", "(1, 5)")
	t.AddDoubleRule(0, 1, 2, 3, 4, 5)
	t.AddRow("(2, 1)", "(2, 2)", "(2, 3)", "(2, 4)", "(2, 5)")
	t.AddDoubleRule(1, 2, 3, 4)
	t.AddRow("(3, 1)", "(3, 2)", "(3, 3)", "(3, 4)", "(3, 5)")
	t.AddDoubleRule(0, 1, 2, 3, 4, 5)
	fmt.Printf("Output:\n%v", t)
	// Output: ""
}

// ----------------------------------------------------------------------------
func ExampleTable_14() {

	t, _ := NewTable("|||c|||c|||c|||c|||c|||")
	t.AddDoubleRule(0, 1, 2, 3, 4, 5)
	t.AddRow("(1, 1)", "(1, 2)", "(1, 3)", "(1, 4)", "(1, 5)")
	t.AddDoubleRule(1, 2, 3, 4)
	t.AddRow("(2, 1)", "(2, 2)", "(2, 3)", "(2, 4)", "(2, 5)")
	t.AddDoubleRule(0, 1, 2, 3, 4, 5)
	t.AddRow("(3, 1)", "(3, 2)", "(3, 3)", "(3, 4)", "(3, 5)")
	t.AddDoubleRule(1, 2, 3, 4)
	fmt.Printf("Output:\n%v", t)
	// Output: ""
}

func ExampleTable_15() {

	t, _ := NewTable("|||c|||c|||c|||c|||c|||")
	t.AddDoubleRule(1, 2, 3, 4)
	t.AddRow("(1, 1)", "(1, 2)", "(1, 3)", "(1, 4)", "(1, 5)")
	t.AddDoubleRule(0, 1, 2, 3, 4, 5)
	t.AddRow("(2, 1)", "(2, 2)", "(2, 3)", "(2, 4)", "(2, 5)")
	t.AddDoubleRule(1, 2, 3, 4)
	t.AddRow("(3, 1)", "(3, 2)", "(3, 3)", "(3, 4)", "(3, 5)")
	t.AddDoubleRule(0, 1, 2, 3, 4, 5)
	fmt.Printf("Output:\n%v", t)
	// Output: ""
}
