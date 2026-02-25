package canonicaljson

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

// Marshal encodes v to canonical JSON bytes according to the project's rules:
// - object keys (maps) sorted lexicographically
// - structs serialized using json tag names when present, in declaration order
// - arrays/slices preserved in order
// - no insignificant whitespace
// This implementation intentionally avoids third-party packages.
func Marshal(v any) ([]byte, error) {
	buf := &bytes.Buffer{}
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return []byte("null"), nil
	}
	if err := encodeValue(rv, buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encodeValue(v reflect.Value, buf *bytes.Buffer) error {
	if !v.IsValid() {
		buf.WriteString("null")
		return nil
	}
	// unwrap interfaces and pointers
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer {
		if v.IsNil() {
			buf.WriteString("null")
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
		return nil
	case reflect.String:
		buf.WriteString(strconv.Quote(v.String()))
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.WriteString(strconv.FormatInt(v.Int(), 10))
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		buf.WriteString(strconv.FormatUint(v.Uint(), 10))
		return nil
	case reflect.Float32, reflect.Float64:
		// Use strconv with a stable format
		buf.WriteString(strconv.FormatFloat(v.Float(), 'g', -1, 64))
		return nil
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return errors.New("canonicaljson: map keys must be strings")
		}
		keys := make([]string, 0, v.Len())
		for _, kv := range v.MapKeys() {
			keys = append(keys, kv.String())
		}
		sort.Strings(keys)
		buf.WriteByte('{')
		for i, k := range keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(strconv.Quote(k))
			buf.WriteByte(':')
			if err := encodeValue(v.MapIndex(reflect.ValueOf(k)), buf); err != nil {
				return err
			}
		}
		buf.WriteByte('}')
		return nil
	case reflect.Slice, reflect.Array:
		buf.WriteByte('[')
		n := v.Len()
		for i := 0; i < n; i++ {
			if i > 0 {
				buf.WriteByte(',')
			}
			if err := encodeValue(v.Index(i), buf); err != nil {
				return err
			}
		}
		buf.WriteByte(']')
		return nil
	case reflect.Struct:
		buf.WriteByte('{')
		t := v.Type()
		first := true
		for i := 0; i < v.NumField(); i++ {
			f := t.Field(i)
			if f.PkgPath != "" { // unexported
				continue
			}
			tag := f.Tag.Get("json")
			name := f.Name
			if tag != "" {
				// tag may include options like `json:"foo,omitempty"`
				if tag == "-" {
					continue
				}
				// extract name before comma
				if comma := indexComma(tag); comma >= 0 {
					name = tag[:comma]
				} else {
					name = tag
				}
				if name == "" {
					name = f.Name
				}
			}
			if !first {
				buf.WriteByte(',')
			}
			first = false
			buf.WriteString(strconv.Quote(name))
			buf.WriteByte(':')
			if err := encodeValue(v.Field(i), buf); err != nil {
				return err
			}
		}
		buf.WriteByte('}')
		return nil
	case reflect.Invalid:
		buf.WriteString("null")
		return nil
	default:
		return fmt.Errorf("canonicaljson: unsupported kind %s", v.Kind())
	}
}

func indexComma(tag string) int {
	for i := 0; i < len(tag); i++ {
		if tag[i] == ',' {
			return i
		}
	}
	return -1
}
