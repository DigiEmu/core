package replay

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Result is the deterministic replay output for v0.5.
// It intentionally stays minimal and stable.
type Result struct {
	Canonical []byte
	Size      int
	Source    string
	SHA256    string
}

type jsonResult struct {
	Canonical string `json:"canonical"`
	Size      int    `json:"size"`
	Source    string `json:"source"`
	SHA256    string `json:"sha256"`
}

func newResult(data []byte, source string) (Result, error) {
	if len(data) == 0 {
		return Result{}, fmt.Errorf("replay: empty input")
	}
	if source == "" {
		return Result{}, fmt.Errorf("replay: empty source")
	}

	out := make([]byte, len(data))
	copy(out, data)

	sum := sha256.Sum256(out)

	return Result{
		Canonical: out,
		Size:      len(out),
		Source:    source,
		SHA256:    hex.EncodeToString(sum[:]),
	}, nil
}

// Marshal returns a deterministic JSON representation of the replay result.
func (r Result) Marshal() ([]byte, error) {
	if len(r.Canonical) == 0 {
		return nil, fmt.Errorf("replay: empty canonical")
	}
	if r.Size <= 0 {
		return nil, fmt.Errorf("replay: invalid size")
	}
	if r.Source == "" {
		return nil, fmt.Errorf("replay: empty source")
	}
	if r.SHA256 == "" {
		return nil, fmt.Errorf("replay: empty sha256")
	}

	return json.Marshal(jsonResult{
		Canonical: string(r.Canonical),
		Size:      r.Size,
		Source:    r.Source,
		SHA256:    r.SHA256,
	})
}

// FromBytes returns a deterministic replay result from canonical bundle bytes.
func FromBytes(data []byte) (Result, error) {
	return newResult(data, "bytes")
}

// FromReader reads all bytes from r and returns a deterministic replay result.
func FromReader(r io.Reader) (Result, error) {
	if r == nil {
		return Result{}, fmt.Errorf("replay: nil reader")
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return Result{}, fmt.Errorf("replay: read reader: %w", err)
	}

	return newResult(data, "reader")
}

// FromFile reads a file fully and returns a deterministic replay result.
func FromFile(path string) (Result, error) {
	if path == "" {
		return Result{}, fmt.Errorf("replay: empty path")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Result{}, fmt.Errorf("replay: read file: %w", err)
	}

	return newResult(data, "file")
}
