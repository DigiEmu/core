package verify

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadBundleFromPath_StripsBOM(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "snapshot.json")

	// JSON content with UTF-8 BOM prefix
	content := []byte{0xEF, 0xBB, 0xBF}
	content = append(content, []byte(`{"ref":"demo","expected_hash_v1":"abc","state": {"x":1}}`)...)

	if err := os.WriteFile(p, content, 0644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	sb, err := loadBundleFromPath(p, "demo")
	if err != nil {
		t.Fatalf("loadBundleFromPath failed: %v", err)
	}
	if sb.Ref != "demo" {
		t.Fatalf("expected ref demo, got %s", sb.Ref)
	}
	if sb.ExpectedHashV1 != "abc" {
		t.Fatalf("expected expected_hash_v1 abc, got %s", sb.ExpectedHashV1)
	}
}
