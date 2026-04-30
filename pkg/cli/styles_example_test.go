package cli

import core "dappco.re/go"

func ExampleDefaultTableStyle() {
	core.Println("DefaultTableStyle")
	// Output: DefaultTableStyle
}

func ExampleFormatAge() {
	core.Println("FormatAge")
	// Output: FormatAge
}

func ExampleNewTable() {
	core.Println("NewTable")
	// Output: NewTable
}

func ExamplePad() {
	core.Println("Pad")
	// Output: Pad
}

func ExampleTable_AddRow() {
	core.Println("Table_AddRow")
	// Output: Table_AddRow
}

func ExampleTable_Render() {
	core.Println("Table_Render")
	// Output: Table_Render
}

func ExampleTable_String() {
	core.Println("Table_String")
	// Output: Table_String
}

func ExampleTable_WithBorders() {
	core.Println("Table_WithBorders")
	// Output: Table_WithBorders
}

func ExampleTable_WithCellStyle() {
	core.Println("Table_WithCellStyle")
	// Output: Table_WithCellStyle
}

func ExampleTable_WithMaxWidth() {
	core.Println("Table_WithMaxWidth")
	// Output: Table_WithMaxWidth
}

func ExampleTruncate() {
	core.Println(Truncate("abcdef", 4))
	// Output: a...
}
