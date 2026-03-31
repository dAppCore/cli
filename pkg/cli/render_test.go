package cli

import (
	"strings"
	"testing"
)

func TestCompositeRender_Good(t *testing.T) {
	UseRenderFlat()
	composite := Layout("HCF")
	composite.H("Header content").C("Body content").F("Footer content")

	output := composite.String()
	if !strings.Contains(output, "Header content") {
		t.Errorf("Render flat: expected 'Header content' in output, got %q", output)
	}
	if !strings.Contains(output, "Body content") {
		t.Errorf("Render flat: expected 'Body content' in output, got %q", output)
	}
}

func TestCompositeRender_Bad(t *testing.T) {
	// Rendering an empty composite should not panic and return empty string.
	composite := Layout("HCF")
	output := composite.String()
	if output != "" {
		t.Errorf("Empty composite render: expected empty string, got %q", output)
	}
}

func TestCompositeRender_Ugly(t *testing.T) {
	// RenderSimple and RenderBoxed styles add separators between sections.
	UseRenderSimple()
	defer UseRenderFlat()

	composite := Layout("HCF")
	composite.H("top").C("middle").F("bottom")
	output := composite.String()
	if output == "" {
		t.Error("RenderSimple: expected non-empty output")
	}

	UseRenderBoxed()
	output = composite.String()
	if output == "" {
		t.Error("RenderBoxed: expected non-empty output")
	}
}
