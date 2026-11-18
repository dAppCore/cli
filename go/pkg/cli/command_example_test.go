package cli

import core "github.com/dappcore/go"

func ExampleRegisterCommand() {
	core.Println("RegisterCommand")
	// Output: RegisterCommand
}

func ExampleRequireArgs() {
	core.Println("RequireArgs")
	// Output: RequireArgs
}

func ExampleRequireExactArgs() {
	core.Println("RequireExactArgs")
	// Output: RequireExactArgs
}
