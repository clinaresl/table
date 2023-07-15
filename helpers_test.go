// -*- coding: utf-8 -*-
// helpers_test.go
// -----------------------------------------------------------------------------
//
// Started on <sáb 19-12-2020 22:45:26.735542876 (1608414326)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

package table

import (
	"reflect"
	"testing"
)

func Test_stripLastSeparator(t *testing.T) {
	type args struct {
		spec, rexp string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{

		// COLUMN SPECIFICATION

		// No vertical separator in the last column

		// One column
		{args: args{spec: "l",
			rexp: colSpecRegex},
			want:  "l",
			want1: ""},

		{args: args{spec: "c",
			rexp: colSpecRegex},
			want:  "c",
			want1: ""},

		{args: args{spec: "r",
			rexp: colSpecRegex},
			want:  "r",
			want1: ""},

		{args: args{spec: "p{120}",
			rexp: colSpecRegex},
			want:  "p{120}",
			want1: ""},

		{args: args{spec: "L{120}",
			rexp: colSpecRegex},
			want:  "L{120}",
			want1: ""},

		{args: args{spec: "C{120}",
			rexp: colSpecRegex},
			want:  "C{120}",
			want1: ""},

		{args: args{spec: "R{120}",
			rexp: colSpecRegex},
			want:  "R{120}",
			want1: ""},

		// Two columns
		{args: args{spec: "l|c",
			rexp: colSpecRegex},
			want:  "l|c",
			want1: ""},

		{args: args{spec: "|l|c",
			rexp: colSpecRegex},
			want:  "|l|c",
			want1: ""},

		{args: args{spec: "c|l",
			rexp: colSpecRegex},
			want:  "c|l",
			want1: ""},

		{args: args{spec: "|c|l",
			rexp: colSpecRegex},
			want:  "|c|l",
			want1: ""},

		{args: args{spec: "r|p{120}",
			rexp: colSpecRegex},
			want:  "r|p{120}",
			want1: ""},

		{args: args{spec: "|r|p{120}",
			rexp: colSpecRegex},
			want:  "|r|p{120}",
			want1: ""},

		{args: args{spec: "p{120}|L{120}",
			rexp: colSpecRegex},
			want:  "p{120}|L{120}",
			want1: ""},

		{args: args{spec: "|p{120}|L{120}",
			rexp: colSpecRegex},
			want:  "|p{120}|L{120}",
			want1: ""},

		{args: args{spec: "L{120}|C{120}",
			rexp: colSpecRegex},
			want:  "L{120}|C{120}",
			want1: ""},

		{args: args{spec: "|L{120}|C{120}",
			rexp: colSpecRegex},
			want:  "|L{120}|C{120}",
			want1: ""},

		{args: args{spec: "C{120}|R{120}",
			rexp: colSpecRegex},
			want:  "C{120}|R{120}",
			want1: ""},

		{args: args{spec: "|C{120}|R{120}",
			rexp: colSpecRegex},
			want:  "|C{120}|R{120}",
			want1: ""},

		{args: args{spec: "R{120}|l",
			rexp: colSpecRegex},
			want:  "R{120}|l",
			want1: ""},

		{args: args{spec: "|R{120}|l",
			rexp: colSpecRegex},
			want:  "|R{120}|l",
			want1: ""},

		// With a vertical separator in the last column

		// One column
		{args: args{spec: "l|",
			rexp: colSpecRegex},
			want:  "l",
			want1: "│"},

		{args: args{spec: "c||",
			rexp: colSpecRegex},
			want:  "c",
			want1: "║"},

		{args: args{spec: "r|||",
			rexp: colSpecRegex},
			want:  "r",
			want1: "┃"},

		{args: args{spec: "p{120} | ",
			rexp: colSpecRegex},
			want:  "p{120}",
			want1: " │ "},

		{args: args{spec: "L{120} || ",
			rexp: colSpecRegex},
			want:  "L{120}",
			want1: " ║ "},

		{args: args{spec: "C{120} ||| ",
			rexp: colSpecRegex},
			want:  "C{120}",
			want1: " ┃ "},

		{args: args{spec: "R{120}  |",
			rexp: colSpecRegex},
			want:  "R{120}",
			want1: "  │"},

		// Two columns
		{args: args{spec: "l|c ",
			rexp: colSpecRegex},
			want:  "l|c",
			want1: " "},

		{args: args{spec: "|l|c |",
			rexp: colSpecRegex},
			want:  "|l|c",
			want1: " │"},

		{args: args{spec: "c|l||",
			rexp: colSpecRegex},
			want:  "c|l",
			want1: "║"},

		{args: args{spec: "|c|l||| ",
			rexp: colSpecRegex},
			want:  "|c|l",
			want1: "┃ "},

		{args: args{spec: "r|p{120}    |    ",
			rexp: colSpecRegex},
			want:  "r|p{120}",
			want1: "    │    "},

		{args: args{spec: "|r|p{120}||  ",
			rexp: colSpecRegex},
			want:  "|r|p{120}",
			want1: "║  "},

		{args: args{spec: "p{120}|L{120}   |||",
			rexp: colSpecRegex},
			want:  "p{120}|L{120}",
			want1: "   ┃"},

		{args: args{spec: "|p{120}|L{120} | ",
			rexp: colSpecRegex},
			want:  "|p{120}|L{120}",
			want1: " │ "},

		{args: args{spec: "L{120}|C{120} ||",
			rexp: colSpecRegex},
			want:  "L{120}|C{120}",
			want1: " ║"},

		{args: args{spec: "|L{120}|C{120}||   ",
			rexp: colSpecRegex},
			want:  "|L{120}|C{120}",
			want1: "║   "},

		{args: args{spec: "C{120}|R{120}   |||      ",
			rexp: colSpecRegex},
			want:  "C{120}|R{120}",
			want1: "   ┃      "},

		{args: args{spec: "|C{120}|R{120} |",
			rexp: colSpecRegex},
			want:  "|C{120}|R{120}",
			want1: " │"},

		{args: args{spec: "R{120}|l| ",
			rexp: colSpecRegex},
			want:  "R{120}|l",
			want1: "│ "},

		{args: args{spec: "|R{120}|l | ",
			rexp: colSpecRegex},
			want:  "|R{120}|l",
			want1: " │ "},

		// ROW SPECIFICATION

		// No horizontal separator in the last row

		// One row
		{args: args{spec: "t",
			rexp: rowSpecRegex},
			want:  "t",
			want1: ""},

		{args: args{spec: "c",
			rexp: rowSpecRegex},
			want:  "c",
			want1: ""},

		{args: args{spec: "b",
			rexp: rowSpecRegex},
			want:  "b",
			want1: ""},

		// Two rows

		{args: args{spec: "t|t",
			rexp: rowSpecRegex},
			want:  "t|t",
			want1: ""},

		{args: args{spec: "c|t",
			rexp: rowSpecRegex},
			want:  "c|t",
			want1: ""},

		{args: args{spec: "b|t",
			rexp: rowSpecRegex},
			want:  "b|t",
			want1: ""},

		{args: args{spec: "t|c",
			rexp: rowSpecRegex},
			want:  "t|c",
			want1: ""},

		{args: args{spec: "c|c",
			rexp: rowSpecRegex},
			want:  "c|c",
			want1: ""},

		{args: args{spec: "b|c",
			rexp: rowSpecRegex},
			want:  "b|c",
			want1: ""},

		{args: args{spec: "t|b",
			rexp: rowSpecRegex},
			want:  "t|b",
			want1: ""},

		{args: args{spec: "c|b",
			rexp: rowSpecRegex},
			want:  "c|b",
			want1: ""},

		{args: args{spec: "b|b",
			rexp: rowSpecRegex},
			want:  "b|b",
			want1: ""},

		// With a horizontal separator in the last row

		// One row
		{args: args{spec: "t|",
			rexp: rowSpecRegex},
			want:  "t",
			want1: "│"},

		{args: args{spec: "c||",
			rexp: rowSpecRegex},
			want:  "c",
			want1: "║"},

		{args: args{spec: "b|||",
			rexp: rowSpecRegex},
			want:  "b",
			want1: "┃"},

		// Two rows

		{args: args{spec: "t|t |",
			rexp: rowSpecRegex},
			want:  "t|t",
			want1: " │"},

		{args: args{spec: "c|t|| ",
			rexp: rowSpecRegex},
			want:  "c|t",
			want1: "║ "},

		{args: args{spec: "b|t ||| ",
			rexp: rowSpecRegex},
			want:  "b|t",
			want1: " ┃ "},

		{args: args{spec: "t|c|  ",
			rexp: rowSpecRegex},
			want:  "t|c",
			want1: "│  "},

		{args: args{spec: "c|c  || ",
			rexp: rowSpecRegex},
			want:  "c|c",
			want1: "  ║ "},

		{args: args{spec: "b|c   |||",
			rexp: rowSpecRegex},
			want:  "b|c",
			want1: "   ┃"},

		{args: args{spec: "t|b |    ",
			rexp: rowSpecRegex},
			want:  "t|b",
			want1: " │    "},

		{args: args{spec: "c|b>>||<<",
			rexp: rowSpecRegex},
			want:  "c|b",
			want1: ">>║<<"},

		{args: args{spec: "b|b····",
			rexp: rowSpecRegex},
			want:  "b|b",
			want1: "····"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := stripLastSeparator(tt.args.spec, tt.args.rexp)
			if got != tt.want {
				t.Errorf("stripLastSeparator() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("stripLastSeparator() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSplitParagraph(t *testing.T) {
	type args struct {
		str   string
		width int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{

		{args: args{str: "\nEn un lugar de la Mancha de cuyo nombre no quiero acordarme ...\n\n\nY colorín colorado, este cuento se ha acabado\n\n",
			width: 1},
			want: []string{"",
				"E",
				"n",
				"u",
				"n",
				"l",
				"u",
				"g",
				"a",
				"r",
				"d",
				"e",
				"l",
				"a",
				"M",
				"a",
				"n",
				"c",
				"h",
				"a",
				"d",
				"e",
				"c",
				"u",
				"y",
				"o",
				"n",
				"o",
				"m",
				"b",
				"r",
				"e",
				"n",
				"o",
				"q",
				"u",
				"i",
				"e",
				"r",
				"o",
				"a",
				"c",
				"o",
				"r",
				"d",
				"a",
				"r",
				"m",
				"e",
				".",
				".",
				".",
				"",
				"",
				"Y",
				"c",
				"o",
				"l",
				"o",
				"r",
				"í",
				"n",
				"c",
				"o",
				"l",
				"o",
				"r",
				"a",
				"d",
				"o",
				",",
				"e",
				"s",
				"t",
				"e",
				"c",
				"u",
				"e",
				"n",
				"t",
				"o",
				"s",
				"e",
				"h",
				"a",
				"a",
				"c",
				"a",
				"b",
				"a",
				"d",
				"o",
				""}},

		{args: args{str: "\nEn un lugar de la Mancha de cuyo nombre no quiero acordarme ...\n\n\nY colorín colorado, este cuento se ha acabado\n\n",
			width: 5},
			want: []string{"",
				"En un",
				"lugar",
				"de la",
				"Manch",
				"a de",
				"cuyo",
				"nombr",
				"e no",
				"quier",
				"o",
				"acord",
				"arme",
				"...",
				"",
				"",
				"Y",
				"color",
				"ín",
				"color",
				"ado,",
				"este",
				"cuent",
				"o se",
				"ha",
				"acaba",
				"do",
				""}},

		{args: args{str: "\nEn un lugar de la Mancha de cuyo nombre no quiero acordarme ...\n\n\nY colorín colorado, este cuento se ha acabado\n\n",
			width: 10},
			want: []string{"",
				"En un",
				"lugar de",
				"la Mancha",
				"de cuyo",
				"nombre no",
				"quiero",
				"acordarme",
				"...",
				"",
				"",
				"Y colorín",
				"colorado,",
				"este",
				"cuento se",
				"ha acabado",
				""}},

		{args: args{str: "\nEn un lugar de la Mancha de cuyo nombre no quiero acordarme ...\n\n\nY colorín colorado, este cuento se ha acabado\n\n",
			width: 20},
			want: []string{"",
				"En un lugar de la",
				"Mancha de cuyo",
				"nombre no quiero",
				"acordarme ...",
				"",
				"",
				"Y colorín colorado,",
				"este cuento se ha",
				"acabado",
				""}},

		{args: args{str: "\nEn un lugar de la Mancha de cuyo nombre no quiero acordarme ...\n\n\nY colorín colorado, este cuento se ha acabado\n\n",
			width: 30},
			want: []string{"",
				"En un lugar de la Mancha de",
				"cuyo nombre no quiero",
				"acordarme ...",
				"",
				"",
				"Y colorín colorado, este",
				"cuento se ha acabado",
				""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitParagraph(tt.args.str, tt.args.width); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitParagraph() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func Test_physicalToLogical(t *testing.T) {
// 	type args struct {
// 		s  string
// 		pi int
// 	}
// 	tests := []struct {
// 		name   string
// 		args   args
// 		wantLi int
// 	}{

// 		// trivial example with the empty string
// 		{args: args{s: "",
// 			pi: -1},
// 			wantLi: -1},

// 		{args: args{s: "",
// 			pi: 0},
// 			wantLi: -1},

// 		{args: args{s: "",
// 			pi: 1},
// 			wantLi: -1},

// 		// examples with no ANSI color codes
// 		{args: args{s: "Gladiator in arena consilium capit",
// 			pi: -1},
// 			wantLi: -1},

// 		{args: args{s: "Gladiator in arena consilium capit",
// 			pi: 0},
// 			wantLi: 0},

// 		{args: args{s: "Gladiator in arena consilium capit",
// 			pi: 16},
// 			wantLi: 16},

// 		{args: args{s: "Gladiator in arena consilium capit",
// 			pi: 33},
// 			wantLi: 33},

// 		{args: args{s: "Gladiator in arena consilium capit",
// 			pi: 34},
// 			wantLi: -1},

// 		// examples with only one ANSI color code
// 		{args: args{s: "\033[38;2;160;10;10mGladiator in arena consilium capit",
// 			pi: -1},
// 			wantLi: -1},

// 		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit",
// 			pi: -1},
// 			wantLi: -1},

// 		{args: args{s: "\033[38;2;160;10;10mGladiator in arena consilium capit",
// 			pi: 17},
// 			wantLi: 0},

// 		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit",
// 			pi: 0},
// 			wantLi: 0},

// 		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit",
// 			pi: 18},
// 			wantLi: 1},

// 		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit",
// 			pi: 26},
// 			wantLi: 9},

// 		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit",
// 			pi: 50},
// 			wantLi: 33},

// 		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit",
// 			pi: 51},
// 			wantLi: -1},

// 		// examples with two ANSI color codes
// 		{args: args{s: "\033[38;2;160;10;10mGladiator \033[0min arena consilium capit",
// 			pi: -1},
// 			wantLi: -1},

// 		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
// 			pi: -1},
// 			wantLi: -1},

// 		{args: args{s: "\033[38;2;160;10;10mGladiator \033[0min arena consilium capit",
// 			pi: 17},
// 			wantLi: 0},

// 		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
// 			pi: 0},
// 			wantLi: 0},

// 		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
// 			pi: 18},
// 			wantLi: 1},

// 		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
// 			pi: 26},
// 			wantLi: 9},

// 		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
// 			pi: 31},
// 			wantLi: 10},

// 		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
// 			pi: 54},
// 			wantLi: 33},

// 		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
// 			pi: 55},
// 			wantLi: -1},

// 		{args: args{s: "\033[38;2;160;10;10mGladiator in arena consilium capit\033[0m",
// 			pi: -1},
// 			wantLi: -1},

// 		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit\033[0m",
// 			pi: -1},
// 			wantLi: -1},

// 		{args: args{s: "\033[38;2;160;10;10mGladiator in arena consilium capit\033[0m",
// 			pi: 17},
// 			wantLi: 0},

// 		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit\033[0m",
// 			pi: 0},
// 			wantLi: 0},

// 		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit\033[0m",
// 			pi: 18},
// 			wantLi: 1},

// 		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit\033[0m",
// 			pi: 26},
// 			wantLi: 9},

// 		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit\033[0m",
// 			pi: 31},
// 			wantLi: 14},

// 		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit\033[0m",
// 			pi: 50},
// 			wantLi: 33},

// 		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit\033[0m",
// 			pi: 51},
// 			wantLi: -1},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if gotLi := physicalToLogical(tt.args.s, tt.args.pi); gotLi != tt.wantLi {
// 				t.Errorf("physicalToLogical() = %v, want %v", gotLi, tt.wantLi)
// 			}
// 		})
// 	}
// }

func Test_countPrintableRuneInString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
	}{

		// empty string (without and with color ANSI codes)
		{args: args{s: ""},
			wantCount: 0},

		{args: args{s: "\033[38;2;160;10;10m\033[0m"},
			wantCount: 0},

		// non-emtpy strings without color ANSI codes
		{args: args{s: "Gladiator in arena consilium capit"},
			wantCount: 34},

		// non-emtpy strings with one color ANSI codes
		{args: args{s: "\033[38;2;160;10;10mGladiator in arena consilium capit"},
			wantCount: 34},

		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit"},
			wantCount: 34},

		{args: args{s: "Gladiator in arena consilium capit\033[38;2;160;10;10m"},
			wantCount: 34},

		// non-emtpy strings with two color ANSI codes
		{args: args{s: "\033[38;2;160;10;10m\033[0mGladiator in arena consilium capit"},
			wantCount: 34},

		{args: args{s: "\033[38;2;160;10;10mGladiator \033[0min arena consilium capit"},
			wantCount: 34},

		{args: args{s: "\033[38;2;160;10;10mGladiator in arena consilium capit\033[0m"},
			wantCount: 34},

		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit"},
			wantCount: 34},

		{args: args{s: "G\033[38;2;160;10;10mladiator in arena\033[0m consilium capit"},
			wantCount: 34},

		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit\033[0m"},
			wantCount: 34},

		{args: args{s: "Gladiator in arena consilium capit\033[38;2;160;10;10m\033[0m"},
			wantCount: 34},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCount := countPrintableRuneInString(tt.args.s); gotCount != tt.wantCount {
				t.Errorf("countPrintableRuneInString() = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}

func Test_getRune(t *testing.T) {
	type args struct {
		s string
		i int
	}
	tests := []struct {
		name    string
		args    args
		want    rune
		wantErr bool
	}{

		// simple tests: errors are not tested
		{args: args{s: "Gladiator in arena consilium capit",
			i: 0},
			want: 'G'},

		{args: args{s: "\033[38;2;160;10;10mGladiator in arena consilium capit",
			i: 0},
			want: 'G'},

		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit",
			i: 0},
			want: 'G'},

		{args: args{s: "Gladiator in arena consilium capit",
			i: 1},
			want: 'l'},

		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit",
			i: 1},
			want: 'l'},

		{args: args{s: "Gl\033[38;2;160;10;10madiator in arena consilium capit",
			i: 1},
			want: 'l'},

		{args: args{s: "Gladiator in arena consilium capit",
			i: 32},
			want: 'i'},

		{args: args{s: "Gladiator in arena consilium cap\033[38;2;160;10;10mit",
			i: 32},
			want: 'i'},

		{args: args{s: "Gladiator in arena consilium capi\033[38;2;160;10;10mt",
			i: 32},
			want: 'i'},

		{args: args{s: "Gladiator in arena consilium capit",
			i: 33},
			want: 't'},

		{args: args{s: "Gladiator in arena consilium capi\033[38;2;160;10;10mt",
			i: 33},
			want: 't'},

		{args: args{s: "Gladiator in arena consilium capit\033[38;2;160;10;10m",
			i: 33},
			want: 't'},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRune(tt.args.s, tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRune() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getRune() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_insertRune(t *testing.T) {
	type args struct {
		s string
		i int
		r rune
	}
	tests := []struct {
		name string
		args args
		want string
	}{

		// tests with empty strings
		{args: args{s: "",
			i: 0, r: '🏯'},
			want: ""},

		{args: args{s: "R",
			i: 0, r: '🏯'},
			want: "🏯"},

		// tests with strings that contain no ANSI color codes
		{args: args{s: "Gladiator in━arena consilium capit",
			i: 0, r: '🏯'},
			want: "🏯ladiator in━arena consilium capit"},

		{args: args{s: "Gladiator in━arena consilium capit",
			i: 20, r: '🏯'},
			want: "Gladiator in━arena🏯consilium capit"},

		{args: args{s: "Gladiator in━arena consilium capit",
			i: 35, r: '🏯'},
			want: "Gladiator in━arena consilium capi🏯"},

		{args: args{s: "Gladiator in━arena consilium capit",
			i: 36, r: '🏯'},
			want: "Gladiator in━arena consilium capit"},

		// tests with strings that contain one ANSI color code
		{args: args{s: "\033[38;2;160;10;10mGladiator in━arena consilium capit",
			i: 0, r: '🏯'},
			want: "🏯[38;2;160;10;10mGladiator in━arena consilium capit"},

		{args: args{s: "\033[38;2;160;10;10mGladiator in━arena consilium capit",
			i: 17, r: '🏯'},
			want: "\033[38;2;160;10;10m🏯ladiator in━arena consilium capit"},

		{args: args{s: "\033[38;2;160;10;10mGladiator in━arena consilium capit",
			i: 37, r: '🏯'},
			want: "\033[38;2;160;10;10mGladiator in━arena🏯consilium capit"},

		{args: args{s: "\033[38;2;160;10;10mGladiator in━arena consilium capit",
			i: 52, r: '🏯'},
			want: "\033[38;2;160;10;10mGladiator in━arena consilium capi🏯"},

		{args: args{s: "\033[38;2;160;10;10mGladiator in━arena consilium capit",
			i: 53, r: '🏯'},
			want: "\033[38;2;160;10;10mGladiator in━arena consilium capit"},

		// tests with strings that contain two ANSI color codes
		{args: args{s: "\033[38;2;160;10;10mGladiator\033[0m in━arena consilium capit",
			i: 0, r: '🏯'},
			want: "🏯[38;2;160;10;10mGladiator\033[0m in━arena consilium capit"},

		{args: args{s: "\033[38;2;160;10;10mGladiator\033[0m in━arena consilium capit",
			i: 17, r: '🏯'},
			want: "\033[38;2;160;10;10m🏯ladiator\033[0m in━arena consilium capit"},

		{args: args{s: "\033[38;2;160;10;10mGladiator\033[0m in━arena consilium capit",
			i: 41, r: '🏯'},
			want: "\033[38;2;160;10;10mGladiator\033[0m in━arena🏯consilium capit"},

		{args: args{s: "\033[38;2;160;10;10mGladiator\033[0m in━arena consilium capit",
			i: 56, r: '🏯'},
			want: "\033[38;2;160;10;10mGladiator\033[0m in━arena consilium capi🏯"},

		{args: args{s: "\033[38;2;160;10;10mGladiator\033[0m in━arena consilium capit",
			i: 57, r: '🏯'},
			want: "\033[38;2;160;10;10mGladiator\033[0m in━arena consilium capit"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := insertRune(tt.args.s, tt.args.i, tt.args.r); got != tt.want {
				t.Errorf("insertRune() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_logicalToPhysical(t *testing.T) {
	type args struct {
		s  string
		li int
	}
	tests := []struct {
		name             string
		args             args
		wantPi1, wantPi2 int
		wantSout         string
	}{

		// logicalToPhysical can be invoked with either "force=true" or
		// "force=false". In case force is given, the input string is extended
		// with as many blank characters as necessary as to have exactly a
		// number of printable characters equal to the logical position seeked.
		// logicalToPhysical returns then the physical position of the logical
		// location requested and also the resulting string which might have
		// been modified or not.
		//
		// Thus, wantPi1 is the physical location obtained with "force=False".
		// The resulting string is not checked as it is never modified
		//
		//
		// wantPi2 and wantSout are the returned values when "force=true"
		//
		// Instead of creating different tests for these two cases, the body
		// below the definition of all test cases invoke each test case both
		// with force=true and force=false

		// trivial example with the empty string
		{args: args{s: "",
			li: -1},
			wantPi1:  -1,
			wantPi2:  -1,
			wantSout: ""},

		{args: args{s: "",
			li: 0},
			wantPi1:  -1,
			wantPi2:  0,
			wantSout: " "},

		{args: args{s: "",
			li: 1},
			wantPi1:  -1,
			wantPi2:  1,
			wantSout: "  "},

		// examples with no ANSI color codes
		{args: args{s: "Gladiator━in arena consilium capit",
			li: -1},
			wantPi1:  -1,
			wantPi2:  -1,
			wantSout: "Gladiator━in arena consilium capit"},

		{args: args{s: "Gladiator━in arena consilium capit",
			li: 0},
			wantPi1:  0,
			wantPi2:  0,
			wantSout: "Gladiator━in arena consilium capit"},

		{args: args{s: "Gladiator━in arena consilium capit",
			li: 16},
			wantPi1:  18,
			wantPi2:  18,
			wantSout: "Gladiator━in arena consilium capit"},

		{args: args{s: "Gladiator━in arena consilium capit",
			li: 33},
			wantPi1:  35,
			wantPi2:  35,
			wantSout: "Gladiator━in arena consilium capit"},

		{args: args{s: "Gladiator━in arena consilium capit",
			li: 34},
			wantPi1:  -1,
			wantPi2:  36,
			wantSout: "Gladiator━in arena consilium capit "},

		// examples with only one ANSI color code
		{args: args{s: "\033[38;2;160;10;10mGladiator━in arena consilium capit",
			li: -1},
			wantPi1:  -1,
			wantPi2:  -1,
			wantSout: "\033[38;2;160;10;10mGladiator━in arena consilium capit"},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit",
			li: -1},
			wantPi1:  -1,
			wantPi2:  -1,
			wantSout: "G\033[38;2;160;10;10mladiator━in arena consilium capit"},

		{args: args{s: "\033[38;2;160;10;10mGladiator━in arena consilium capit",
			li: 0},
			wantPi1:  17,
			wantPi2:  17,
			wantSout: "\033[38;2;160;10;10mGladiator━in arena consilium capit"},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit",
			li: 0},
			wantPi1:  0,
			wantPi2:  0,
			wantSout: "G\033[38;2;160;10;10mladiator━in arena consilium capit"},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit",
			li: 1},
			wantPi1:  18,
			wantPi2:  18,
			wantSout: "G\033[38;2;160;10;10mladiator━in arena consilium capit"},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit",
			li: 9},
			wantPi1:  26,
			wantPi2:  26,
			wantSout: "G\033[38;2;160;10;10mladiator━in arena consilium capit"},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit",
			li: 33},
			wantPi1:  52,
			wantPi2:  52,
			wantSout: "G\033[38;2;160;10;10mladiator━in arena consilium capit"},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit",
			li: 34},
			wantPi1:  -1,
			wantPi2:  53,
			wantSout: "G\033[38;2;160;10;10mladiator━in arena consilium capit "},

		// examples with two ANSI color codes
		{args: args{s: "\033[38;2;160;10;10mGladiator━\033[0min arena consilium capit",
			li: -1},
			wantPi1:  -1,
			wantPi2:  -1,
			wantSout: "\033[38;2;160;10;10mGladiator━\033[0min arena consilium capit"},

		{args: args{s: "G\033[38;2;160;10;10mladiator━\033[0min arena consilium capit",
			li: -1},
			wantPi1:  -1,
			wantPi2:  -1,
			wantSout: "G\033[38;2;160;10;10mladiator━\033[0min arena consilium capit"},

		{args: args{s: "\033[38;2;160;10;10mGladiator━\033[0min arena consilium capit",
			li: 0},
			wantPi1:  17,
			wantPi2:  17,
			wantSout: "\033[38;2;160;10;10mGladiator━\033[0min arena consilium capit"},

		{args: args{s: "G\033[38;2;160;10;10mladiator━\033[0min arena consilium capit",
			li: 0},
			wantPi1:  0,
			wantPi2:  0,
			wantSout: "G\033[38;2;160;10;10mladiator━\033[0min arena consilium capit"},

		{args: args{s: "G\033[38;2;160;10;10mladiator━\033[0min arena consilium capit",
			li: 1},
			wantPi1:  18,
			wantPi2:  18,
			wantSout: "G\033[38;2;160;10;10mladiator━\033[0min arena consilium capit"},

		{args: args{s: "G\033[38;2;160;10;10mladiator━\033[0min arena consilium capit",
			li: 9},
			wantPi1:  26,
			wantPi2:  26,
			wantSout: "G\033[38;2;160;10;10mladiator━\033[0min arena consilium capit"},

		{args: args{s: "G\033[38;2;160;10;10mladiator━\033[0min arena consilium capit",
			li: 10},
			wantPi1:  33,
			wantPi2:  33,
			wantSout: "G\033[38;2;160;10;10mladiator━\033[0min arena consilium capit"},

		{args: args{s: "G\033[38;2;160;10;10mladiator━\033[0min arena consilium capit",
			li: 33},
			wantPi1:  56,
			wantPi2:  56,
			wantSout: "G\033[38;2;160;10;10mladiator━\033[0min arena consilium capit"},

		{args: args{s: "G\033[38;2;160;10;10mladiator━\033[0min arena consilium capit",
			li: 34},
			wantPi1:  -1,
			wantPi2:  57,
			wantSout: "G\033[38;2;160;10;10mladiator━\033[0min arena consilium capit "},

		{args: args{s: "\033[38;2;160;10;10mGladiator━in arena consilium capit\033[0m",
			li: -1},
			wantPi1:  -1,
			wantPi2:  -1,
			wantSout: "\033[38;2;160;10;10mGladiator━in arena consilium capit\033[0m"},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m",
			li: -1},
			wantPi1:  -1,
			wantPi2:  -1,
			wantSout: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m"},

		{args: args{s: "\033[38;2;160;10;10mGladiator━in arena consilium capit\033[0m",
			li: 0},
			wantPi1:  17,
			wantPi2:  17,
			wantSout: "\033[38;2;160;10;10mGladiator━in arena consilium capit\033[0m"},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m",
			li: 0},
			wantPi1:  0,
			wantPi2:  0,
			wantSout: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m"},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m",
			li: 1},
			wantPi1:  18,
			wantPi2:  18,
			wantSout: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m"},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m",
			li: 9},
			wantPi1:  26,
			wantPi2:  26,
			wantSout: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m"},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m",
			li: 14},
			wantPi1:  33,
			wantPi2:  33,
			wantSout: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m"},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m",
			li: 33},
			wantPi1:  52,
			wantPi2:  52,
			wantSout: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m"},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m",
			li: 34},
			wantPi1:  -1,
			wantPi2:  57,
			wantSout: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m "},
	}
	for _, tt := range tests {

		// Tests run without modifying the input string
		t.Run(tt.name, func(t *testing.T) {
			if gotPi, _ := logicalToPhysical(tt.args.s, tt.args.li, false); gotPi != tt.wantPi1 {
				t.Errorf("[force = %v] logicalToPhysical() = %v, want %v", false, gotPi, tt.wantPi1)
			}
		})

		// Tests run modifying the input string if needed
		t.Run(tt.name, func(t *testing.T) {
			if gotPi, gotSout := logicalToPhysical(tt.args.s, tt.args.li, true); gotPi != tt.wantPi2 || gotSout != tt.wantSout {
				t.Errorf("[force = %v] logicalToPhysical() = %v/'%v', want %v/'%v'", true, gotPi, gotSout, tt.wantPi2, tt.wantSout)
			}
		})
	}
}

