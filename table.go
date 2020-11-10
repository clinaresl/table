package table

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ----------------------------------------------------------------------------
// Table
// ----------------------------------------------------------------------------

// Functions
// ----------------------------------------------------------------------------

// NewTable creates a new table from the column specification and optionally, a
// row specification.
//
// The column specification consists of a string which specifies the separator
// and style of each column, i.e., the horizontal alignment. The different
// available horizontal alignments are given below:
//
// 1. 'l': the contents of the column are ragged left
//
// 2. 'c': the contents of the column are horizontally centered
//
// 3. 'r': the contents of the column are ragged right
//
// 4. 'p{NUMBER}': the width of the column is fixed to NUMBER and the contents
// are split across various lines if needed
//
// 5. 'L{NUMBER}'/'C{NUMBER}'/'R{NUMBER}': the width of the column does not
// exceed NUMBER columns and the contents are ragged left/centered/ragged right
// respectively
//
// In addition, the column speciication might contain other characters which are
// then added to the contents as well.
//
// If a second string is given, then it is interpreted as the row specification.
// The different available vertical alignments are given below:
//
// 1. 't': the contents of the row are aligned to the top
//
// 2. 'c': the contents of the row are vertically centered
//
// 3. 'b': the contents of the row are aligned to the bottom
//
// By default, all rows are vertically aligned to the top, so that in case a row
// specification is given, then it has to refer to as many columns as the column
// specification given first or less ---the rest being aligned to the top.
// Otherwise, an error is returned. In contraposition to the column
// specification, the row specification can use only the specifiers given above
// and only those. Otherwise, an error is returned.
//
// It returns a pointer to a table which can then be used to access its
// services. In case either the column or row specification could not be
// processed it returns an error.
func NewTable(spec ...string) (*Table, error) {

	// error-checking
	if len(spec) == 0 || len(spec) > 2 {
		return &Table{}, errors.New("NewTable accepts only either one or two string arguments")
	}

	// -- initialization
	var columns []column
	var colspec, rowspec string

	// capture the args given by the user. Note that at this point spec contains
	// either one or two strings only
	if len(spec) >= 1 {
		colspec = spec[0]
	}
	if len(spec) == 2 {
		rowspec = spec[1]
	}

	// first things first, verify that the given column specification is
	// syntactically correct
	re := regexp.MustCompile(colSpecRegexAll)
	if !re.MatchString(colspec) {
		return &Table{}, errors.New("invalid column specification")
	}

	// the specification of the table is processed with a regular expression
	// which should be used to consume the whole string
	re = regexp.MustCompile(colSpecRegex)
	for {

		// get the next column and, if none is found, then exit
		recol := re.FindStringIndex(colspec)
		if recol == nil {
			break
		}
		nxtcol, err := newColumn(colspec[recol[0]:recol[1]])
		if err != nil {
			return &Table{}, err
		}

		// and add it to the columns of this table
		columns = append(columns, *nxtcol)

		// and now move forward in the column specification string
		colspec = colspec[recol[1]:]
	}

	// maybe the column specification string is not empty here. Any remainings
	// are interpreted as the separator of a last column which contains no text
	// and which have no format
	if colspec != "" {
		columns = append(columns,
			column{sep: colspec,
				hformat: style{},
				vformat: style{}})
	}

	// Before returning, process the separators of all columns to make the
	// appropriate substitutions
	for j := range columns {
		columns[j].sep = strings.ReplaceAll(columns[j].sep, "|||", "┃")
		columns[j].sep = strings.ReplaceAll(columns[j].sep, "||", "║")
		columns[j].sep = strings.ReplaceAll(columns[j].sep, "|", "│")
	}

	// Now, process all the vertical separators in case any has been given
	if vertFmt, err := getVerticalStyles(rowspec); err != nil {
		return &Table{}, err
	} else {

		// first, verify that the number of vertical specifiers is less or equal
		// than the number of columns given in the column specification, if any
		// was given
		if len(vertFmt) > 0 && len(vertFmt) > len(columns) {
			return &Table{}, fmt.Errorf("The number of columns given in the row specification (%v) must be less or equal than %v, the number of columns given in the column specification",
				len(vertFmt), len(columns))
		}
		for j, jstyle := range vertFmt {
			columns[j].vformat = jstyle
		}
	}

	// Note that the only information initialized in the creation of a table are
	// the columns
	return &Table{columns: columns}, nil
}

