package core

import "testing"

func TestSeverityValid(t *testing.T) {
	valid := []Severity{SeverityOK, SeverityInfo, SeverityWarning, SeverityCritical, SeverityUnknown}
	for _, s := range valid {
		if !s.Valid() {
			t.Fatalf("expected %q to be valid", s)
		}
	}
	if Severity("bad").Valid() {
		t.Fatal("unexpected valid severity")
	}
}
