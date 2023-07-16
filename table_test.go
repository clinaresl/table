// -*- coding: utf-8 -*-
// table_test.go
// -----------------------------------------------------------------------------
//
// Started on <sáb 19-12-2020 22:45:26.735542876 (1608414326)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

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
		colspec, rowspec string
	}
	tests := []struct {
		name      string
		args      args
		wantTable *Table
		wantError error
	}{

		// invalid column specifications
		{args: args{"", ""},
			wantTable: &Table{},
			wantError: errors.New("invalid column specification")},

		{args: args{"|", ""},
			wantTable: &Table{},
			wantError: errors.New("invalid column specification")},

		{args: args{"c", "bb"},
			wantTable: &Table{},
			wantError: errors.New("The number of columns given in the row specification (2) must be less or equal than 1, the number of columns given in the column specification")},

		{args: args{"c", "x"},
			wantTable: &Table{},
			wantError: errors.New("'x' is an incorrect vertical format")},

		// correct column specifications

		// horizontal style / no vertical style
		{args: args{"c", ""},
			wantTable: &Table{columns: []column{{hformat: style{alignment: 'c'},
				vformat: style{alignment: 't'}}}},
			wantError: nil},

		{args: args{"l", ""},
			wantTable: &Table{columns: []column{{hformat: style{alignment: 'l'},
				vformat: style{alignment: 't'}}}},
			wantError: nil},

		{args: args{"r", ""},
			wantTable: &Table{columns: []column{{hformat: style{alignment: 'r'},
				vformat: style{alignment: 't'}}}},
			wantError: nil},

		{args: args{"p{10}", ""},
			wantTable: &Table{columns: []column{{hformat: style{alignment: 'p', arg: 10},
				vformat: style{alignment: 't'}}}},
			wantError: nil},

		{args: args{"|c", ""},
			wantTable: &Table{columns: []column{{sep: "│",
				hformat: style{alignment: 'c'},
				vformat: style{alignment: 't'}}}},
			wantError: nil},

		{args: args{"cc", ""},
			wantTable: &Table{columns: []column{{
				hformat: style{alignment: 'c'},
				vformat: style{alignment: 't'}},
				{hformat: style{alignment: 'c'},
					vformat: style{alignment: 't'}}}},
			wantError: nil},

		{args: args{"l l", ""},
			wantTable: &Table{columns: []column{{
				hformat: style{alignment: 'l'},
				vformat: style{alignment: 't'}},
				{sep: " ",
					hformat: style{alignment: 'l'},
					vformat: style{alignment: 't'}}}},
			wantError: nil},

		{args: args{"| c ||| c", ""},
			wantTable: &Table{columns: []column{{sep: "│ ",
				hformat: style{alignment: 'c'},
				vformat: style{alignment: 't'}},
				{sep: " ┃ ",
					hformat: style{alignment: 'c'},
					vformat: style{alignment: 't'}}}},
			wantError: nil},

		{args: args{"| c ||| l || p{4} ||| p{100}  rr|", ""},
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

		// horizontal style / vertical style
		{args: args{"c", "t"},
			wantTable: &Table{columns: []column{{hformat: style{alignment: 'c'},
				vformat: style{alignment: 't'}}}},
			wantError: nil},

		{args: args{"l", "c"},
			wantTable: &Table{columns: []column{{hformat: style{alignment: 'l'},
				vformat: style{alignment: 'c'}}}},
			wantError: nil},

		{args: args{"r", "b"},
			wantTable: &Table{columns: []column{{hformat: style{alignment: 'r'},
				vformat: style{alignment: 'b'}}}},
			wantError: nil},

		{args: args{"p{10}", "c"},
			wantTable: &Table{columns: []column{{hformat: style{alignment: 'p', arg: 10},
				vformat: style{alignment: 'c'}}}},
			wantError: nil},

		{args: args{"|c", "t"},
			wantTable: &Table{columns: []column{{sep: "│",
				hformat: style{alignment: 'c'},
				vformat: style{alignment: 't'}}}},
			wantError: nil},

		{args: args{"cc", "cb"},
			wantTable: &Table{columns: []column{{
				hformat: style{alignment: 'c'},
				vformat: style{alignment: 'c'}},
				{hformat: style{alignment: 'c'},
					vformat: style{alignment: 'b'}}}},
			wantError: nil},

		{args: args{"l l", "tb"},
			wantTable: &Table{columns: []column{{
				hformat: style{alignment: 'l'},
				vformat: style{alignment: 't'}},
				{sep: " ",
					hformat: style{alignment: 'l'},
					vformat: style{alignment: 'b'}}}},
			wantError: nil},

		{args: args{"| c ||| c", "ct"},
			wantTable: &Table{columns: []column{{sep: "│ ",
				hformat: style{alignment: 'c'},
				vformat: style{alignment: 'c'}},
				{sep: " ┃ ",
					hformat: style{alignment: 'c'},
					vformat: style{alignment: 't'}}}},
			wantError: nil},

		{args: args{"| c ||| l || p{4} ||| p{100}  rr|", "tcb"},
			wantTable: &Table{columns: []column{{sep: "│ ",
				hformat: style{alignment: 'c'},
				vformat: style{alignment: 't'}},
				{sep: " ┃ ",
					hformat: style{alignment: 'l'},
					vformat: style{alignment: 'c'}},
				{sep: " ║ ",
					hformat: style{alignment: 'p', arg: 4},
					vformat: style{alignment: 'b'}},
				{sep: " ┃ ",
					hformat: style{alignment: 'p', arg: 100},
					vformat: style{alignment: 't'}},
				{sep: "  ",
					hformat: style{alignment: 'r'},
					vformat: style{alignment: 't'}},
				{sep: "",
					hformat: style{alignment: 'r'},
					vformat: style{alignment: 't'}},
				{sep: "│",
					hformat: style{},
					vformat: style{}}}},
			wantError: nil},

		{args: args{"| c ||| l || p{4} ||| p{100}  rr|", "tcbbct"},
			wantTable: &Table{columns: []column{{sep: "│ ",
				hformat: style{alignment: 'c'},
				vformat: style{alignment: 't'}},
				{sep: " ┃ ",
					hformat: style{alignment: 'l'},
					vformat: style{alignment: 'c'}},
				{sep: " ║ ",
					hformat: style{alignment: 'p', arg: 4},
					vformat: style{alignment: 'b'}},
				{sep: " ┃ ",
					hformat: style{alignment: 'p', arg: 100},
					vformat: style{alignment: 'b'}},
				{sep: "  ",
					hformat: style{alignment: 'r'},
					vformat: style{alignment: 'c'}},
				{sep: "",
					hformat: style{alignment: 'r'},
					vformat: style{alignment: 't'}},
				{sep: "│",
					hformat: style{},
					vformat: style{}}}},
			wantError: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// prepare the args to create a new table. At least, the colspec
			// shall be given and optionally, if a rowspec is given it should be
			// added
			specs := make([]string, 2)
			specs[0] = tt.args.colspec
			if tt.args.rowspec != "" {
				specs[1] = tt.args.rowspec
			}

			// create a new table with the given arguments
			got, err := NewTable(specs...)

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
	// Output: Black lives matter
	//
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
	// Output:
	// Output:
	// │  Year  ║     Year     ┃     Year     │
	// │  1979  ║     2013     ┃     2018     │
	// │ Ariane ║     Gaia     ┃    Aeolus    │
	// │        ║ Proba Series ┃ Bepicolombo  │
	// │        ║    Swarm     ┃ Metop Series │
	//
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
	t.AddRow("...", "...")
	t.AddRow("testflag", "testing flags")
	t.AddRow("testfunc", "testing functions")
	fmt.Printf("Output:\n%v", t)
	// Output:
	// Output:
	// bug        start a bug report
	// build      compile packages and
	//            dependencies
	// clean      remove object files and
	//            cached files
	// doc        show documentation for
	//            package or symbol
	// env        print Go environment
	//            information
	// fix        update packages to use
	//            new APIs
	// ...        ...
	// testflag   testing flags
	// testfunc   testing functions
	//
}

// The following example shows how to add single rules to various parts of a
// table
func ExampleTable_3() {

	t, _ := NewTable("l | r ")
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
	// Output:
	// Output:
	// ━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━
	// Country       │    Population
	// ──────────────┼───────────────
	// China         │ 1,394,015,977
	// India         │ 1,326,093,247
	// United States │   329,877,505
	// Indonesia     │   267,026,366
	// Pakistan      │   233,500,636
	// Nigeria       │   214,028,302
	// ━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━
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

// The following table tests the vertical formatting and also ANSI color escape
// sequences both for the contents and the separators
func ExampleTable_16() {

	t, err := NewTable("\033[38;2;160;10;10m| c \033[38;2;10;160;10m| c \033[38;2;80;80;160m| c \033[38;2;160;80;40m|\033[0m", "cb")
	if err != nil {
		log.Fatalln(" NewTable: Fatal error!")
	}
	t.AddRow("\033[38;2;206;10;0mPlayer\033[0m", "\033[38;2;10;206;0mYear\033[0m", "\033[38;2;100;0;206mTournament\033[0m")
	t.AddSingleRule()
	t.AddRow("\033[38;5;206mRafa\033[0m \033[31;1;4mNadal\033[0m", "2010", "French Open\nWimbledon\nUS Open")
	t.AddSingleRule()
	t.AddRow("Roger Federer", "2007", "\033[38;2;255;82;197;48;2;155;106;0mAustralian Open\033[0m\nWimbledon\nUS Open")
	t.AddSingleRule()

	fmt.Printf("Output:\n%v", t)
	// Output:
}

// Tables are stringers and AddRow adds the output of a Sprintf operation. As a
// result, tables can be nested
func ExampleTable_17() {

	board1, _ := NewTable("||cccccccc||")
	board1.AddDoubleRule()
	board1.AddRow("\u265c", "\u265e", "\u265d", "\u265b", "\u265a", "\u265d", "", "\u265c")
	board1.AddRow("\u265f", "\u265f", "\u265f", "\u265f", "\u2592", "\u265f", "\u265f", "\u265f")
	board1.AddRow("", "\u2592", "", "\u2592", "", "\u265e", "", "\u2592")
	board1.AddRow("\u2592", "", "\u2592", "", "\u265f", "", "\u2592", "")
	board1.AddRow("", "\u2592", "", "\u2592", "\u2659", "\u2659", "", "\u2592")
	board1.AddRow("\u2592", "", "\u2658", "", "\u2592", "", "\u2592", "")
	board1.AddRow("\u2659", "\u2659", "\u2659", "\u2659", "", "\u2592", "\u2659", "\u2659")
	board1.AddRow("\u2656", "", "\u2657", "\u2655", "\u2654", "\u2657", "\u2658", "\u2656")
	board1.AddDoubleRule()

	board2, _ := NewTable("||cccccccc||")
	board2.AddDoubleRule()
	board2.AddRow("\u265c", "\u265e", "\u265d", "\u265b", "\u265a", "\u265d", "\u265e", "\u265c")
	board2.AddRow("\u265f", "\u265f", "\u265f", "", "\u265f", "\u265f", "\u265f", "\u265f")
	board2.AddRow("", "\u2592", "", "\u2592", "", "\u2592", "", "\u2592")
	board2.AddRow("\u2592", "", "\u2592", "", "\u2592", "", "\u2592", "")
	board2.AddRow("", "\u2592", "", "\u2659", "\u265f", "\u2592", "", "\u2592")
	board2.AddRow("\u2592", "", "\u2592", "", "\u2592", "\u2659", "\u2592", "")
	board2.AddRow("\u2659", "\u2659", "\u2659", "\u2592", "", "\u2592", "\u2659", "\u2659")
	board2.AddRow("\u2656", "\u2658", "\u2657", "\u2655", "\u2654", "\u2657", "\u2658", "\u2656")
	board2.AddDoubleRule()

	t, _ := NewTable("| c | c  c |", "cct")
	t.AddSingleRule()
	t.AddRow("ECO Code", "Moves", "Board")
	t.AddSingleRule()
	t.AddRow("C26 Vienna Game: Vienna Gambit", "1.e4 e5 2.♘c3 ♞6 3.f4", board1)
	t.AddRow("D00 Blackmar-Diemer Gambit: Gedult Gambit", "1.e4 d5 2.d4 exd4 3.f3", board2)
	t.AddSingleRule()

	fmt.Printf("Output:\n%v", t)
	// Output:
}

// Tables are stringers and AddRow adds the output of a Sprintf operation. As a
// result, tables can be also stacked
func ExampleTable_18() {

	example, _ := NewTable("|||c|||c|||c|||c|||c|||")
	example.AddDoubleRule(1, 2, 3, 4)
	example.AddRow("(1, 1)", "(1, 2)", "(1, 3)", "(1, 4)", "(1, 5)")
	example.AddDoubleRule(0, 1, 2, 3, 4, 5)
	example.AddRow("(2, 1)", "(2, 2)", "(2, 3)", "(2, 4)", "(2, 5)")
	example.AddDoubleRule(1, 2, 3, 4)
	example.AddRow("(3, 1)", "(3, 2)", "(3, 3)", "(3, 4)", "(3, 5)")
	example.AddDoubleRule(0, 1, 2, 3, 4, 5)

	code, _ := NewTable("l")
	code.AddRow("t, _ := NewTable(\"|||c|||c|||c|||c|||c|||\")")
	code.AddRow("example.AddDoubleRule(1, 2, 3, 4)")
	code.AddRow("example.AddRow(\"(1, 1)\", \"(1, 2)\", \"(1, 3)\", \"(1, 4)\", \"(1, 5)\")")
	code.AddRow("example.AddDoubleRule(0, 1, 2, 3, 4, 5)")
	code.AddRow("example.AddRow(\"(2, 1)\", \"(2, 2)\", \"(2, 3)\", \"(2, 4)\", \"(2, 5)\")")
	code.AddRow("example.AddDoubleRule(1, 2, 3, 4)")
	code.AddRow("example.AddRow(\"(3, 1)\", \"(3, 2)\", \"(3, 3)\", \"(3, 4)\", \"(3, 5)\")")
	code.AddRow("example.AddDoubleRule(0, 1, 2, 3, 4, 5)")

	t, _ := NewTable("c | l")
	t.AddRow(example, code)

	fmt.Printf("Output:\n%v", t)
	// Output:
}

// The following table tests multicolumns where the width of the table columns
// are enough to accommodate its contents and viceversa
func ExampleTable_19() {

	t, _ := NewTable("l | r ")
	t.AddThickRule()
	t.AddRow(Multicolumn(2, "c", "Demographics 2020"))
	t.AddRow("Country", "Population")
	t.AddSingleRule()
	t.AddRow("China", "1,394,015,977")
	t.AddRow("India", "1,326,093,247")
	t.AddRow("United States", "329,877,505")
	t.AddRow("Indonesia", "267,026,366")
	t.AddRow("Pakistan", "233,500,636")
	t.AddRow("Nigeria", "214,028,302")
	t.AddSingleRule()
	t.AddRow(Multicolumn(2, "l", "Source: https://www.worldometers.info/"))
	t.AddThickRule()
	fmt.Printf("Output:\n%v", t)
	// Output:
}

// The following table tests different multicolumns in the same row.
// Example taken from the TeX chapter of stack exchange:
//
// https://tex.stackexchange.com/questions/314025/making-stats-table-with-multicolumn-and-cline
func ExampleTable_20() {

	t, _ := NewTable("l c c || c c")
	t.AddRow(Multicolumn(5, "c", "Table 2: Overall Results"))
	t.AddThickRule()
	t.AddRow("", Multicolumn(2, "c", "Females"), Multicolumn(2, "c", "Males"))
	t.AddSingleRule(1, 5)
	t.AddRow("Treatment", "Mortality", "Mean\nPressure", "Mortality", "Mean\nPressure")
	t.AddSingleRule()
	t.AddRow("Placebo", 0.21, 163, 0.22, 164)
	t.AddRow("ACE Inhibitor", 0.13, 142, 0.15, 144)
	t.AddRow("Hydralazine", 0.17, 143, 0.16, 140)
	t.AddThickRule()
	t.AddRow(Multicolumn(5, "c", "Adapted from\nhttps://tex.stackexchange.com/questions/314025/making-stats-table-with-multicolumn-and-cline"))
	t.AddSingleRule()
	fmt.Printf("Output:\n%v", t)
	// Output:
}

// The next example shows how multicolumns can be used to modify the style of a
// separator or the horizontal alignment of any cell
func ExampleTable_21() {

	t, _ := NewTable("    r   l c")
	t.AddRow(Multicolumn(3, "    c", "♁ Earth"))
	t.AddThickRule()
	t.AddRow(Multicolumn(3, "    C{30}", "\033[37;3mEarth is the third planet from the Sun and the only astronomical object known to harbor life\033[0m"))
	t.AddSingleRule()
	t.AddRow(Multicolumn(1, "   |c", "Feature"),
		Multicolumn(1, "   c", "Measure"),
		Multicolumn(1, "c|", "Unit"))
	t.AddSingleRule()
	t.AddRow("Aphelion", 152100000, "km")
	t.AddRow("Perihelion", 147095000, "km")
	t.AddRow("Eccentricity", 0.0167086)
	t.AddRow("Orbital period", 365.256363004)
	t.AddRow("Semi-major axis", 149598023, "km")
	t.AddSingleRule()
	t.AddRow(Multicolumn(3, "   │c│", "\033[37;3mData provided by Wikipedia\033[0m"))
	t.AddSingleRule()
	fmt.Printf("Output:\n%v", t)
	// Output:
}

// A wide variety of multicolumns that group table columns
func ExampleTable_22() {

	t, _ := NewTable("|c||c|||c|c||c|||c|")
	t.AddSingleRule()
	t.AddRow("(1, 1)", "(1, 2)", "(1, 3)", "(1, 4)", "(1, 5)", "(1, 6)")
	t.AddDoubleRule()
	t.AddRow(Multicolumn(2, "|c", "(2, 1)"), "(2, 3)", "(2, 4)", "(2, 5)", "(2, 6)")
	t.AddThickRule()
	t.AddRow("(3, 1)", Multicolumn(2, "|c||", "(3, 2)"), "(3, 4)", "(3, 5)", "(3, 6)")
	t.AddSingleRule()
	t.AddRow("(4, 1)", "(4, 2)", Multicolumn(2, "|c", "(4, 3)"), "(4, 5)", "(4, 6)")
	t.AddDoubleRule()
	t.AddRow("(5, 1)", "(5, 2)", "(5, 3)", Multicolumn(2, "||c", "(5, 4)"), "(5, 6)")
	t.AddThickRule()
	t.AddRow("(6, 1)", "(6, 2)", "(6, 3)", "(6, 4)", Multicolumn(2, "|c", "(6, 5)"))
	t.AddSingleRule()
	t.AddRow(Multicolumn(2, "|c", "(7, 1)"), Multicolumn(2, "|c", "(7,2)"), "(7, 3)", "(7, 4)")
	t.AddDoubleRule()
	t.AddRow(Multicolumn(2, "|c", "(8, 1)"), "(8, 2)", Multicolumn(2, "||c", "(8, 3)"), "(8, 4)")
	t.AddThickRule()
	t.AddRow(Multicolumn(2, "|c", "(9, 1)"), "(9, 2)", "(9, 3)", Multicolumn(2, "|c", "(9, 4)"))
	t.AddSingleRule()
	t.AddRow("(10, 1)", Multicolumn(2, "|c||", "(10, 2)**"), Multicolumn(2, "c", "(10, 3)"), "(10, 4)")
	t.AddDoubleRule()
	t.AddRow("(11, 1)", Multicolumn(2, "|c||", "(11, 2)"), "(11, 3)", Multicolumn(2, "|c", "(11, 4)"))
	t.AddThickRule()
	t.AddRow("(12, 1)", "(12,2)", Multicolumn(2, "|c", "(12, 3)"), Multicolumn(2, "|c", "(12, 4)"))
	t.AddSingleRule()
	t.AddRow(Multicolumn(2, "|c", "(13, 1)"), Multicolumn(2, "|c", "(13, 2)"), Multicolumn(2, "|c", "(13, 3)"))
	t.AddDoubleRule()
	t.AddRow(Multicolumn(3, "|c||", "(14, 1)"), "(14, 2)", "(14, 3)", "(14, 4)")
	t.AddThickRule()
	t.AddRow("(15, 1)", Multicolumn(3, "|c", "(15, 2)"), "(15, 3)", "(15, 4)")
	t.AddSingleRule()
	t.AddRow("(16, 1)", "(16, 2)", Multicolumn(3, "|c", "(16, 3)"), "(16, 4)")
	t.AddDoubleRule()
	t.AddRow("(17, 1)", "(17, 2)", "(17, 4)", Multicolumn(3, "||c", "(17, 3)"))
	t.AddThickRule()
	t.AddRow(Multicolumn(3, "|c||", "(18, 1)**"), Multicolumn(2, "c", "(18, 2)"), "(18, 3)")
	t.AddSingleRule()
	t.AddRow(Multicolumn(3, "|c||", "(19, 1)"), "(19, 2)", Multicolumn(2, "|c", "(19, 3)"))
	t.AddDoubleRule()
	t.AddRow("(20, 1)", Multicolumn(3, "|c", "(20, 2)"), Multicolumn(2, "|c", "(20, 3)"))
	t.AddThickRule()
	t.AddRow(Multicolumn(2, "|c", "(21, 1)"), Multicolumn(3, "|c", "(21, 2)"), "(21, 3)")
	t.AddSingleRule()
	t.AddRow(Multicolumn(2, "|c", "(22, 1)"), "(22, 2)", Multicolumn(3, "||c", "(22, 3)"))
	t.AddDoubleRule()
	t.AddRow("(23, 1)", Multicolumn(2, "|c", "(23, 2)"), Multicolumn(3, "|c", "(23, 3)"))
	t.AddThickRule()
	t.AddRow(Multicolumn(3, "|c", "(24, 1)"), Multicolumn(3, "|c", "(24, 2)"))
	t.AddSingleRule()
	t.AddRow(Multicolumn(6, "|c", "Merging columns"))
	fmt.Printf("Output:\n%v", t)
	// Output:
}

// A wide variety of multicolumns that split table columns
func ExampleTable_23() {

	t, _ := NewTable("|c|c|c|")
	t.AddRow("Column #1", "Column #2", "Column #3")
	t.AddSingleRule()
	t.AddRow(Multicolumn(1, "|c||c|||", "(1, 1)", "(1, 2)"), "(1, 3)", "(1, 4)")
	t.AddDoubleRule()
	t.AddRow("(2, 1)", Multicolumn(1, "|c||c|||", "(2, 2)", "(2, 3)"), "(2, 4)")
	t.AddThickRule()
	t.AddRow("(3, 1)", "(3, 2)", Multicolumn(1, "|c||c|||", "(3, 3)", "(3, 4)"))
	t.AddSingleRule()
	t.AddRow(Multicolumn(1, "|c||c|c||", "(4, 1)", "(4, 2)", "(4, 3)"), "(4, 4)", "(4, 5)")
	t.AddDoubleRule()
	t.AddRow("(5, 1)", Multicolumn(1, "|c||c|c||", "(5, 2)", "(5, 3)", "(5, 4)"), "(5, 5)")
	t.AddThickRule()
	t.AddRow("(6, 1)", "(6, 2)", Multicolumn(1, "|c||c|c||", "(6, 3)", "(6, 4)", "(6, 5)"))
	t.AddSingleRule()
	t.AddRow(Multicolumn(3, "|c", "Splitting columns"))
	fmt.Printf("Output:\n%v", t)
	// Output:
}
