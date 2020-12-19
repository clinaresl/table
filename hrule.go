package table

import (
	"log"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Processing a cell means transforming logical rows into physical ones by
// splitting its contents across several (physical) rows, and also adding blank
// lines so that the result satisfies the vertical format of the column where it
// has to be shown, if and only if the height of the corresponding row is larger
// than the number of physical rows necessary to display the contents of the
// cell. To properly process a cell it is necessary to get a pointer to the
// table, and also the integer indices to the row and column of the cell
func (h hrule) Process(t *Table, irow, jcol int) []formatter {

	// the result of processing a horizontal rule consists of a single file
	// which is then stored in a single string. The result will then be
	// transformed into a slice of arrays containing only this string
	var splitters string

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
	// to use (either the one from the preceding column or the one following
	// immediately after) in case a rune found in the separator has to be
	// substituted
	offset := -1

	// in addition, it is necessary to consider that the last column introduces
	// an asymmetry which comes from the fact that it might contain or not a
	// vertical separator (in spite of it containing data or not). If no
	// vertical separator is present, then runes in the separator should be
	// substituted by them if no splitter is found; otherwise, the horizontal
	// separator used in the preceding column should be used
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
			west, _ = utf8.DecodeRuneInString(string(t.cells[irow][jcol-1].(hrule)))
			if west == horizontal_blank {
				west = none
			}
		}

		// east
		if jcol < t.GetNbColumns() {
			east, _ = utf8.DecodeRuneInString(string(t.cells[irow][jcol].(hrule)))
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
				brkrule, _ := utf8.DecodeRuneInString(string(t.cells[irow][jcol+offset].(hrule)))
				splitters += string(brkrule)
			}
		} else {

			// if a correct separator was returned, then add it to the string to
			// return
			splitters += string(splitter)
		}
	}

	// and return the string computed so far but in the form of a slice
	// containing only one line. Mind the trick: the splitters (which are a
	// standard string) are casted into a hrule to enable the specific format of
	// horizontal rules
	return []formatter{hrule(splitters)}
}

// Cells are also formatted (physical) line by line where each physical line is
// the result of processing cell (irow, jcol) and should be given in the
// receiver of this method. Each invocation returns a string where each
// (physical) line is forrmatted according to the horizontal format
// specification of the j-th column.
func (h hrule) Format(t *Table, irow, jcol int) string {

	// The result of formatting a horizontal rule consists of prefixing the
	// horizontal rule used in the data column with the splitters given in this
	// horizontal rule

	// the only task to do consists of repeating the separator as many times as
	// the width of this column after the splitters
	rule, ok := t.cells[irow][jcol].(hrule)
	if !ok {
		log.Fatalf(" The formatter in location (%v, %v) could not be casted into a rule!", irow, jcol)
	}
	return string(h) + strings.Repeat(string(rule), t.columns[jcol].width)
}
