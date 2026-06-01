package conformance

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Result holds the outcome for a single conformance case
type Result struct {
	CaseName string `json:"case_name"`
	// Result is the expected Verify Result value from the case file (PASS/FAIL/ERROR).
	// Keep for backward-compatibility with older consumers.
	Result string `json:"result"`
	// ExpectedResult mirrors Result and makes intent clearer for consumers.
	ExpectedResult string `json:"expected_result"`
	// CasePassed indicates whether the conformance runner considered the
	// case successful (i.e., the expected_verify_result.json was structurally
	// valid and satisfied required fields). This is distinct from the
	// Verify Result value which may be PASS/FAIL/ERROR.
	CasePassed          bool   `json:"case_passed"`
	ReasonCode          string `json:"reason_code"`
	VerifyResultVersion string `json:"verify_result_version"`
	Err                 string `json:"error,omitempty"`
}

// DiscoverCases finds case directories under the provided root directory.
// root is expected to point to testdata/core_2_conformance directory.
func DiscoverCases(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("read root %s: %w", root, err)
	}
	var cases []string
	for _, e := range entries {
		if e.IsDir() {
			// ensure required files exist
			dir := filepath.Join(root, e.Name())
			if _, err := os.Stat(filepath.Join(dir, "input.json")); err != nil {
				continue
			}
			if _, err := os.Stat(filepath.Join(dir, "expected_verify_result.json")); err != nil {
				continue
			}
			cases = append(cases, dir)
		}
	}
	return cases, nil
}

// RunCase performs minimal structural validation of expected_verify_result.json
// and returns a Result. It does not execute verification logic; it validates
// that the expected verify result is structurally coherent for conformance purposes.
func RunCase(caseDir string) Result {
	r := Result{CaseName: filepath.Base(caseDir)}

	b, err := os.ReadFile(filepath.Join(caseDir, "expected_verify_result.json"))
	if err != nil {
		r.Err = fmt.Sprintf("read expected_verify_result: %v", err)
		return r
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		r.Err = fmt.Sprintf("invalid json: %v", err)
		return r
	}

	// minimal checks
	if v, ok := m["result"].(string); ok {
		r.Result = v
	} else {
		r.Err = "missing or invalid 'result' field"
		return r
	}
	if rc, ok := m["reason_code"].(string); ok && rc != "" {
		r.ReasonCode = rc
	} else {
		r.Err = "missing or empty 'reason_code'"
		return r
	}
	if vr, ok := m["verify_result_version"].(string); ok && vr != "" {
		r.VerifyResultVersion = vr
	} else {
		r.Err = "missing or empty 'verify_result_version'"
		return r
	}

	// check allowed result values
	switch r.Result {
	case "PASS", "FAIL", "ERROR":
		// ok
	default:
		r.Err = "invalid 'result' value; must be PASS, FAIL, or ERROR"
		return r
	}

	// mirror expected result and mark case passed when no structural errors
	r.ExpectedResult = r.Result
	r.CasePassed = (r.Err == "")

	return r
}

// RunAll discovers and runs all conformance cases under root.
// root should point to testdata/core_2_conformance; discovery is non-recursive.
func RunAll(root string) ([]Result, error) {
	dirs, err := DiscoverCases(root)
	if err != nil {
		return nil, err
	}
	var out []Result
	for _, d := range dirs {
		res := RunCase(d)
		out = append(out, res)
	}
	return out, nil
}

// ValidateCaseFiles returns an error if required files are missing in the case directory.
func ValidateCaseFiles(caseDir string) error {
	required := []string{"input.json", "expected_verify_result.json"}
	for _, f := range required {
		p := filepath.Join(caseDir, f)
		if _, err := os.Stat(p); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("missing required file: %s", f)
			}
			return fmt.Errorf("stat %s: %w", p, err)
		}
	}
	return nil
}

// Helper to walk directories - returns directories with required files.
func DiscoverCasesByWalk(root string) ([]string, error) {
	var dirs []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		// skip root
		if path == root {
			return nil
		}
		// ensure immediate child dirs only
		rel, _ := filepath.Rel(root, path)
		if rel == "." || filepath.Dir(rel) != "." {
			return nil
		}
		if err := ValidateCaseFiles(path); err == nil {
			dirs = append(dirs, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dirs, nil
}
