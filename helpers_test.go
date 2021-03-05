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

func Test_physicalToLogical(t *testing.T) {
	type args struct {
		s  string
		pi int
	}
	tests := []struct {
		name   string
		args   args
		wantLi int
	}{

		// trivial example with the empty string
		{args: args{s: "",
			pi: -1},
			wantLi: -1},

		{args: args{s: "",
			pi: 0},
			wantLi: -1},

		{args: args{s: "",
			pi: 1},
			wantLi: -1},

		// examples with no ANSI color codes
		{args: args{s: "Gladiator in arena consilium capit",
			pi: -1},
			wantLi: -1},

		{args: args{s: "Gladiator in arena consilium capit",
			pi: 0},
			wantLi: 0},

		{args: args{s: "Gladiator in arena consilium capit",
			pi: 16},
			wantLi: 16},

		{args: args{s: "Gladiator in arena consilium capit",
			pi: 33},
			wantLi: 33},

		{args: args{s: "Gladiator in arena consilium capit",
			pi: 34},
			wantLi: -1},

		// examples with only one ANSI color code
		{args: args{s: "\033[38;2;160;10;10mGladiator in arena consilium capit",
			pi: -1},
			wantLi: -1},

		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit",
			pi: -1},
			wantLi: -1},

		{args: args{s: "\033[38;2;160;10;10mGladiator in arena consilium capit",
			pi: 17},
			wantLi: 0},

		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit",
			pi: 0},
			wantLi: 0},

		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit",
			pi: 18},
			wantLi: 1},

		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit",
			pi: 26},
			wantLi: 9},

		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit",
			pi: 50},
			wantLi: 33},

		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit",
			pi: 51},
			wantLi: -1},

		// examples with two ANSI color codes
		{args: args{s: "\033[38;2;160;10;10mGladiator \033[0min arena consilium capit",
			pi: -1},
			wantLi: -1},

		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
			pi: -1},
			wantLi: -1},

		{args: args{s: "\033[38;2;160;10;10mGladiator \033[0min arena consilium capit",
			pi: 17},
			wantLi: 0},

		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
			pi: 0},
			wantLi: 0},

		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
			pi: 18},
			wantLi: 1},

		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
			pi: 26},
			wantLi: 9},

		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
			pi: 31},
			wantLi: 10},

		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
			pi: 54},
			wantLi: 33},

		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
			pi: 55},
			wantLi: -1},

		{args: args{s: "\033[38;2;160;10;10mGladiator in arena consilium capit\033[0m",
			pi: -1},
			wantLi: -1},

		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit\033[0m",
			pi: -1},
			wantLi: -1},

		{args: args{s: "\033[38;2;160;10;10mGladiator in arena consilium capit\033[0m",
			pi: 17},
			wantLi: 0},

		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit\033[0m",
			pi: 0},
			wantLi: 0},

		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit\033[0m",
			pi: 18},
			wantLi: 1},

		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit\033[0m",
			pi: 26},
			wantLi: 9},

		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit\033[0m",
			pi: 31},
			wantLi: 14},

		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit\033[0m",
			pi: 50},
			wantLi: 33},

		{args: args{s: "G\033[38;2;160;10;10mladiator in arena consilium capit\033[0m",
			pi: 51},
			wantLi: -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLi := physicalToLogical(tt.args.s, tt.args.pi); gotLi != tt.wantLi {
				t.Errorf("physicalToLogical() = %v, want %v", gotLi, tt.wantLi)
			}
		})
	}
}

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
		name   string
		args   args
		wantPi int
	}{
		// trivial example with the empty string
		{args: args{s: "",
			li: -1},
			wantPi: -1},

		{args: args{s: "",
			li: 0},
			wantPi: -1},

		{args: args{s: "",
			li: 1},
			wantPi: -1},

		// examples with no ANSI color codes
		{args: args{s: "Gladiator━in arena consilium capit",
			li: -1},
			wantPi: -1},

		{args: args{s: "Gladiator━in arena consilium capit",
			li: 0},
			wantPi: 0},

		{args: args{s: "Gladiator━in arena consilium capit",
			li: 16},
			wantPi: 16},

		{args: args{s: "Gladiator━in arena consilium capit",
			li: 33},
			wantPi: 33},

		{args: args{s: "Gladiator━in arena consilium capit",
			li: 34},
			wantPi: -1},

		// examples with only one ANSI color code
		{args: args{s: "\033[38;2;160;10;10mGladiator━in arena consilium capit",
			li: -1},
			wantPi: -1},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit",
			li: -1},
			wantPi: -1},

		{args: args{s: "\033[38;2;160;10;10mGladiator━in arena consilium capit",
			li: 0},
			wantPi: 17},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit",
			li: 0},
			wantPi: 0},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit",
			li: 1},
			wantPi: 18},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit",
			li: 9},
			wantPi: 26},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit",
			li: 33},
			wantPi: 50},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit",
			li: 34},
			wantPi: -1},

		// examples with two ANSI color codes
		{args: args{s: "\033[38;2;160;10;10mGladiator \033[0min arena consilium capit",
			li: -1},
			wantPi: -1},

		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
			li: -1},
			wantPi: -1},

		{args: args{s: "\033[38;2;160;10;10mGladiator \033[0min arena consilium capit",
			li: 0},
			wantPi: 17},

		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
			li: 0},
			wantPi: 0},

		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
			li: 1},
			wantPi: 18},

		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
			li: 9},
			wantPi: 26},

		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
			li: 10},
			wantPi: 31},

		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
			li: 33},
			wantPi: 54},

		{args: args{s: "G\033[38;2;160;10;10mladiator \033[0min arena consilium capit",
			li: 34},
			wantPi: -1},

		{args: args{s: "\033[38;2;160;10;10mGladiator━in arena consilium capit\033[0m",
			li: -1},
			wantPi: -1},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m",
			li: -1},
			wantPi: -1},

		{args: args{s: "\033[38;2;160;10;10mGladiator━in arena consilium capit\033[0m",
			li: 0},
			wantPi: 17},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m",
			li: 0},
			wantPi: 0},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m",
			li: 1},
			wantPi: 18},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m",
			li: 9},
			wantPi: 26},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m",
			li: 14},
			wantPi: 31},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m",
			li: 33},
			wantPi: 50},

		{args: args{s: "G\033[38;2;160;10;10mladiator━in arena consilium capit\033[0m",
			li: 34},
			wantPi: -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPi := logicalToPhysical(tt.args.s, tt.args.li); gotPi != tt.wantPi {
				t.Errorf("logicalToPhysical() = %v, want %v", gotPi, tt.wantPi)
			}
		})
	}
}
