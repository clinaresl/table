package table

import (
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
