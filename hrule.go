package table

// Horizontal rules are formatters and therefore, they provide methods for both
// processing and formatting
//
// Processing a horizontal rule does nothing at all as horizontal rules do not
// require space by themselves. Hence, they return a silce of strings comprising
// only one line with no character at all
func (h hrule) Process(col column) []string {

	return []string{""}
}

// Formatting a horizontal rule implies repeating the rune used to draw the
// horizontal rule as many times as required so that it takes the width of the
// column
func (h hrule) Format(col column) string {
	rule := ""
	for i := 0; i < col.width; i++ {
		rule += string(h)
	}

	return rule
}
