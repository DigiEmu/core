package usecases

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"digiemu-core/internal/kernel/domain"
	"digiemu-core/internal/kernel/ports"
)

// snapshotCanonicalLines builds a stable, line-based representation.
// It avoids relying on json.Marshal map iteration order.
// Format (one record per line):
//
//	UNIT|<unitID>|<key>|<title>|<desc>|<headVersionID>
//	VER|<id>|<prev>|<contentHash>
//	AUD|<id>|<type>|<atUnix>|<actor>|<unit>|<ver>|<dataCanonical>
func snapshotCanonicalLines(u ports.UnitDTO, vs []ports.VersionDTO) []string {
	lines := make([]string, 0, 1+len(vs))
	lines = append(lines, fmt.Sprintf(
		"UNIT|%s|%s|%s|%s|%s",
		u.ID, u.Key, u.Title, u.Description, u.HeadVersionID,
	))
	for _, v := range vs {
		lines = append(lines, fmt.Sprintf(
			"VER|%s|%s|%s",
			v.ID, v.PrevVersionID, v.ContentHash,
		))
	}
	return lines
}

func auditCanonicalLines(evs []domain.AuditEvent) ([]string, error) {
	lines := make([]string, 0, len(evs))
	for _, ev := range evs {
		dataCanon, err := canonicalJSON(ev.Data)
		if err != nil {
			return nil, err
		}
		lines = append(lines, fmt.Sprintf(
			"AUD|%s|%s|%d|%s|%s|%s|%s",
			ev.ID, ev.Type, ev.AtUnix, ev.ActorID, ev.UnitID, ev.VersionID, dataCanon,
		))
	}
	return lines, nil
}

func sha256HexFromLines(lines []string) string {
	joined := strings.Join(lines, "\n") + "\n"
	sum := sha256.Sum256([]byte(joined))
	return hex.EncodeToString(sum[:])
}

// canonicalJSON renders arbitrary JSON-like structures with stable map key ordering.
// Supports: nil, bool, float64/int-like, string, []any, map[string]any, map[string]string, json.RawMessage.
func canonicalJSON(v any) (string, error) {
	switch t := v.(type) {
	case nil:
		return "null", nil
	case string:
		b, _ := json.Marshal(t)
		return string(b), nil
	case bool:
		if t {
			return "true", nil
		}
		return "false", nil
	case json.RawMessage:
		var x any
		if err := json.Unmarshal(t, &x); err != nil {
			b, _ := json.Marshal(string(t))
			return string(b), nil
		}
		return canonicalJSON(x)
	case []any:
		parts := make([]string, 0, len(t))
		for _, it := range t {
			s, err := canonicalJSON(it)
			if err != nil {
				return "", err
			}
			parts = append(parts, s)
		}
		return "[" + strings.Join(parts, ",") + "]", nil
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			ks, _ := json.Marshal(k)
			vs, err := canonicalJSON(t[k])
			if err != nil {
				return "", err
			}
			parts = append(parts, string(ks)+":"+vs)
		}
		return "{" + strings.Join(parts, ",") + "}", nil
	case map[string]string:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			ks, _ := json.Marshal(k)
			vs, _ := json.Marshal(t[k])
			parts = append(parts, string(ks)+":"+string(vs))
		}
		return "{" + strings.Join(parts, ",") + "}", nil
	default:
		b, err := json.Marshal(t)
		if err != nil {
			return "", err
		}
		var x any
		if err := json.Unmarshal(b, &x); err != nil {
			return "", err
		}
		return canonicalJSON(x)
	}
}
