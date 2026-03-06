package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type identityDoc struct {
	Name      string `json:"name"`
	Algorithm string `json:"algorithm"`
}

func ensureLocalIdentity() (identityDoc, ed25519.PublicKey, ed25519.PrivateKey, error) {
	baseDir := ".digiemu"
	keyPath := filepath.Join(baseDir, "identity.key")
	pubPath := filepath.Join(baseDir, "identity.pub")
	metaPath := filepath.Join(baseDir, "identity.json")

	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return identityDoc{}, nil, nil, err
	}

	keyInfo, keyErr := os.Stat(keyPath)
	pubInfo, pubErr := os.Stat(pubPath)
	metaInfo, metaErr := os.Stat(metaPath)

	if keyErr == nil && pubErr == nil && metaErr == nil && !keyInfo.IsDir() && !pubInfo.IsDir() && !metaInfo.IsDir() {
		keyHex, err := os.ReadFile(keyPath)
		if err != nil {
			return identityDoc{}, nil, nil, err
		}
		pubHex, err := os.ReadFile(pubPath)
		if err != nil {
			return identityDoc{}, nil, nil, err
		}
		metaBytes, err := os.ReadFile(metaPath)
		if err != nil {
			return identityDoc{}, nil, nil, err
		}

		var meta identityDoc
		if err := json.Unmarshal(metaBytes, &meta); err != nil {
			return identityDoc{}, nil, nil, err
		}
		if meta.Algorithm != "ed25519" {
			return identityDoc{}, nil, nil, fmt.Errorf("unsupported identity algorithm: %q", meta.Algorithm)
		}

		privBytes, err := hex.DecodeString(strings.TrimSpace(string(keyHex)))
		if err != nil {
			return identityDoc{}, nil, nil, err
		}
		pubBytes, err := hex.DecodeString(strings.TrimSpace(string(pubHex)))
		if err != nil {
			return identityDoc{}, nil, nil, err
		}
		if len(privBytes) != ed25519.PrivateKeySize {
			return identityDoc{}, nil, nil, fmt.Errorf("invalid identity private key size: %d", len(privBytes))
		}
		if len(pubBytes) != ed25519.PublicKeySize {
			return identityDoc{}, nil, nil, fmt.Errorf("invalid identity public key size: %d", len(pubBytes))
		}

		return meta, ed25519.PublicKey(pubBytes), ed25519.PrivateKey(privBytes), nil
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return identityDoc{}, nil, nil, err
	}

	meta := identityDoc{
		Name:      "local",
		Algorithm: "ed25519",
	}

	metaBytes, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return identityDoc{}, nil, nil, err
	}

	if err := os.WriteFile(keyPath, []byte(hex.EncodeToString(priv)), 0o600); err != nil {
		return identityDoc{}, nil, nil, err
	}
	if err := os.WriteFile(pubPath, []byte(hex.EncodeToString(pub)), 0o644); err != nil {
		return identityDoc{}, nil, nil, err
	}
	if err := os.WriteFile(metaPath, metaBytes, 0o644); err != nil {
		return identityDoc{}, nil, nil, err
	}

	return meta, pub, priv, nil
}

func loadLocalIdentityPublic() (identityDoc, ed25519.PublicKey, error) {
	metaPath := filepath.Join(".digiemu", "identity.json")
	pubPath := filepath.Join(".digiemu", "identity.pub")

	metaBytes, err := os.ReadFile(metaPath)
	if err != nil {
		return identityDoc{}, nil, fmt.Errorf("load identity metadata: %w", err)
	}
	pubHex, err := os.ReadFile(pubPath)
	if err != nil {
		return identityDoc{}, nil, fmt.Errorf("load identity public key: %w", err)
	}

	var meta identityDoc
	if err := json.Unmarshal(metaBytes, &meta); err != nil {
		return identityDoc{}, nil, err
	}
	if meta.Algorithm != "ed25519" {
		return identityDoc{}, nil, fmt.Errorf("unsupported identity algorithm: %q", meta.Algorithm)
	}

	pubBytes, err := hex.DecodeString(strings.TrimSpace(string(pubHex)))
	if err != nil {
		return identityDoc{}, nil, err
	}
	if len(pubBytes) != ed25519.PublicKeySize {
		return identityDoc{}, nil, fmt.Errorf("invalid identity public key size: %d", len(pubBytes))
	}

	return meta, ed25519.PublicKey(pubBytes), nil
}
