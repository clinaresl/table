package table

import (
	"reflect"
	"testing"
)

func Test_content_Process(t *testing.T) {
	type args struct {
		col column
	}
	tests := []struct {
		name string
		c    content
		args args
		want []string
	}{

		{c: "",
			args: args{col: column{}},
			want: []string{""}},

		{c: "Black lives matter",
			args: args{col: column{}},
			want: []string{"Black lives matter"}},

		{c: "Black\nlives\nmatter",
			args: args{col: column{}},
			want: []string{"Black", "lives", "matter"}},

		{c: "",
			args: args{col: column{sep: "|", hformat: style{alignment: 'l'}}},
			want: []string{""}},

		{c: "",
			args: args{col: column{sep: "|", hformat: style{alignment: 'c'}}},
			want: []string{""}},

		{c: "",
			args: args{col: column{sep: "|", hformat: style{alignment: 'r'}}},
			want: []string{""}},

		{c: "Black lives matter",
			args: args{col: column{sep: "|", hformat: style{alignment: 'l'}}},
			want: []string{"Black lives matter"}},

		{c: "Black lives matter",
			args: args{col: column{sep: "|", hformat: style{alignment: 'c'}}},
			want: []string{"Black lives matter"}},

		{c: "Black lives matter",
			args: args{col: column{sep: "|", hformat: style{alignment: 'r'}}},
			want: []string{"Black lives matter"}},

		{c: "Black\nlives\nmatter",
			args: args{col: column{sep: "|", hformat: style{alignment: 'l'}}},
			want: []string{"Black", "lives", "matter"}},

		{c: "Black\nlives\nmatter",
			args: args{col: column{sep: "|", hformat: style{alignment: 'c'}}},
			want: []string{"Black", "lives", "matter"}},

		{c: "Black\nlives\nmatter",
			args: args{col: column{sep: "|", hformat: style{alignment: 'r'}}},
			want: []string{"Black", "lives", "matter"}},

		{c: "Black\nlives\nmatter\n",
			args: args{col: column{sep: "|", hformat: style{alignment: 'l'}}},
			want: []string{"Black", "lives", "matter", ""}},

		{c: "Black\nlives\nmatter\n",
			args: args{col: column{sep: "|", hformat: style{alignment: 'c'}}},
			want: []string{"Black", "lives", "matter", ""}},

		{c: "Black\nlives\nmatter\n",
			args: args{col: column{sep: "|", hformat: style{alignment: 'r'}}},
			want: []string{"Black", "lives", "matter", ""}},

		{c: "Black\n\nlives\nmatter",
			args: args{col: column{sep: "|", hformat: style{alignment: 'l'}}},
			want: []string{"Black", "", "lives", "matter"}},

		{c: "Black\nlives\n\nmatter",
			args: args{col: column{sep: "|", hformat: style{alignment: 'c'}}},
			want: []string{"Black", "lives", "", "matter"}},

		{c: "Black\nlives\nmatter\n\nAlways",
			args: args{col: column{sep: "|", hformat: style{alignment: 'r'}}},
			want: []string{"Black", "lives", "matter", "", "Always"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Process(tt.args.col, 0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("content.Process() = %v, want %v", got, tt.want)
			}
		})
	}
}
