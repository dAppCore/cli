package cli

import core "github.com/dappcore/go"

func ExampleT() {
	core.Println(T("example.key"))
	// Output: example.key
}
