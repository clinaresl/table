package table

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Contents are formatters and therefore, they provide methods for both
// processing and formatting
//
// Processing a content involves splitting it across several rows if needed so
// that it satisfies the format of the column where it has to be shown
func (c content) Process(col column) []string {

	// if a paragraph alignment (p, C, L, R) modifier is used for this specific
	// column, then split it the content
	if col.hformat.alignment == 'p' ||
		col.hformat.alignment == 'C' ||
		col.hformat.alignment == 'L' ||
		col.hformat.alignment == 'R' {
		return splitParagraph(string(c), col.hformat.arg)
	}

	// if, on the other hand, a newline character has been provided, split the
	// content as well according to the newline characters
	re := regexp.MustCompile(newlineRegex)
	return re.Split(string(c), -1)
}

// Formatting a cell implies adding blank characters to a physical line so that
// it satisfies the format of the column where it has to be shown
func (c content) Format(col column) string {

	// in case it is necessary, the prefix and suffix contain a string of blank
	// characters to insert properly so that the contents satisfy the format of
	// this column
	var prefix, suffix string

	// compute the prefix to use for representing the contents of this column
	if unicode.ToLower(rune(col.hformat.alignment)) == 'c' {
		prefix = strings.Repeat(" ", (col.width-utf8.RuneCountInString(string(c)))/2)
	}
	if unicode.ToLower(rune(col.hformat.alignment)) == 'r' {
		prefix = strings.Repeat(" ", col.width-utf8.RuneCountInString(string(c)))
	}

	// compute the suffix to use for representing the contents of this column
	if unicode.ToLower(rune(col.hformat.alignment)) == 'c' {

		// note that in this case an additional character is added, i.e.,
		// centered strings are ragged left in case the difference is and odd
		// number
		suffix = strings.Repeat(" ", (col.width-utf8.RuneCountInString(string(c)))/2)
		suffix += strings.Repeat(" ", (col.width-utf8.RuneCountInString(string(c)))%2)
	}
	if unicode.ToLower(rune(col.hformat.alignment)) == 'l' || col.hformat.alignment == 'p' {
		suffix = strings.Repeat(" ", col.width-utf8.RuneCountInString(string(c)))
	}

	// and return the concatenation of the prefix, the content and the suffix
	return prefix + string(c) + suffix
}
