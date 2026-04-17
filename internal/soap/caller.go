package soap

import "context"

// Caller is the interface for making SOAP calls. Handlers depend on this for testability.
type Caller interface {
	Call(ctx context.Context, method string, params []Param) ([]byte, error)
}

// Param represents one SOAP method parameter.
// Value types determine XML encoding:
//   - string, int, int64, bool → simple element with xsi:type
//   - map[string]any           → ns2:Map encoding
//   - []map[string]any         → SOAP-ENC:Array of Maps
//   - []int64 / []string       → SOAP-ENC:Array of primitives
//   - float64                  → xsd:float element
//   - nil                      → xsi:nil="true"
type Param struct {
	Name  string
	Value any
}
