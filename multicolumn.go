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

// Return a new instance of a multicolumn. The first parameter is the number of
// columns that are grouped under the specification given next. Immediately
// after an arbitrary number of arguments satisfying the null interface can be
// given which substitute the given number of columns.
func Multicolumn(nbcolumns int, spec string, args ...interface{}) *multicolumn {

	// and return an instance of a multicolumn with no error. Note that the output is initially empty
	return &multicolumn{
		nbcolumns: nbcolumns,
		spec:      spec,
		args:      args,
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

	// Processing a multicolumn is truly easy. It just suffices creating an
	// ancilliary table with the specification given in the multicolumn. The
	// output of the process is just the lines that result from printing it

	// -- initialization
	var result []formatter

	// Create the ancilliary table and if any error is reported panic (as the
	// process step does not return any errors)
	if t, err := NewTable(m.spec); err != nil {
		panic("Panic while processing a multicolumn. It was not possible to create the ancilliary table!")
	} else {

		// if no error was reported during the creating of the table, just add a
		// row with all the arguments given in the multicolumn
		t.AddRow(m.args...)

		// and store all lines as different multicolumns where only the output
		// of each line is stored separately
		for _, line := range strings.Split(fmt.Sprintf("%v", t), "\n") {

			// note that only each line is computed separately
			result = append(result, formatter(multicolumn{output: line}))
		}
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
