package cli

import core "dappco.re/go"

func ExampleAnsiStyle_Background() {
	core.Println("AnsiStyle_Background")
	// Output: AnsiStyle_Background
}

func ExampleAnsiStyle_Italic() {
	core.Println("AnsiStyle_Italic")
	// Output: AnsiStyle_Italic
}

func ExampleAnsiStyle_Underline() {
	core.Println("AnsiStyle_Underline")
	// Output: AnsiStyle_Underline
}

func ExampleColorEnabled() {
	old := ColorEnabled()
	defer SetColorEnabled(old)
	SetColorEnabled(false)
	core.Println(ColorEnabled())
	// Output: false
}

func ExampleSetColorEnabled() {
	old := ColorEnabled()
	defer SetColorEnabled(old)
	SetColorEnabled(false)
	core.Println(NewStyle().Bold().Render("ready"))
	// Output: ready
}

func ExampleNewStyle() {
	old := ColorEnabled()
	defer SetColorEnabled(old)
	SetColorEnabled(false)
	core.Println(NewStyle().Bold().Render("ready"))
	// Output: ready
}

func ExampleAnsiStyle_Bold() {
	old := ColorEnabled()
	defer SetColorEnabled(old)
	SetColorEnabled(false)
	core.Println(NewStyle().Bold().Render("ready"))
	// Output: ready
}

func ExampleAnsiStyle_Dim() {
	old := ColorEnabled()
	defer SetColorEnabled(old)
	SetColorEnabled(false)
	core.Println(NewStyle().Dim().Render("ready"))
	// Output: ready
}

func ExampleAnsiStyle_Foreground() {
	old := ColorEnabled()
	defer SetColorEnabled(old)
	SetColorEnabled(false)
	core.Println(NewStyle().Foreground("#00ff00").Render("ready"))
	// Output: ready
}

func ExampleAnsiStyle_Render() {
	old := ColorEnabled()
	defer SetColorEnabled(old)
	SetColorEnabled(false)
	core.Println(NewStyle().Bold().Render("ready"))
	// Output: ready
}