// Methods
// ----------------------------------------------------------------------------

// -- Private

// return the full sequence of splitters of a horizontal rule with a column
// separator. This method simply processes each rune of the separator taking
// into account other surrounding objects of the table. If one rune in the
// column separator has a corresponding splitter then a substitution is
// performed; otherwise, the given horizontal rule is used.
func (t *Table) getFullSplitter(irow, jcol int) (splitters string) {

	// define variables for storing the runes to the west, east, north and south
	// of each rune in the column separator
	var west, east, north, south rune

	// get the separator to process which is the one given in the j-th column
	sep := t.columns[jcol].sep

	// search for ANSI color escape sequences
	re := regexp.MustCompile(ansiColorRegex)
	colindexes := re.FindAllStringIndex(sep, -1)

	// locate at the first color and annotate how many have been found
	colind, nbcolors := 0, len(colindexes)

	// the following value should be equal to -1 if we have not found a vertical
	// separator yet and 0 otherwise. It is used to decide what horizontal rule
	// to use (either the one from the preceding the column or the one following
	// immediately after) in case a rune found in the separator has to be
	// substituted
	offset := -1

	// in addition, it is necessary to consider that the last column introduces
	// an asymmetry which comes from the fact that it might contain or not a
	// vertical separator (in spite of it containing data or not). If no
	// vertical separator is present, then runes in the separator should be
	// substituted by them if no splitter is found; otherwise, the horizontal
	// separator used in the preceding column should be used instead
	hasVerticalSeparator := containsVerticalSeparator(sep)

	// To do this it is mandatory to compute all runes to the west, east, north
	// and south in this specific order for each rune in the separator
	for idx, irune := range getRunes(sep) {

		// ANSI color escape sequences have to be directly copied to the
		// splitters
		if colind < nbcolors && idx >= colindexes[colind][0] {

			// if the ANSI color escape sequence starts right here then copy it
			// to the splitter
			if idx == colindexes[colind][0] {
				splitters += sep[colindexes[colind][0]:colindexes[colind][1]]
			}

			// if this position ends the entire ANSI color sequence, then move
			// to the next color
			if idx == colindexes[colind][1]-1 {
				colind += 1
			}

			// and skip the treatment of this rune (character)
			continue
		}

		// in case this rune is a vertical separator then make sure the offset
		// is 0. Note that if the string used in the separator contains more
		// than one vertical separator, the next column is considered
		// immediately after the first vertical separator, in spite of the
		// number of them
		if isVerticalSeparator(irune) && offset == -1 {
			offset = 0
		}

		// west
		if jcol == 0 {
			west = none
		} else {
			west = rune(t.cells[irow][jcol-1].(hrule))
			if west == horizontal_blank {
				west = none
			}
		}

		// east
		if jcol < t.GetNbColumns() {
			east = rune(t.cells[irow][jcol].(hrule))
			if east == horizontal_blank {
				east = none
			}
		} else {
			east = none
		}

		// north. The current separator is used as both the north and south
		// separator in case it does not fall out of bounds
		if irow > 0 {
			north = irune
		} else {
			north = none
		}
		if irow < t.GetNbRows()-1 {
			south = irune
		} else {
			south = none
		}

		// now, use the runes surrounding this one to access the splitter to use
		if splitter, err := getSingleSplitter(west, east, north, south); err != nil {

			// if it turns out that no splitter is known for this combination of
			// the west, east, north and south runes, then substitute it
			// properly. In general, the rune used instead is the one used in
			// the horizontal rule of either the preceding or next column
			// (depending whether a vertical separator has been processed yet).
			// Obviously, if the rune currently in process appears either before
			// or after the table then use it instead for the substitution
			if (offset == -1 && jcol == 0) ||
				(jcol == len(t.columns)-1 &&
					(!hasVerticalSeparator || offset == 0)) {
				splitters += string(irune)
			} else {
				splitters += string(rune(t.cells[irow][jcol+offset].(hrule)))
			}
		} else {

			// if a correct separator was returned, then add it to the string to
			// return
			splitters += string(splitter)
		}
	}

	// and return the string computed so far
	return
}

