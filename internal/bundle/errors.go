package bundle

import "errors"

var ErrExpectedAlreadySet = errors.New("expected_hash_v1 already set")
var ErrSnapshotNotFound = errors.New("snapshot not found")
var ErrSnapshotInvalidJSON = errors.New("snapshot invalid json")
var ErrInvalidNewHash = errors.New("invalid new hash")
