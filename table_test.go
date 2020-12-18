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
