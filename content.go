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
// that it satisfies the format of the column where it has to be shown. In
// addition, if the number of rows is strictly positive, then the content is
// vertically formatted
func (c content) Process(col column, nbrows int) (result []string) {

	// if a paragraph alignment (p, C, L, R) modifier is used for this specific
	// column, then split it the content
	if col.hformat.alignment == 'p' ||
		col.hformat.alignment == 'C' ||
		col.hformat.alignment == 'L' ||
		col.hformat.alignment == 'R' {
		result = splitParagraph(string(c), col.hformat.arg)
	}

	// if, on the other hand, a newline character has been provided, split the
	// content as well according to the newline characters
	re := regexp.MustCompile(newlineRegex)
	result = re.Split(string(c), -1)

	// if the number of physical rows is strictly positive, then proceed to
	// vertical format these contents
	if nbrows > 0 {

		var prefix, suffix int
		if unicode.ToLower(rune(col.vformat.alignment)) == 'c' {
			prefix = (nbrows - len(result)) / 2
		}
		if unicode.ToLower(rune(col.vformat.alignment)) == 'b' {
			prefix = nbrows - len(result)
		}

		if unicode.ToLower(rune(col.vformat.alignment)) == 't' {
			suffix = nbrows - len(result)
		}
		if unicode.ToLower(rune(col.vformat.alignment)) == 'c' {
			suffix = (nbrows - len(result)) / 2
			suffix += (nbrows - len(result)) % 2
		}

		// and now add the corresponding number of blank lines as required
		for iline := 0; iline < prefix; iline++ {
			result = prepend("", result)
		}
		for iline := 0; iline < suffix; iline++ {
			result = append(result, "")
		}
	}

	return
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
