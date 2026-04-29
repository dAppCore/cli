package cli

import core "dappco.re/go"

func ExampleT() {
	core.Println(T("example.key"))
	// Output: example.key
}
