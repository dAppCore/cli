package cli

import core "dappco.re/go"

func ExampleComposite_Regions() {
	composite := Layout("HC")
	count := 0
	for range composite.Regions() {
		count++
	}
	core.Println(count)
	// Output: 2
}

func ExampleComposite_Slots() {
	composite := Layout("H")
	count := 0
	for range composite.Slots() {
		count++
	}
	core.Println(count)
	// Output: 1
}

func ExampleStringBlock_Render() {
	core.Println(StringBlock("ready").Render())
	// Output: ready
}

func ExampleLayout() {
	composite := Layout("C").C("body")
	core.Println(core.Trim(composite.String()))
	// Output: body
}

func ExampleParseVariant() {
	result := ParseVariant("HCF")
	core.Println(result.OK)
	// Output: true
}

func ExampleComposite_H() {
	composite := Layout("H").H("header")
	core.Println(core.Trim(composite.String()))
	// Output: header
}

func ExampleComposite_L() {
	composite := Layout("L").L("left")
	core.Println(core.Trim(composite.String()))
	// Output: left
}

func ExampleComposite_C() {
	composite := Layout("C").C("content")
	core.Println(core.Trim(composite.String()))
	// Output: content
}

func ExampleComposite_R() {
	composite := Layout("R").R("right")
	core.Println(core.Trim(composite.String()))
	// Output: right
}

func ExampleComposite_F() {
	composite := Layout("F").F("footer")
	core.Println(core.Trim(composite.String()))
	// Output: footer
}