// Add the given horizontal rule to the table from a start column to and end
// column. Any number of pairs (start, end) can be given. If no column is given,
// the horizontal rule takes the entire width of the table. If the end column of
// a pair goes beyond the width of the table, the horizontal rule is drawn up to
// the last column
//
// In case it is not possible to process the given specification an informative
// error is returned
func (t *Table) addRule(rule hrule, cols ...int) error {

	// obviously, the number of columns should be an even number
	if len(cols)%2 != 0 {
		return errors.New("An even number of columns must be given when adding a horizontal rule")
	}

	// if no columns are given, then take the entire width of the table
	if len(cols) == 0 {
		cols = []int{0, len(t.columns)}
	}

	// rules are internally stored as a sequence of horizontal rules, one for
	// each column predefined in the table. First, make the horizontal rule to
	// contain the blank
	var icells []formatter
	for i := 0; i < t.GetNbColumns(); i++ {
		icells = append(icells, hrule(horizontal_blank))
	}

	// next, overwrite only those columns that are within a pair (start, end)
	for i := 0; i < len(cols); i += 2 {

		// if the end column is strictly less than the start column then return
		// an error
		if cols[i] > cols[i+1] {
			return errors.New("The end column of a horizontal rule should be strictly larger or equal than the start column")
		}

		// and now update the horizontal rule of the columns between the start and end only
		for j := cols[i]; j < cols[i+1] && j < t.GetNbColumns(); j++ {
			icells[j] = rule
		}
	}

	// and now add these rules to the contents to format and also an additional
	// row with a height always equal to one
	t.cells = append(t.cells, icells)
	t.rows = append(t.rows, row{height: 1})

	// and return no error
	return nil
}

// -- Public

// Add a new line of data to the bottom of the column. This function accepts an
// arbitrary number of arguments that satisfy the null interface. If the number
// of elements given exceeds the number of columns of the table an error is
// immediately issued
func (t *Table) AddRow(cells ...interface{}) error {

	// if the number of elements given exceeds the number of columns then
	// immediately raised an error
	if t.GetNbColumns() < len(cells) {
		return fmt.Errorf("The number of elements given (%v) exceeds the number of columns (%v)",
			len(cells), t.GetNbColumns())
	}

	// otherwise, process all cells given and add them to the table as cells
	// which can be formatted. 'i' is the column index and height is
	// the number of physical rows required to draw this row
	var i, height int
	icells := make([]formatter, t.GetNbColumns())
	for ; i < t.GetNbColumns() && i < len(cells); i++ {

		// add the content of the i-th column with a string that represents it
		text := fmt.Sprintf("%v", cells[i])
		icells[i] = content(text)

		// process the contents of this cell, and update the number of physical
		// rows required to show this line. As the number of physical rows is
		// unknown at this stage, a value equal to zero is given.
		contents := icells[i].Process(t.columns[i], 0)
		height = max(height, len(contents))

		// in addition update the number of physical columns required to draw
		// this cell
		for _, line := range contents {

			// if this column is a pagraph then use the width defined.
			// Otherwise, take the maximum width among all columns
			if t.columns[i].hformat.alignment == 'p' {
				t.columns[i].width = t.columns[i].hformat.arg
			} else {
				t.columns[i].width = max(t.columns[i].width, countPrintableRuneInString(line))
			}
		}
	}

	// now, if not all columns were given then automatically add empty cells
	for ; i < t.GetNbColumns(); i++ {
		icells[i] = content("")
	}

	// add this cells to this table, along with the number of physical rows
	// required to draw it
	t.cells = append(t.cells, icells)
	t.rows = append(t.rows, row{height: height})

	// and exit with no error
	return nil
}

