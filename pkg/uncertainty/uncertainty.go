package uncertainty

// Level describes uncertainty on a coarse public scale.
// v1.0 intentionally keeps this minimal and stable.
type Level string

const (
	Low    Level = "low"
	Medium Level = "medium"
	High   Level = "high"
)

// Annotation attaches uncertainty information to an object reference.
type Annotation struct {
	Level Level    `json:"level"`
	Note  string   `json:"note,omitempty"`
	Refs  []string `json:"refs,omitempty"`
}

func (a Annotation) Validate() error {
	switch a.Level {
	case Low, Medium, High:
		return nil
	default:
		return ErrInvalidLevel
	}
}
