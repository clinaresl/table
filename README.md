# Introduction

This package implements means for drawing data in tabular form and it is
intended as a substitution of the Go standard package `tabwriter`. Its design is
based on the functionality of tables in LaTeX but extends its functionality in
various ways through a very simple interface

It honours [UTF-8 characters](https://www.utf8-chartable.de/), [ANSI color escape sequences](https://stackoverflow.com/questions/4842424/list-of-ansi-color-escape-sequences), full/partial
horizontal rules, various vertical and horizontal alignment options, and
multicolumns.

Remarkably, it prints any *stringer* and as tables are stringers, tables can be
nested to any degree.


# Installation 

Clone and install the `table` package with the following command:

    $ go get github.com/clinaresl/table
    
To try the different examples given in the package change dir to
`$GOPATH/github.com/clinaresl/table` and type:

    $ go test
   
# Usage #

This section provides various examples of usage with the hope of providing a
flavour of the different capabilities of the package. For a full description of
the package check out the technical documentation.

## First step: Create a table ##

Before inserting data to a new table it is necessary to create it first:

```Go
	t, err := NewTable("l   p{25}")
	if err != nil {
		log.Fatalln(" NewTable: Fatal error!")
	}
```

This snippet creates a table with two columns. The first one displays its
contents ragged left, whereas the second one takes a fixed width of 25
characters to display the contents of each cell and, in case a cell exceeds the
available width, its contents are shown left ragged in as many lines as needed.
In case it was not possible to successfully process the *column specification*,
an error is immediately returned.

The following table shows all the different options for specifying the format of
a single column:

| Syntax | Purpose |
|:------:|:-------:|
|  `l`   | the contents of the column are ragged left |
|  `c`   | the contents of the column are horizontally aligned |
|  `r`   | the contents of the column are ragged right |
|  `p{NUMBER}` | the cell takes a fixed with equal to *NUMBER** characters and the contents are split across various lines if needed |
|  `L{NUMBER}` | the width of the column does not exceed *NUMBER** characters and the contents are ragged left |
|  `C{NUMBER}` | the width of the column does not exceed *NUMBER** characters and the contents are centered |
|  `R{NUMBER}` | the width of the column does not exceed *NUMBER** characters and the contents are ragged right |

The *column specification* allows the usage of `|`, e.g.:

``` Go
	t, _ := NewTable("|c|c|c|c|c|")
```

creates a table with five different columns all separated by a single vertical
separator. It is possible to create also *double* and *thick* vertical
separators using `||` and `|||` respectively. It is also possible to provide any
other character (e.g., blank spaces or tabs) either before or after any column.
These are then copied either before or after the contents of each cell in each
row.

In case a second string is given to `NewTable` it is interpreted as the *row
specification*:

```Go
	t, _ := NewTable("| c | c  c |", "cct")
```

This line (where no error checking is performed!) creates three different
columns whose contents are horizontally centered surrounded by a single space
and with vertical single separators between adjacent columns and before and
after the first and last column. In addition, it sets the *vertical alignment*
of each cell as follows: the contents of the first and second columns are
vertically centered (`c`), whereas the contents of the last column are pushed to
the top of the cell ---`t`.

| Syntax | Purpose |
|:------:|:-------:|
|  `t`   | the contents of the column are aligned to the top |
|  `c`   | the contents of the column are vertically aligned |
|  `b`   | the contents of the column are aligned to the bottom |

By default, all columns are vertically aligned to the top. In case a *row
specification* is given it must refer to as many columns as there are in the
*column specification* given first or less. In contraposition to the *column
specification*, the *row specification* can only consist of any of the modifies
shown above.
   
`NewTable` returns a pointer to `Table` which can be used next for adding data
to it and, in the end, printing it.
 
# License #

table is free software: you can redistribute it and/or modify it
under the terms of the GNU General Public License as published by the
Free Software Foundation, either version 3 of the License, or (at your
option) any later version.

table is distributed in the hope that it will be useful, but WITHOUT
ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License
for more details.

You should have received a copy of the GNU General Public License
along with table.  If not, see <http://www.gnu.org/licenses/>.


# Author #

Carlos Linares Lopez <carlos.linares@uc3m.es>  
Computer Science Department <https://www.inf.uc3m.es/en>  
Universidad Carlos III de Madrid <https://www.uc3m.es/home>
