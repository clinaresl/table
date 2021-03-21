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
	"log"
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
// Importantly, the column specification is allowed to end with any vertical
// separator and no column specifier. In this case, the last separator is used
// as the separator of the cell following the multirow immediately after, as in
// LaTeX
func NewMulticolumn(nbcolumns int, spec string, args ...interface{}) (multicolumn, error) {

	// First things first, use only the specification that does not contain a
	// last separator without a column
	newspec, lastsep := stripLastSeparator(spec)

	// create a table with the processed column specification, i.e., the one
	// that does only contain columns preceded by a vertical separator (if any
	// is given)
	t, err := NewTable(newspec)
	if err != nil {

		// Of course, if creating the table produces any error abort immediately
		return multicolumn{}, err
	}

	// finally, return an instance of a multicolumn with no error. Note that
	// both the initial column and the output are initially empty
	return multicolumn{
		nbcolumns: nbcolumns,
		spec:      newspec,
		lastsep:   lastsep,
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

	// Processing a multicolumn is truly easy. It just suffices to add all
	// arguments given and to return a new multicolumn for each physical line

	// -- initialization
	var result []formatter

	// add all arguments to the multicolumn table
	m.table.AddRow(m.args...)

	// and store all lines as different multicolumns where only the output
	// of each line is stored separately
	for _, line := range strings.Split(fmt.Sprintf("%v", m.table), "\n") {

		// note that only each line is computed separately. In addition,
		// other information is passed to the multicolumn to be formatted
		result = append(result, formatter(multicolumn{
			jinit:     m.jinit,
			nbcolumns: m.nbcolumns,
			lastsep:   m.lastsep,
			table:     m.table,
			output:    line}))
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

	log.Printf("line: %v\n", m.output)
	log.Printf("tWidth: %v\n", t.getColumnsWidth(m.jinit, m.nbcolumns))
	log.Printf("\t init: %v\n", m.jinit)
	log.Printf("\t #cols: %v\n", m.nbcolumns)
	log.Printf("mWidth: %v\n", m.table.getColumnsWidth(0, len(m.table.columns)))
	log.Printf("\t init: %v\n", 0)
	log.Printf("\t #cols: %v\n\n", len(m.table.columns))

	// Formatting a multicolumn consists of simply returning its output string
	// but, in case this multicolumn spans until the right margin of the table,
	// then a separator has to be added in case this multicolumn has any
	if m.jinit+m.nbcolumns < len(t.columns) {
		return m.output
	}
	return m.output + m.lastsep
}
