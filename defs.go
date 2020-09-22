// This package implements means for drawing data in tabular form. It is
// strongly based on the definition of tables in LaTeX but extends its
// functionality in various ways.
package table

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

const none = 0

const horizontal_single = '\u2500' // ─
const horizontal_double = '\u2550' // ═
const horizontal_thick = '\u2501'  // ━

const vertical_single = '\u2502' // │
const vertical_double = '\u2551' // ║
const vertical_thick = '\u2503'  // ┃

// The splitter is defined as an association of four different runes: the west,
// east, north and south runes of the splitter:
//
//        north
// west  splitter  east
//        south
//
// this creates an association of ASCII/UTF-8 characters defined by hand below.
var splitterUTF8 = map[rune]map[rune]map[rune]map[rune]rune{

	none: {
		none: {
			none: {
				none:            none,
				vertical_single: none,
				vertical_double: none,
				vertical_thick:  none,
			},
			vertical_single: {
				none:            none,
				vertical_single: none,
				vertical_double: none,
				vertical_thick:  none,
			},
			vertical_double: {
				none:            none,
				vertical_single: none,
				vertical_double: none,
				vertical_thick:  none,
			},
			vertical_thick: {
				none:            none,
				vertical_single: none,
				vertical_double: none,
				vertical_thick:  none,
			},
		},
		horizontal_single: {
			none: {
				none:            none,
				vertical_single: '\u250c', // ┌
				vertical_double: '\u2553', // ╓
				vertical_thick:  '\u250e', // ┎
			},
			vertical_single: {
				none:            '\u2514', // └
				vertical_single: '\u251c', // ├
				vertical_double: '\u251f', // ┟: south double not supported!
				vertical_thick:  '\u251f', // ┟: south double not supported!
			},
			vertical_double: {
				none:            '\u2559', // ╙
				vertical_single: '\u251e', // ┞: north double not supported!
				vertical_double: '\u255f', // ╟
				vertical_thick:  '\u2520', // ┠: north double not supported!
			},
			vertical_thick: {
				none:            '\u2516', // ┖
				vertical_single: '\u251e', // ┞
				vertical_double: '\u2520', // ┠: south double not supported!
				vertical_thick:  '\u2520', // ┠
			},
		},
		horizontal_double: {
			none: {
				none:            none,
				vertical_single: '\u2552', // ╒
				vertical_double: '\u2554', // ╔
				vertical_thick:  '\u250f', // ┏: east double not supported!
			},
			vertical_single: {
				none:            '\u2558', // ╘
				vertical_single: '\u255e', // ╞
				vertical_double: '\u2522', // ┢: east/south double not supported!
				vertical_thick:  '\u2522', // ┢: east double not supported!
			},
			vertical_double: {
				none:            '\u255a', // ╚
				vertical_single: '\u2521', // ┡: east/north double not supported!
				vertical_double: '\u2560', // ╠
				vertical_thick:  '\u2523', // ┣: east/north double not supported!
			},
			vertical_thick: {
				none:            '\u2517', // ┗: easth double not supported!
				vertical_single: '\u2521', // ┡: east double not supported!
				vertical_double: '\u2523', // ┣: east/south double not supported!
				vertical_thick:  '\u2523', // ┣: east double not supported!
			},
		},
		horizontal_thick: {
			none: {
				none:            none,
				vertical_single: '\u250d', // ┍
				vertical_double: '\u250f', // ┏: south double not supported!
				vertical_thick:  '\u250f', // ┏: south double not supported!
			},
			vertical_single: {
				none:            '\u2515', // ┕
				vertical_single: '\u251d', // ┝
				vertical_double: '\u2522', // ┢: south double not supported!
				vertical_thick:  '\u2522', // ┢
			},
			vertical_double: {
				none:            '\u2517', // ┗: north double not supported
				vertical_single: '\u2521', // ┡: north double not supported!
				vertical_double: '\u2523', // ┣: north/south double not supported!
				vertical_thick:  '\u2523', // ┣: north double not supported!
			},
			vertical_thick: {
				none:            '\u2517', // ┗
				vertical_single: '\u2521', // ┡
				vertical_double: '\u2523', // ┣: south double not supported
				vertical_thick:  '\u2523', // ┣
			},
		},
	},

	horizontal_single: {
		none: {
			none: {
				none:            '\u2500', // ─
				vertical_single: '\u2510', // ┐
				vertical_double: '\u2556', // ╖
				vertical_thick:  '\u2512', // ┒
			},
			vertical_single: {
				none:            '\u2518', // ┘
				vertical_single: '\u2524', // ┤
				vertical_double: '\u2527', // ┧: south double not supported!
				vertical_thick:  '\u2527', // ┧: south double not supported!
			},
			vertical_double: {
				none:            '\u255c', // ╜
				vertical_single: '\u2526', // ┦: north double not supported!
				vertical_double: '\u2562', // ╢: north/south double not supported!
				vertical_thick:  '\u2528', // ┨: north double not supported!
			},
			vertical_thick: {
				none:            '\u251a', // ┚
				vertical_single: '\u2526', // ┦
				vertical_double: '\u2528', // ┨: south double not supported!
				vertical_thick:  '\u2528', // ┨
			},
		},
		horizontal_single: {
			none: {
				none:            '\u2500', // ─
				vertical_single: '\u252c', // ┬
				vertical_double: '\u2565', // ╥
				vertical_thick:  '\u2530', // ┰
			},
			vertical_single: {
				none:            '\u2534', // ┴
				vertical_single: '\u253c', // ┼
				vertical_double: '\u2541', // ╁: south double not supported!
				vertical_thick:  '\u2541', // ╁: south double not supported!
			},
			vertical_double: {
				none:            '\u2568', // ⢲: SHOULD NOT BE BRAILLE!!
				vertical_single: '\u2540', // ╀: north double not supported!
				vertical_double: '\u2542', // ╂: north/south double not supported!
				vertical_thick:  '\u2542', // ╂: north double not supported!
			},
			vertical_thick: {
				none:            '\u2538', // ┸
				vertical_single: '\u2540', // ╀
				vertical_double: '\u2542', // ╂: south double not supported!
				vertical_thick:  '\u2542', // ╂
			},
		},
		horizontal_double: {
			none: {
				none:            '\u2534', // ┴
				vertical_single: '\u253c', // ┼
				vertical_double: '\u2541', // ╁: south double not supported!
				vertical_thick:  '\u2541', // ╁: south double not supported!
			},
			vertical_single: {
				none:            '\u2536', // ┶: east double not supported!
				vertical_single: '\u253e', // ┾: east dobule not supported!
				vertical_double: '\u2546', // ╆: east/south double not supported!
				vertical_thick:  '\u2546', // ╆: east double not supported!
			},
			vertical_double: {
				none:            '\u253a', // ┺: easth double not supported!
				vertical_single: '\u2544', // ╄: east/north double not supported!
				vertical_double: '\u254a', // ╊: east/north/south double not supported!
				vertical_thick:  '\u254a', // ╊: east/north double not supported!
			},
			vertical_thick: {
				none:            '\u253a', // ┺: easth double not supported!
				vertical_single: '\u2544', // ╄: east double not supported!
				vertical_double: '\u254a', // ╊: east/south double not supported!
				vertical_thick:  '\u254a', // ╊: east double not supported!
			},
		},
		horizontal_thick: {
			none: {
				none:            '\u2534', // ┴
				vertical_single: '\u253c', // ┼
				vertical_double: '\u2541', // ╁: south double not supported!
				vertical_thick:  '\u2541', // ╁: south double not supported!
			},
			vertical_single: {
				none:            '\u2536', // ┶
				vertical_single: '\u253e', // ┾
				vertical_double: '\u2546', // ╆: south double not supported!
				vertical_thick:  '\u2546', // ╆
			},
			vertical_double: {
				none:            '\u2536', // ┶: north double not supported
				vertical_single: '\u2544', // ╄: north double not supported!
				vertical_double: '\u254a', // ╊: north/south double not supported!
				vertical_thick:  '\u254a', // ╊: north double not supported!
			},
			vertical_thick: {
				none:            '\u253a', // ┺
				vertical_single: '\u2544', // ╄
				vertical_double: '\u254a', // ╊: south double not supported
				vertical_thick:  '\u254a', // ╊
			},
		},
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
// consists then of a number of lines for displaying its contents. Horizontal
// separators are represented as rows of height 1
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
