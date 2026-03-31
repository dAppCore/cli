package cli

import "testing"

func TestParseVariant_Good(t *testing.T) {
	composite, err := ParseVariant("H[LC]F")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if _, ok := composite.regions[RegionHeader]; !ok {
		t.Error("Expected Header region")
	}
	if _, ok := composite.regions[RegionFooter]; !ok {
		t.Error("Expected Footer region")
	}

	headerSlot := composite.regions[RegionHeader]
	if headerSlot.child == nil {
		t.Error("Header should have child layout for H[LC]")
	} else {
		if _, ok := headerSlot.child.regions[RegionLeft]; !ok {
			t.Error("Child should have Left region")
		}
	}
}

func TestParseVariant_Bad(t *testing.T) {
	// Invalid region character.
	_, err := ParseVariant("X")
	if err == nil {
		t.Error("Expected error for invalid region character 'X'")
	}

	// Unmatched bracket.
	_, err = ParseVariant("H[C")
	if err == nil {
		t.Error("Expected error for unmatched bracket")
	}
}

func TestParseVariant_Ugly(t *testing.T) {
	// Empty variant should produce empty composite without panic.
	composite, err := ParseVariant("")
	if err != nil {
		t.Fatalf("Empty variant should not error: %v", err)
	}
	if len(composite.regions) != 0 {
		t.Errorf("Empty variant should have no regions, got %d", len(composite.regions))
	}
}
