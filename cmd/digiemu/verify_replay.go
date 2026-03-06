package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	pkgreplay "digiemu-core/pkg/replay"
)

type verifyReplayBundleDoc struct {
	Version  string `json:"version"`
	Ref      string `json:"ref"`
	Snapshot string `json:"snapshot"`
	Meta     string `json:"meta"`
	Trace    string `json:"trace"`
}

type verifyReplayTraceDoc struct {
	SHA256       string `json:"sha256"`
	ReplaySHA256 string `json:"replay_sha256"`
	Mode         string `json:"mode"`
}

func runVerifyReplayWithIO(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		fmt.Fprintln(stderr, "usage: digiemu verify replay <bundle.json>")
		return 2
	}

	bundlePath := args[0]
	bundleBytes, err := os.ReadFile(bundlePath)
	if err != nil {
		fmt.Fprintf(stderr, "verify replay: %v\n", err)
		return 1
	}

	var bundle verifyReplayBundleDoc
	if err := json.Unmarshal(bundleBytes, &bundle); err != nil {
		fmt.Fprintln(stderr, "verify replay: invalid json")
		return 1
	}
	if bundle.Version == "" || bundle.Ref == "" || bundle.Snapshot == "" || bundle.Trace == "" {
		fmt.Fprintln(stderr, "verify replay: invalid schema")
		return 1
	}

	baseDir := filepath.Dir(bundlePath)
	snapshotPath := filepath.Join(baseDir, bundle.Snapshot)
	tracePath := filepath.Join(baseDir, bundle.Trace)

	snapshotBytes, err := os.ReadFile(snapshotPath)
	if err != nil {
		fmt.Fprintf(stderr, "verify replay: read snapshot: %v\n", err)
		return 1
	}

	traceBytes, err := os.ReadFile(tracePath)
	if err != nil {
		fmt.Fprintf(stderr, "verify replay: read trace: %v\n", err)
		return 1
	}
	var trace verifyReplayTraceDoc
	if err := json.Unmarshal(traceBytes, &trace); err != nil {
		fmt.Fprintln(stderr, "verify replay: invalid trace json")
		return 1
	}

	r1, err := pkgreplay.FromBytes(snapshotBytes)
	if err != nil {
		fmt.Fprintf(stderr, "verify replay: %v\n", err)
		return 1
	}
	out1, err := r1.Marshal()
	if err != nil {
		fmt.Fprintf(stderr, "verify replay: marshal result: %v\n", err)
		return 1
	}

	r2, err := pkgreplay.FromBytes(snapshotBytes)
	if err != nil {
		fmt.Fprintf(stderr, "verify replay: %v\n", err)
		return 1
	}
	out2, err := r2.Marshal()
	if err != nil {
		fmt.Fprintf(stderr, "verify replay: marshal result: %v\n", err)
		return 1
	}

	if !bytes.Equal(out1, out2) {
		fmt.Fprintln(stderr, "verify replay: output not deterministic")
		return 1
	}

	if trace.ReplaySHA256 != "" {
		sum := sha256.Sum256(out1)
		actual := hex.EncodeToString(sum[:])
		if actual != trace.ReplaySHA256 {
			fmt.Fprintln(stderr, "replay mismatch")
			return 1
		}
	}

	_ = trace.Mode
	_ = trace.SHA256

	fmt.Fprintln(stdout, "OK replay deterministic")
	return 0
}
