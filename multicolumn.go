// -*- coding: utf-8 -*-
// multicolumn.go
// -----------------------------------------------------------------------------
//
// Started on <sáb 19-12-2020 22:45:26.735542876 (1608414326)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

package table

import (
	"fmt"
	"strings"
)

// ----------------------------------------------------------------------------
// Multicolumn
// ----------------------------------------------------------------------------

// Functions
// ----------------------------------------------------------------------------

// Multicolumns are meant to be inserted as ordinary cells in a data row. This
// function is intended to create and separately store multicolumns to be used
// later.
//
// Return a new instance of a multicolumn. The first parameter is the number of
// columns that are grouped under the specification given next. Immediately
// after an arbitrary number of arguments can be given which are formatted
// according to the column specification given.
//
// Importantly, the column specification is not allowed to end with any vertical
// separator ---or, in other words, all columns given in the column
// specification of the multicolumn must be data columns. Indeed, the row is
// continued with the vertical separator given in the first column of the table
// after the multicolumn. If the column specification ends with a vertical
// separator, or if any other error is found, an error is immediately raised
func NewMulticolumn(nbcolumns int, spec string, args ...interface{}) (multicolumn, error) {

	// create a table with the given column specification
	t, err := NewTable(spec)
	if err != nil {

		// Of course, if creating the table produces any error abort immediately
		return multicolumn{}, err
	}

	// Otherwise, verify that all columns in this table are data columns
	if len(t.columns) != t.GetNbColumns() {
		return multicolumn{}, fmt.Errorf("Invalid column specification of a multicolumn: '%v'", spec)
	}

	// finally, return an instance of a multicolumn with no error. Note that
	// both the initial column and the output are initially empty
	return multicolumn{
		nbcolumns: nbcolumns,
		spec:      spec,
		table:     *t,
		args:      args,
	}, nil
}

// Multicolumns are meant to be inserted as ordinary cells in a data row. This
// function is intended to be used straight ahead with AddRow.
//
// Return a new instance of a multicolumn. The first parameter is the number of
// columns that are grouped under the specification given next. Immediately
// after an arbitrary number of arguments can be given which are formatted
// according to the column specification given.
//
// This function uses NewMulticolumn and, if an error is returned, then it
// panics.
func Multicolumn(nbcolumns int, spec string, args ...interface{}) multicolumn {

	// create a new multi column
	if mcolumn, err := NewMulticolumn(nbcolumns, spec, args...); err != nil {

		// if an error is found, automatically panic. There is nothing better to
		// do as this function is intended to be used directly when adding
		// contents to a row.
		panic(err)
	} else {

		// if no error is spotted, then return a multicolumn
		return mcolumn
	}
}

// Methods
// ----------------------------------------------------------------------------

// Multicolumns are formatters and thus they should be both processed and
// formatted

// Processing a cell means transforming logical rows into physical ones by
// splitting its contents across several (physical) rows, and also adding blank
// lines so that the result satisfies the vertical format of the column where it
// has to be shown, if and only if the height of the corresponding row is larger
// than the number of physical rows necessary to display the contents of the
// cell. To properly process a cell it is necessary to get a pointer to the
// table, and also the integer indices to the row and column of the cell
func (m multicolumn) Process(t *Table, irow, jcol int) []formatter {

	// Processing a multicolumn is truly easy. It just suffices to add all arguments given

	// -- initialization
	var result []formatter

	// add all arguments to the multicolumn table
	m.table.AddRow(m.args...)

	// and store all lines as different multicolumns where only the output
	// of each line is stored separately
	for _, line := range strings.Split(fmt.Sprintf("%v", m.table), "\n") {

		// note that only each line is computed separately
		result = append(result, formatter(multicolumn{output: line}))
	}

	// and return the result computed so far
	return result
}

// Cells are also formatted (physical) line by line where each physical line is
// the result of processing cell (irow, jcol) and should be given in the
// receiver of this method. Each invocation returns a string where each
// (physical) line is forrmatted according to the horizontal format
// specification of the j-th column.
func (m multicolumn) Format(t *Table, irow, jcol int) string {

	// Formatting a multicolumn consists of simply returning its output string
	return m.output
}
