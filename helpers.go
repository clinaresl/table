// -*- coding: utf-8 -*-
// helpers.go
// -----------------------------------------------------------------------------
//
// Started on <sáb 19-12-2020 22:45:26.735542876 (1608414326)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

package table

import (
	"errors"
	"fmt"
	"regexp"
	"unicode"
	"unicode/utf8"
)

// Functions
// ----------------------------------------------------------------------------

// Return the maximum of two integers
func max(n, m int) int {
	if n > m {
		return n
	}
	return m
}

// process the given column specification and return a slice of instances of
// columns properly initialized. In case the parsing was not possible an error
// is returned
func getColumns(colspec string) ([]column, error) {

	// --initialization
	var columns []column

	// the specification is processed with a regular expression which should be
	// used to consume the whole string
	re := regexp.MustCompile(colSpecRegex)
	for {

		// get the next column and, if none is found, then exit
		recol := re.FindStringIndex(colspec)
		if recol == nil {
			break
		}
		nxtcol, err := newColumn(colspec[recol[0]:recol[1]])
		if err != nil {
			return []column{}, err
		}

		// and add it to the columns of this table
		columns = append(columns, *nxtcol)

		// and now move forward in the column specification string
		colspec = colspec[recol[1]:]
	}

	// maybe the column specification string is not empty here. Any remainings
	// are interpreted as the separator of a last column which contains no text
	// and which has no format
	if colspec != "" {
		columns = append(columns,
			column{sep: colspec,
				hformat: style{},
				vformat: style{}})
	}

	// and return the slice of columns along with no error
	return columns, nil
}

// return true if and only if the given rune is recognized as a vertical
// separator as defined in this package and false otherwise
func isVerticalSeparator(r rune) bool {
	return r == '│' || r == '║' || r == '┃'
}

// return true if and only if the given string contains a vertical separator as
// defined in this package and false otherwise
func containsVerticalSeparator(sep string) bool {

	// for each rune in the given string
	for _, r := range sep {

		// if this rune is a vertical separator then return true immediately
		if isVerticalSeparator(r) {
			return true
		}
	}

	// otherwise, return false
	return false
}

// Just cast a slice of strings into a slice of contents
func strToContent(input []string) (output []content) {

	for _, str := range input {
		output = append(output, content(str))
	}
	return
}

// Return the number of runes in the given string which are printable
func countPrintableRuneInString(s string) (count int) {

	// first things first, remove all ANSI color escape sequences
	re := regexp.MustCompile(ansiColorRegex)
	s = re.ReplaceAllString(s, "")

	// Initialize the counter
	count = 0

	// and now, count all both printable and graphic runes in the resulting
	// string after removing the ANSI color escape sequences
	for _, r := range s {
		if unicode.IsGraphic(r) && unicode.IsPrint(r) {
			count += 1
		}
	}
	return
}

// the following function returns a slice of strings with the same contents than
// the input string (with some spaces removed) such that the length of each
// string is the larger one less or equal than the given width
func splitParagraph(str string, width int) (result []string) {

	// iterate over all runes of the input string
	for len(str) > 0 {

		// while processing a substring, keep track of the number of runes in it
		// and also the location of the last byte to add to it. In addition, it
		// is required to store the position of the rune to start considering in
		// the next cycle
		var nbrunes, end, nxt int
		for pos, rune := range str {

			// accept this rune
			nbrunes += 1

			// in case this is a space (including utf-8 spaces) then remember
			// the location of the last position to include in the current
			// substring
			if unicode.IsSpace(rune) {
				end, nxt = pos, utf8.RuneLen(rune)

				// and, in case this is a newline character, then exit
				// immediately from the inner loop
				if rune == '\n' {
					break
				}
			}

			// If the maximum number of runes to add has been reached then break
			// avoiding adding more runes
			if nbrunes >= width {

				// if no breaking point has been found before then add all runes
				// until the current location
				if end == 0 {
					end, nxt = pos+utf8.RuneLen(rune), 0
				}

				// If the character immediately after this one is a space then
				// add all runes until this location also
				nxtrune, _ := utf8.DecodeRuneInString(str[pos+utf8.RuneLen(rune):])
				if unicode.IsSpace(nxtrune) {
					end, nxt = pos+utf8.RuneLen(rune), utf8.RuneLen(rune)
				}

				break
			}

			// Finally, if the whole string has been exhausted, then add it
			// until the end
			if pos+utf8.RuneLen(rune) >= len(str) {
				end, nxt = len(str), 0
			}
		}

		// add the substring from the beginning of the input string until the
		// end
		result = append(result, str[:end])

		// and move forward in the string
		str = str[end+nxt:]
	}

	return
}

// return the rune that splits the four regions north-west, north-east,
// south-west and south-east as stored in the map of splitters with no error. If
// such splitter does not exist, it returns none with an error
func getSingleSplitter(west, east, north, south rune) (rune, error) {

	// check for the existence of the west rune. In case it does not exist,
	// return an error
	if _, ok := splitterUTF8[west]; !ok {
		return none, errors.New("No splitter found")
	}

	// east
	if _, ok := splitterUTF8[west][east]; !ok {
		return none, errors.New("No splitter found")
	}

	// north
	if _, ok := splitterUTF8[west][east][north]; !ok {
		return none, errors.New("No splitter found")
	}

	// south
	if _, ok := splitterUTF8[west][east][north][south]; !ok {
		return none, errors.New("No splitter found")
	}

	// and return the corresponding splitter which, at this point, is guaranteed
	// to exist
	return splitterUTF8[west][east][north][south], nil
}

// return a slice with all runes in a string
func getRunes(s string) (runes []rune) {

	// for all runes in the string
	for _, r := range s {

		// add this rune to the slice of runes to return
		runes = append(runes, r)
	}
	return
}

// return a slice of vertical specifications as a slice of styles. In case the
// row specification is incorrect, an error is returned and the contents of the
// result are undetermined
func getVerticalStyles(rowspec string) ([]style, error) {

	var result []style

	// while the row specification is not empty. Yeah, the row specification
	// should not consist of runes but just simple ascii characters. Still, we
	// traverse the string as runes
	for _, rune := range rowspec {
		switch rune {
		case 't', 'b', 'c':
			result = append(result, style{alignment: byte(rune)})
		default:
			return result, fmt.Errorf("'%v' is an incorrect vertical format", string(rune))
		}
	}

	return result, nil
}

// The following function prepends the given argument to the slice of contents
// given second
func prepend(item content, data []content) []content {

	// just add an item to the slice, copy all items shifting them all by one
	// position to the right and overwrite the first item
	data = append(data, item)
	copy(data[1:], data)
	data[0] = item

	return data
}

// Evenly increment the width of all columns given in the slice of columns so
// that their accumulated sum is incremented by n
func distribute(n int, columns []column) {

	// get the integer quotient and remainder of n and the length of the slice
	quotient, remainder := n/len(columns), n%len(columns)

	// distribute the quotient among all columns
	for idx, _ := range columns {
		columns[idx].width += quotient
	}

	// and now distribute the remainder among the first columns
	for idx := 0; idx < remainder; idx++ {
		columns[idx].width += 1
	}
}