// Add a single horizontal rule to the table from a start column to and end
// column. Any number of pairs (start, end) can be given. If no column is given,
// the horizontal rule takes the entire width of the table.
//
// In case it is not possible to process the given specification an informative
// error is returned
func (t *Table) AddSingleRule(cols ...int) error {

	return t.addRule(horizontal_single, cols...)
}

// Add a double horizontal rule to the table from a start column to and end
// column. Any number of pairs (start, end) can be given. If no column is given,
// the horizontal rule takes the entire width of the table.
//
// In case it is not possible to process the given specification an informative
// error is returned
func (t *Table) AddDoubleRule(cols ...int) error {

	return t.addRule(horizontal_double, cols...)
}

// Add a thick horizontal rule to the table from a start column to and end
// column. Any number of pairs (start, end) can be given. If no column is given,
// the horizontal rule takes the entire width of the table.
//
// In case it is not possible to process the given specification an informative
// error is returned
func (t *Table) AddThickRule(cols ...int) error {

	return t.addRule(horizontal_thick, cols...)
}

// Return the number of columns in a table which contain data.
func (t *Table) GetNbColumns() int {

	// if this table comes with a last column with no contents, then discard it
	if len(t.columns) > 0 && t.columns[len(t.columns)-1].hformat.alignment == 0 {
		return len(t.columns) - 1
	}

	// otherwise, return the number of columns defined
	return len(t.columns)
}

// Return the number of logical rows in a table, i.e., the number of rows given
// by the user. The number of logical rows includes both horizontal separators
// and data lines
func (t *Table) GetNbRows() int {
	return len(t.cells)
}

// Tables are stringers and thus they provide a method to conveniently transform
// its contents into a string
func (t Table) String() (result string) {

	// for each logical row (including the horizontal rules)
	for i, row := range t.rows {

		// the type of an entire logical row is equal to the type of its first
		// item. Thus, to distinguish rules from contents the first element of
		// this row is checked
		switch t.cells[i][0].(type) {

		case hrule:

			// Yes, I know!! Horizontal rules only take one line but ... who
			// knows if this changes in the future?
			for line := 0; line < row.height; line++ {

				// now, draw the horizontal rule for each column. Note that the
				// true number of columns is used to take into account the last
				// one even if it has no content
				for j := 0; j < len(t.columns); j++ {

					// add to the string the splitters of this column and next
					// the horizontal rule as stored in the table
					result += fmt.Sprintf("%v", t.getFullSplitter(i, j))

					// in case this is column has a horizontal rule attached show it as well
					if j < t.GetNbColumns() {

						// this is done by repeating the horizontal rule as many
						// times as the width of this column
						for k := 0; k < t.columns[j].width; k++ {
							result += fmt.Sprintf("%c", t.cells[i][j])
						}
					}
				}
				result += "\n"
			}

		case content:

			// for each physical line of this row
			for line := 0; line < row.height; line++ {

				// and each column in this line, intentionally skipping the last one
				// in case it has no content
				for j := 0; j < t.GetNbColumns(); j++ {

					// Process the contents of this cell. Note that at this
					// point the number of physical rows is known and thus it is
					// given to the Processor so that the contents are
					// vertically formatted as required
					contents := t.cells[i][j].Process(t.columns[j], row.height)

					// and print the contents of this column in the result
					body := content(contents[line]).Format(t.columns[j])
					result += fmt.Sprintf("%v%v", t.columns[j].sep, body)
				}

				// in case a last column is given with no format, then add the
				// separator, otherwise, just simply add the newline
				if t.columns[len(t.columns)-1].hformat.alignment == 0 {
					result += fmt.Sprintf("%v\n", t.columns[len(t.columns)-1].sep)
				} else {
					result += "\n"
				}
			}
		}
	}

	// and return the string computed so far
	return
}
