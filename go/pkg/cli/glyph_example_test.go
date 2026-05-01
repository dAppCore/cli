package cli

import core "dappco.re/go"

func ExampleUseASCII() {
	core.Println("UseASCII")
	// Output: UseASCII
}

func ExampleUseEmoji() {
	core.Println("UseEmoji")
	// Output: UseEmoji
}

func ExampleUseUnicode() {
	core.Println("UseUnicode")
	// Output: UseUnicode
}

func ExampleGlyph() {
	UseASCII()
	defer UseUnicode()
	core.Println(Glyph(":check:"))
	// Output: [OK]
}
