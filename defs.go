// This package implements means for drawing data in tabular form. It is
// strongly based on the definition of tables in LaTeX but extends its
// functionality in various ways.
package table

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Separators
const HORIZONTAL_SINGLE = '\u2500' // ─
const HORIZONTAL_DOUBLE = '\u2550' // ═
const HORIZONTAL_THICK = '\u2501'  // ━

const VERTICAL_SINGLE = '\u2502' // │
const VERTICAL_DOUBLE = '\u2551' // ║
const VERTICAL_THICK = '\u2503'  // ┃

var SPLITTER = map[rune]map[rune]rune{
	VERTICAL_SINGLE: {
		HORIZONTAL_SINGLE: '\u253c', // ┼
		HORIZONTAL_DOUBLE: '\u256a', // ╪
		HORIZONTAL_THICK:  '\u253f', // ┿
	},
	VERTICAL_DOUBLE: {
		HORIZONTAL_SINGLE: '\u256b', // ╫
		HORIZONTAL_DOUBLE: '\u256c', // ╬
		HORIZONTAL_THICK:  '\u256b', // this combination does not exist!
	},
	VERTICAL_THICK: {
		HORIZONTAL_SINGLE: '\u2542', // ╂
		HORIZONTAL_DOUBLE: '\u254b', // this combination does not exist!
		HORIZONTAL_THICK:  '\u254b', // ╋
	},
}

// the following regexp is used to mach an entire column specification string
const specRegexAll = `^([^clrCLRp]*(c|l|r|C\{\d+\}|L\{\d+\}|R\{\d+\}|p\{\d+\}))+`

// and the following regexp is used to match the specification of a single
// column
const specRegex = `^[^clrCLRp]*(c|l|r|C\{\d+\}|L\{\d+\}|R\{\d+\}|p\{\d+\})`

// to extract the format of a single column the following regexp is used
const columnSpecRegex = `(c|l|r|C\{\d+\}|L\{\d+\}|R\{\d+\}|p\{\d+\})`

// in case a paragraph style is used, the following regexp serves to extract the
// numerical argument
const pRegex = `^(C\{\d+\}|L\{\d+\}|R\{\d+\}|p\{\d+\})$`

// to split strings using the newline as a separator
const newlineRegex = `\n`

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// Table is the main type provided by this package. In order to draw data in
// tabular form it is necessary first to create a table with table.NewTable.
// Once a Table has been created, it is then possible to use all services
// provided for them
//
// A table consists of a slice of columns, each one consisting of a separator
// and a content following immediately after, and a number of rows where the
// height of each row and its formatting style are stored. Note that the last
// separator (if any) is stored as a column without content. They also store the
// cells of the table as a bidimensional matrix of contents that can be
// processed and formatted.
type Table struct {
	columns []column
	rows    []row
	cells   [][]formatter
}

// columns do not store the contents, which are given instead in data rows. A
// column consists then of a vertical separator, a number of columns for
// displaying its contents, and the corresponding styles for showing its
// contents both horizontally and vertically.
type column struct {
	sep              string
	width            int
	hformat, vformat style
}

// rows do not store the contents, which are given instead in data rows. A row
// consists then of a number of lines for displaying its contents
type row struct {
	height int
}

// The style of a body specifies how to draw it and it is represented typically
// with a string and, additionally, with a numerical value in case a specific
// style (such as 'p') requires it
type style struct {
	alignment byte // either c, r, l, p, t, b or none which is represented with 0
	arg       int  // in case it is needed, e.g., for p
}

// Contents are simply strings to be shown on each cell
type content string

// Tables also provide the facility to draw horizontal rules which consist of
// the repetition of a specific rune
type hrule rune

// Tables can draw cells provided that they can be both processed and formatted:
//
// Processing a cell means splitting its contents across several (physical) rows
// so that it satisfies the format of the column where it has to be shown
//
// Formatting a cell implies adding blank characters to a physical line so that
// it satisfies the format of the column where it has to be shown
type formatter interface {
	Process(col column) []string
	Format(col column) string
}
