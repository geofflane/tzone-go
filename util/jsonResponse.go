package util

import "encoding/json"

type JsonResponse map[string]interface{}

// Implements the fmt.Stringer interface
// Converts a JsonResponse to a string representation
func (r JsonResponse) String() (s string) {
  b, err := json.Marshal(r)
  if err != nil {
    s = ""
    return
  }
  s = string(b)
  return
}
