// -*- coding: utf-8 -*-
// multicell.go
// -----------------------------------------------------------------------------
//
// Started on <sáb 27-03-2021 12:45:36.225817391 (1616845536)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

package table

import (
	"fmt"
	"strings"
)

// ----------------------------------------------------------------------------
// Multicell
// ----------------------------------------------------------------------------

// Functions
// ----------------------------------------------------------------------------

// Multicells are meant to be inserted as ordinary cells into a table. This
// function is inteded to create and separately store multicells for further use.
//
// Return a new instance of a multicell. The first two parameters are the number
// of columns and rows that are grouped under the multicells which is is
// formatted according to the column and row specifications given next.
// Immediately after an arbitrary number of arguments can be given which are
// formatted according to the specifications given.
//
// Importantly, both the column and row specifications are allowed to end with a
// last separator (i.e., a separator with no content specifier following
// immediately after). If a last separator is given in the column specification,
// it is used as the separator of the next cell; if one is given in the row
// specification, it is then used as the horizontal rule of the next row.
func NewMulticell(nbcolumns, nbrows int, cspec, rspec string, args ...interface{}) (multicell, error) {

	// First things first, strip the last separator, if any is given from the
	// column and row specifications
	cnewspec, clastsep := stripLastSeparator(cspec, colSpecRegex)
	rnewspec, rlastsep := stripLastSeparator(rspec, rowSpecRegex)

	// create a table with the processed column and row specifications
	t, err := NewTable(cnewspec, rspec)
	if err != nil {

		// Of course, if creating the table produces any error abort immediately
		return multicell{}, err
	}

	// finally, return an instance of a multicell with no error. Note that the
	// initial column/row and the output are initially empty
	return multicell{
		nbcolumns: nbcolumns,
		nbrows:    nbrows,
		cspec:     cnewspec,
		rspec:     rnewspec,
		clastsep:  clastsep,
		rlastsep:  rlastsep,
		table:     *t,
		args:      args,
	}, nil
}

// Multicells are meant to be inserted as ordinary cells in a table. This
// function is intended to be used straight ahead with AddRow.
//
// Return a new instance of a multicell. The first two parameters are the number
// of columns and rows that are grouped under the multicells which is is
// formatted according to the column and row specifications given next.
// Immediately after an arbitrary number of arguments can be given which are
// formatted according to the specifications given.
//
// This function uses NewMulticell and, if an error is returned, then it panics.
func Multicell(nbcolumns, nbrows int, cspec, rspec string, args ...interface{}) multicell {

	// create a new multicell
	if mcell, err := NewMulticell(nbcolumns, nbrows, cspec, rspec, args...); err != nil {

		// if an error is found, automatically panic. There is nothing better to
		// do as this function is intended to be used directly when adding
		// contents to a row.
		panic(err)
	} else {

		// if no error is spotted, then return a multicell
		return mcell
	}
}

// Methods
// ----------------------------------------------------------------------------

// Multicells are formatters and thus they should be both processed and
// formatted

// Processing a cell means transforming logical rows into physical ones by
// splitting its contents across several (physical) rows, and also adding blank
// lines so that the result satisfies the vertical format of the column where it
// has to be shown, if and only if the height of the corresponding row is larger
// than the number of physical rows necessary to display the contents of the
// cell. To properly process a cell it is necessary to get a pointer to the
// table, and also the integer indices to the row and column of the cell
func (m multicell) Process(t *Table, irow, jcol int) []formatter {

	// Processing a multicell is truly easy. It just suffices to add all
	// arguments given and to return a new multicell for each physical line.
	// There is only one caveat to consider and it is to modify the table in
	// case another multicolumn or multicell precedes this one

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
			if tm, err := NewTable(mprev.getLastVerticalSep() + m.cspec); err != nil {
				panic(err)
			} else {

				// the following is ugly ... I know :( and it is a little bit of
				// hacking. Multicells are processed after distributing some
				// space among its columns. Thus, if we are re-creating the
				// inner table of a multicell, it is more than a good idea to
				// preserve the widths of all its columns
				for idx, _ := range tm.columns {
					tm.columns[idx].width = m.table.columns[idx].width
				}
				tm.columns[0].width -= countPrintableRuneInString(mprev.getLastVerticalSep())
				m.table = *tm
			}
		}
	}

	// add all arguments to the multicell table. Note that in the case of
	// multicells, arguments are given for the whole table so that it is
	// necessary now to arrange them by rows
	for iditem := 0; iditem < len(m.args); iditem += len(m.table.columns) {
		m.table.AddRow(m.args[iditem:min(iditem+len(m.table.columns), len(m.args))]...)
	}

	// and store all lines as different multicells where only the output of each
	// line is stored separately
	for _, line := range strings.Split(fmt.Sprintf("%v", m.table), "\n") {

		// note that only each line is computed separately. In addition,
		// other information is passed to the multicell to be formatted
		result = append(result, formatter(multicell{
			jinit:     m.jinit,
			nbcolumns: m.nbcolumns,
			iinit:     m.iinit,
			nbrows:    m.nbrows,
			clastsep:  m.clastsep,
			rlastsep:  m.rlastsep,
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
func (m multicell) Format(t *Table, irow, jcol int) string {

	// Formatting a multicell consists of simply returning its output string
	// but, in case this multicell spans until the right margin of the table,
	// then a separator has to be added in case this multicell has any
	if m.jinit+m.nbcolumns < len(t.columns) {
		return m.output
	}
	return m.output + m.clastsep
}

// Multicells are also binders. Binders are any content that either merges
// cells/rows or that splits them such as multicolumns and multicells. The main
// difference between binders and ordinary contents is that the contents of a
// binder are formatted with the assitance of a nested table

// Binders are allowed to modify the vertical horizontal separator of the next
// cell in the table. This applies only if a binder is formatted with a last
// column which contains no column specification, i.e., no data. The following
// service provides such separator or an empty string if none is given
func (m multicell) getLastVerticalSep() string {
	return m.clastsep
}

// Likewise, they are allowed also to modify the horizontal separator of the
// next horizontal rule. This applies only if a binder is formatted with a last
// row which contains no row specification, i.e., no data. The following service
// provides such separator or an empty string if none is given
func (m multicell) getLastHorizontalSep() string {
	return m.rlastsep
}

// Binders span over an arbitrary number of consecutive columns. The
// following services provide the first column they take and the number of
// consecutive columns they span over
func (m multicell) getColumnInit() int {
	return m.jinit
}
func (m multicell) getNbColumns() int {
	return m.nbcolumns
}

// In the same vain, binders span over an arbitrary number of consecutive
// rows. The following services provide the first row they take and the
// number of consecutive rows they span over
func (m multicell) getRowInit() int {
	return m.iinit
}
func (m multicell) getNbRows() int {
	return m.nbrows
}

// The key difference with contents is that binders are processed with the
// contents of a nested table which is returned with the following service
func (m multicell) getTable() *Table {
	return &m.table
}

// Local Variables:
// mode:go
// fill-column:80
// End:
