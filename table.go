// -*- coding: utf-8 -*-
// table.go
// -----------------------------------------------------------------------------
//
// Started on <sáb 19-12-2020 22:45:26.735542876 (1608414326)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

// This package provides means for drawing data in tabular form. It is
// strongly based on the definition of tables in LaTeX but extends its
// functionality in various ways through a very simple interface
//
// It honours UTF-8 characters, ANSI color escape sequences, full/partial
// horizontal rules, and a wide variety of vertical and horizontal alignment
// options.
//
// Remarkably, it prints any stringer and as tables are stringers, tables can be
// nested to any degree.
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
// 1. 'l': the contents of the column are ragged right
//
// 2. 'c': the contents of the column are horizontally centered
//
// 3. 'r': the contents of the column are ragged left
//
//  4. 'p{NUMBER}': the width of the column is fixed to NUMBER positions and the
//     contents are split across various lines if needed
//
//  5. 'L{NUMBER}'/'C{NUMBER}'/'R{NUMBER}': the width of the column does not
//     exceed NUMBER columns and the contents are ragged left/centered/ragged
//     right respectively
//
// In addition, the column specification might contain other characters which are
// then added to the contents as well.
//
// If a second string is given, then it is interpreted as the row specification,
// which specifies the vertical alignment:
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
// NewTable returns a pointer to a table which can then be used to access its
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
	colspec = spec[0]
	if len(spec) == 2 {
		rowspec = spec[1]
	}

	// first things first, verify that the given column specification is
	// syntactically correct
	re := regexp.MustCompile(colSpecRegexAll)
	if !re.MatchString(colspec) {
		return &Table{}, errors.New("invalid column specification")
	}

	// Now, make the appropriate substitutions in case the user used ANSI
	// characters for specifying vertical delimiters
	colspec = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(colspec,
		"|||", "┃"),
		"||", "║"),
		"|", "│")

	// process the column specification given
	columns, err := getColumns(colspec)
	if err != nil {
		return &Table{}, err
	}

	// Before returning, process the separators of all columns to make the
	// appropriate substitutions
	for j := range columns {
		separatorToUTF8(&columns[j].sep)
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
	// contain the blank.
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

	// in case this table contains a last column with no data, add an empty rule
	if len(t.columns) > 0 && t.columns[len(t.columns)-1].hformat.alignment == 0 {
		icells = append(icells, hrule(horizontal_empty))
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
// arbitrary number of arguments. The content shown on the table is the output
// of a Sprintf operation over each argument
//
// If the number of arguments is less than the number of columns, the last cells
// are left empty, unless no argument is given at all in which case no row is
// inserted. Finally, if the number of elements given exceeds the number of
// columns an error is immediately issued
func (t *Table) AddRow(cells ...interface{}) error {

	// if the number of elements given exceeds the number of columns then
	// immediately raised an error
	if t.GetNbColumns() < len(cells) {
		return fmt.Errorf("The number of elements given (%v) exceeds the number of columns (%v)",
			len(cells), t.GetNbColumns())
	}

	// otherwise, process all cells given and add them to the table as cells
	// which can be formatted. 'j' is the column index and height is the number
	// of physical rows required to draw this row, whereas idx is the index of
	// the next cell to process used in the loop that iterates over them
	var j, height int
	icells := make([]formatter, len(t.columns))
	for idx := 0; idx < t.GetNbColumns() && idx < len(cells); idx++ {

		// add the content of the i-th column with a string that represents it
		icells[j] = content(fmt.Sprintf("%v", cells[idx]))

		// process the contents of this cell, and update the number of physical
		// rows required to show this line. Note that this row is added to the
		// bottom of the table
		contents := icells[j].Process(t, len(t.rows), j)
		height = max[int](height, len(contents))

		// in addition update the number of physical columns required to
		// draw this cell.
		//
		// If this column is a paragraph (of any type) then use the width
		// defined
		if t.columns[j].hformat.alignment == 'p' ||
			t.columns[j].hformat.alignment == 'C' ||
			t.columns[j].hformat.alignment == 'L' ||
			t.columns[j].hformat.alignment == 'R' {

			// Importantly, the width of this column should be modified if and
			// only if it is less than the argument given. The reason is that
			// the width of this column might have been increased (e.g., because
			// it is part of a multicell) so that it should not be modified
			// now!!
			if t.columns[j].width < t.columns[j].hformat.arg {
				t.columns[j].width = t.columns[j].hformat.arg
			}
		} else {

			// Otherwise, take the maximum width among all columns
			for _, line := range contents {
				t.columns[j].width = max[int](t.columns[j].width,
					countPrintableRuneInString(string(line.(content))))
			}
		}

		// and move to the next column
		j++
	}
	// now, if not all columns were given then automatically add empty cells.
	// Note that an empty cell is added also to the last column even if it
	// contains no data
	for ; j < len(t.columns); j++ {
		icells[j] = content(horizontal_empty)
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

	return t.addRule(hrule(horizontal_single), cols...)
}

// Add a double horizontal rule to the table from a start column to and end
// column. Any number of pairs (start, end) can be given. If no column is given,
// the horizontal rule takes the entire width of the table.
//
// In case it is not possible to process the given specification an informative
// error is returned
func (t *Table) AddDoubleRule(cols ...int) error {

	return t.addRule(hrule(horizontal_double), cols...)
}

// Add a thick horizontal rule to the table from a start column to and end
// column. Any number of pairs (start, end) can be given. If no column is given,
// the horizontal rule takes the entire width of the table.
//
// In case it is not possible to process the given specification an informative
// error is returned
func (t *Table) AddThickRule(cols ...int) error {

	return t.addRule(hrule(horizontal_thick), cols...)
}

// Return the number of logical columns in a table which contain data.
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
// their contents into a string
func (t Table) String() string {

	// Because of the presence of multicells, each line can print at once an
	// arbitrary number of columns and rows. Hence, it is required to keep track
	// of how many columns are printed in each row and how many rows are printed
	// in each column
	nbcolumns := make([]int, len(t.rows))
	nbrows := make([]int, len(t.columns))

	// Tables are formatted iterating over columns and thus the content of each
	// row is stored separately in a slice of strings which are then
	// concatenated
	var output []string

	// for each logical column
	for j := 0; j < len(t.columns); j++ {

		// initialize the index of the output string that should get the
		// contents of the next physical line
		idx := 0

		// for each logical row
		for i, row := range t.rows {

			// Make sure to skip those columns that have been printed before as
			// a result of a multicell
			if nbcolumns[i] <= j && nbrows[j] <= i {

				// Process the contents of this cell
				contents := t.cells[i][j].Process(&t, i, j)

				// and now for each physical row of this line
				for line := 0; line < row.height; line++ {

					// add this line to the output
					if idx >= len(output) {
						output = append(output, fmt.Sprintf("%v", contents[line].Format(&t, i, j)))
					} else {
						output[idx] += fmt.Sprintf("%v", contents[line].Format(&t, i, j))
					}

					// and move to the next physical line
					idx++
				}
			} else {

				// else, this position has been generated before, so that all
				// that is left is to move to the next physical line
				idx += row.height
			}

			// Now, update the number of columns processed in this logical row
			// and the number of rows processed in this column

			// If this is not a multicell, the next column to process in
			// this row should be updated only in case that we already
			// reached the last one
			if nbcolumns[i] <= j {
				nbcolumns[i]++
			}

			// Likewise, the next row to process in this column is updated
			// only if we already reached the last one
			if nbrows[j] <= i {
				nbrows[j]++
			}
		}
	}

	// insert all splitters
	addSplitters(output)

	// and return the concatenation of all strings in the output string
	// separated by a newline
	return strings.Join(output, "\n")
}
