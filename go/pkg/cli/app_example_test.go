package cli

import core "github.com/dappcore/go"

func ExampleMain() {
	core.Println("Main")
	// Output: Main
}

func ExampleMainWithLocales() {
	core.Println("MainWithLocales")
	// Output: MainWithLocales
}

func ExampleSemVer() {
	core.Println("SemVer")
	// Output: SemVer
}

func ExampleWithAppName() {
	core.Println("WithAppName")
	// Output: WithAppName
}

func ExampleWithLocales() {
	core.Println("WithLocales")
	// Output: WithLocales
}
