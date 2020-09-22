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

// NewTable creates a new table from the column specification. The column
// specification consists of a string which specifies the separator and style of
// each column, i.e., the horizontal alignment. By default, all rows are
// vertically aligned to the top.
//
// The different available alignments are given below:
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
// It returns a pointer to a table which can then be used to access its
// services. In case the column specification could not be processed it returns
// an error.
func NewTable(spec string) (*Table, error) {

	// -- initialization
	var columns []column

	// first things first, verify that the given column specification is
	// syntactically correct
	re := regexp.MustCompile(specRegexAll)
	if !re.MatchString(spec) {
		return &Table{}, errors.New("invalid column specification")
	}

	// the specification of the table is processed with a regular expression
	// which should be used to consume the whole string
	re = regexp.MustCompile(specRegex)
	for {

		// get the next column and, if none is found, then exit
		recol := re.FindStringIndex(spec)
		if recol == nil {
			break
		}
		nxtcol, err := newColumn(spec[recol[0]:recol[1]])
		if err != nil {
			return &Table{}, err
		}

		// and add it to the columns of this table
		columns = append(columns, *nxtcol)

		// and now move forward in the column specification string
		spec = spec[recol[1]:]
	}

	// maybe the column specification string is not empty here. Any remainings
	// are interpreted as the separator of a last column which contains no text
	if spec != "" {
		columns = append(columns,
			column{sep: spec,
				hformat: style{}})
	}

	// Before returning, process the separators of all columns to make the
	// appropriate substitutions
	for j := range columns {
		columns[j].sep = strings.ReplaceAll(columns[j].sep, "|||", "┃")
		columns[j].sep = strings.ReplaceAll(columns[j].sep, "||", "║")
		columns[j].sep = strings.ReplaceAll(columns[j].sep, "|", "│")
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
func (t *Table) getFullSplitter(irow, jcol int, hrule rune, sep string) (splitters string) {

	// define variables for storing the runes to the west, east, north and south
	// of each rune in the column separator
	var west, east, north, south rune

	// get all runes in the input column separator
	runes := getRunes(sep)

	// To do this it is mandatory to compute all runes to the west, east, north
	// and south in this specific order for each rune in the separator
	for _, irune := range runes {

		// west
		if jcol == 0 {
			west = none
		} else {
			west = hrule
		}

		// east. Note that a comparison is performed with the total number of
		// columns instead of using GetNbColumns because this fully ignores the
		// last column when it contains no text
		if jcol < len(t.columns)-1 {
			east = hrule
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

			// otherwise, add the horizontal rule
			splitters += string(hrule)
		} else {

			// if a correct separator was returned, then add it to the string to
			// return
			splitters += string(splitter)
		}
	}

	// and return the string computed so far
	return
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
		// rows required to show this line
		contents := icells[i].Process(t.columns[i])
		height = max(height, len(contents))

		// in addition update the number of physical columns required to draw
		// this cell
		for _, line := range contents {

			// if this column is a pagraph then use the width defined.
			// Otherwise, take the maximum width among all columns
			if t.columns[i].hformat.alignment == 'p' {
				t.columns[i].width = t.columns[i].hformat.arg
			} else {
				t.columns[i].width = max(t.columns[i].width, len(line))
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

// Add a single horizontal rule to the table
func (t *Table) AddSingleRule() {

	// rules are internally stored as a sequence of horizontal rules, one for
	// each column predefined in the table
	var icells []formatter
	for i := 0; i < t.GetNbColumns(); i++ {
		icells = append(icells, hrule(horizontal_single))
	}

	// and now add these rules to the contents to format and also an additional
	// row with a height always equal to one
	t.cells = append(t.cells, icells)
	t.rows = append(t.rows, row{height: 1})
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
		switch val := t.cells[i][0].(type) {

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
					result += fmt.Sprintf("%v", t.getFullSplitter(i, j, rune(val), t.columns[j].sep))

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

			// and each physical line of this row
			for line := 0; line < row.height; line++ {

				// and each column in this line, intentionally skipping the last one
				// in case it has no content
				for j := 0; j < t.GetNbColumns(); j++ {

					// Process the contents of this cell
					contents := t.cells[i][j].Process(t.columns[j])

					// get the text to show in this line. If this cell has no text
					// in the line-th line then format the empty string, otherwise,
					// use the text in contents
					text := ""
					if line < len(contents) {
						text = contents[line]
					}

					// and print the contents of this column in the result
					body := content(text).Format(t.columns[j])
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