func Test_getSingleSplitter(t *testing.T) {
	type args struct {
		west  rune
		east  rune
		north rune
		south rune
	}
	tests := []struct {
		name string
		args args
		want rune
	}{

		// Only the fully supported cases with full rules are tested. Other
		// weird combinations that include the double rule (either horizontal or
		// vertical) or the usage of clines is not tested

		// horizontal single / vertical single
		{args: args{west: none, east: '─', north: none, south: '│'},
			want: '┌'},

		{args: args{west: '─', east: '─', north: none, south: '│'},
			want: '┬'},

		{args: args{west: '─', east: none, north: none, south: '│'},
			want: '┐'},

		{args: args{west: none, east: '─', north: '│', south: '│'},
			want: '├'},

		{args: args{west: '─', east: '─', north: '│', south: '│'},
			want: '┼'},

		{args: args{west: '─', east: none, north: '│', south: '│'},
			want: '┤'},

		{args: args{west: none, east: '─', north: '│', south: none},
			want: '└'},

		{args: args{west: '─', east: '─', north: '│', south: none},
			want: '┴'},

		{args: args{west: '─', east: none, north: '│', south: none},
			want: '┘'},

		// horizontal double / vertical single
		{args: args{west: none, east: '═', north: none, south: '│'},
			want: '╒'},

		{args: args{west: '═', east: '═', north: none, south: '│'},
			want: '╤'},

		{args: args{west: '═', east: none, north: none, south: '│'},
			want: '╕'},

		{args: args{west: none, east: '═', north: '│', south: '│'},
			want: '╞'},

		{args: args{west: '═', east: '═', north: '│', south: '│'},
			want: '╪'},

		{args: args{west: '═', east: none, north: '│', south: '│'},
			want: '╡'},

		{args: args{west: none, east: '═', north: '│', south: none},
			want: '╘'},

		{args: args{west: '═', east: '═', north: '│', south: none},
			want: '\u2567'},

		{args: args{west: '═', east: none, north: '│', south: none},
			want: '╛'},

		// horizontal thick / vertical single
		{args: args{west: none, east: '━', north: none, south: '│'},
			want: '┍'},

		{args: args{west: '━', east: '━', north: none, south: '│'},
			want: '┯'},

		{args: args{west: '━', east: none, north: none, south: '│'},
			want: '┑'},

		{args: args{west: none, east: '━', north: '│', south: '│'},
			want: '┝'},

		{args: args{west: '━', east: '━', north: '│', south: '│'},
			want: '┿'},

		{args: args{west: '━', east: none, north: '│', south: '│'},
			want: '┥'},

		{args: args{west: none, east: '━', north: '│', south: none},
			want: '┕'},

		{args: args{west: '━', east: '━', north: '│', south: none},
			want: '┷'},

		{args: args{west: '━', east: none, north: '│', south: none},
			want: '┙'},

		// horizontal single / vertical double
		{args: args{west: none, east: '─', north: none, south: '║'},
			want: '╓'},

		{args: args{west: '─', east: '─', north: none, south: '║'},
			want: '╥'},

		{args: args{west: '─', east: none, north: none, south: '║'},
			want: '╖'},

		{args: args{west: none, east: '─', north: '║', south: '║'},
			want: '╟'},

		{args: args{west: '─', east: '─', north: '║', south: '║'},
			want: '╫'},

		{args: args{west: '─', east: none, north: '║', south: '║'},
			want: '╢'},

		{args: args{west: none, east: '─', north: '║', south: none},
			want: '╙'},

		{args: args{west: '─', east: '─', north: '║', south: none},
			want: '\u2568'},

		{args: args{west: '─', east: none, north: '║', south: none},
			want: '╜'},

		// horizontal double / vertical double
		{args: args{west: none, east: '═', north: none, south: '║'},
			want: '╔'},

		{args: args{west: '═', east: '═', north: none, south: '║'},
			want: '╦'},

		{args: args{west: '═', east: none, north: none, south: '║'},
			want: '╗'},

		{args: args{west: none, east: '═', north: '║', south: '║'},
			want: '╠'},

		{args: args{west: '═', east: '═', north: '║', south: '║'},
			want: '╬'},

		{args: args{west: '═', east: none, north: '║', south: '║'},
			want: '╣'},

		{args: args{west: none, east: '═', north: '║', south: none},
			want: '╚'},

		{args: args{west: '═', east: '═', north: '║', south: none},
			want: '╩'},

		{args: args{west: '═', east: none, north: '║', south: none},
			want: '╝'},

		// horizontal single / vertical thick
		{args: args{west: none, east: '─', north: none, south: '┃'},
			want: '┎'},

		{args: args{west: '─', east: '─', north: none, south: '┃'},
			want: '┰'},

		{args: args{west: '─', east: none, north: none, south: '┃'},
			want: '┒'},

		{args: args{west: none, east: '─', north: '┃', south: '┃'},
			want: '┠'},

		{args: args{west: '─', east: '─', north: '┃', south: '┃'},
			want: '╂'},

		{args: args{west: '─', east: none, north: '┃', south: '┃'},
			want: '┨'},

		{args: args{west: none, east: '─', north: '┃', south: none},
			want: '┖'},

		{args: args{west: '─', east: '─', north: '┃', south: none},
			want: '┸'},

		{args: args{west: '─', east: none, north: '┃', south: none},
			want: '┚'},

		// horizontal thick / vertical thick
		{args: args{west: none, east: '━', north: none, south: '┃'},
			want: '┏'},

		{args: args{west: '━', east: '━', north: none, south: '┃'},
			want: '┳'},

		{args: args{west: '━', east: none, north: none, south: '┃'},
			want: '┓'},

		{args: args{west: none, east: '━', north: '┃', south: '┃'},
			want: '┣'},

		{args: args{west: '━', east: '━', north: '┃', south: '┃'},
			want: '╋'},

		{args: args{west: '━', east: none, north: '┃', south: '┃'},
			want: '┫'},

		{args: args{west: none, east: '━', north: '┃', south: none},
			want: '┗'},

		{args: args{west: '━', east: '━', north: '┃', south: none},
			want: '┻'},

		{args: args{west: '━', east: none, north: '┃', south: none},
			want: '┛'},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSingleSplitter(tt.args.west, tt.args.east, tt.args.north, tt.args.south); got != tt.want {
				t.Errorf("getSingleSplitter() = %v, want %v", got, tt.want)
			}
		})
	}
}
