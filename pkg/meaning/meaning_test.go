package meaning

import "testing"

// This test is intentionally minimal: it ensures the package compiles and exported
// surface stays stable without importing internals.
func TestPkgMeaning_Basics(t *testing.T) {
	// If you add exported constructors/validators later, extend tests here.
	t.Log("pkg/meaning compile smoke OK")
}
