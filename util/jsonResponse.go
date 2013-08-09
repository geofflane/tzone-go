package util

import "encoding/json"

type JsonResponse map[string]interface{}

func (r JsonResponse) String() (s string) {
  b, err := json.Marshal(r)
  if err != nil {
    s = ""
    return
  }
  s = string(b)
  return
}
