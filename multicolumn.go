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
// function is intended to create and separately store multicolumns for further
// use.
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
	newspec, lastsep := stripLastSeparator(spec, colSpecRegex)

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

	// create a new multicolumn
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
	// arguments given and to return a new multicolumn for each physical line.
	// There is only one caveat to consider and it is to modify the table in
	// case another multicolumn precedes this one

	// -- initialization
	var result []formatter

	// if this multicell starts with no separator in the first column, and is
	// preceded by another multicolumn or multicell which in turn provides a
	// separator in a last column with no body, then use that separator. This
	// involves modifying the table which is stored within this multicell.The
	// reasoning is simple:
	//
	//    1. Multicolumns and multicells are allowed to overwrite the separators
	//    given in the column specification of the table
	//
	//    2. Also, multicolumns and multicells are allowed to affect the
	//    separator to be used in the cell coming immediately after
	//
	// So, if a multicell starts with no separator, then the separator given in
	// the column specification of the table is not used at all and the only
	// chance to create a separator is to use the one given by the preceding
	// cell only if that is a multicolumn or multicell. Note that the following
	// verification is performed only in case the current row has been written
	// to the table (this is not the case when adding rows, but it is when
	// printing the contents of tables)
	if m.table.columns[0].sep == "" && irow < len(t.rows) {
		if mprev := getPreviousBinder(t, irow, jcol); mprev != nil && mprev.getLastVerticalSep() != "" {

			// redo the table using as first separator the one provided in the
			// previous multicolumn/multicell. In case of error (which is
			// unlikely as we are only adding the separator found in the
			// previous multicolumn/multicell) a panic is generated, there's not
			// much we could do at this stage
			if tm, err := NewTable(mprev.getLastVerticalSep() + m.spec); err != nil {
				panic(err)
			} else {

				// the following is ugly ... I know :( and it is a little bit of
				// hacking. Multicolumns can be processed after distributing
				// some space among its columns. Thus, if we are re-creating the
				// inner table of a multicolumn, it is more than a good idea to
				// preserve the widths of all its columns
				for idx, _ := range tm.columns {
					tm.columns[idx].width = m.table.columns[idx].width
				}
				tm.columns[0].width -= countPrintableRuneInString(mprev.getLastVerticalSep())
				m.table = *tm
			}
		}
	}

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

	// Formatting a multicolumn consists of simply returning its output string
	// but, in case this multicolumn spans until the right margin of the table,
	// then a separator has to be added in case this multicolumn has any
	if m.jinit+m.nbcolumns < len(t.columns) {
		return m.output
	}
	return m.output + m.lastsep
}

// Multicolumns are also binders. Binders are any content that either merges
// cells/rows or that splits them such as multicolumns and multicells. The main
// difference between binders and ordinary contents is that the contents of a
// binder are formatted with the assitance of a nested table

// Binders are allowed to modify the vertical horizontal separator of the next
// cell in the table. This applies only if a binder is formatted with a last
// column which contains no column specification, i.e., no data. The following
// service provides such separator or an empty string if none is given
func (m multicolumn) getLastVerticalSep() string {
	return m.lastsep
}

// Likewise, they are allowed also to modify the horizontal separator of the
// next horizontal rule. This applies only if a binder is formatted with a last
// row which contains no row specification, i.e., no data. The following service
// provides such separator or an empty string if none is given
func (m multicolumn) getLastHorizontalSep() string {

	// As multicolumns do not have a row specification, they can not by any
	// means modify the horizontal rule below them
	return ""
}

// Binders span over an arbitrary number of consecutive columns. The
// following services provide the first column they take and the number of
// consecutive columns they span over
func (m multicolumn) getColumnInit() int {
	return m.jinit
}
func (m multicolumn) getNbColumns() int {
	return m.nbcolumns
}

// In the same vain, binders span over an arbitrary number of consecutive
// rows. The following services provide the first row they take and the
// number of consecutive rows they span over
func (m multicolumn) getRowInit() int {

	// Multicolumns only take one line. Hence, they return -1 to indicate that
	// no line has been given
	return -1
}
func (m multicolumn) getNbRows() int {

	// Multicolumns only take one line
	return 1
}

// The key difference with contents is that binders are processed with the
// contents of a nested table which is returned with the following service
func (m multicolumn) getTable() *Table {
	return &m.table
}
