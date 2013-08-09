package util

import "testing"

func TestMarshalStringToJson(t *testing.T) {
  resp := JsonResponse{"test":"val"}
  json := resp.String()
  jsonT{t}.assertJsonMatches("{\"test\":\"val\"}", json)
}

func TestMarshalIntToJson(t *testing.T) {
  resp := JsonResponse{"test": 1}
  json := resp.String()
  jsonT{t}.assertJsonMatches("{\"test\":1}", json)
}

type jsonT struct {
  *testing.T
}

func (t jsonT) assertJsonMatches(expected string, actual string) {
  if expected != actual {
    t.Errorf("Not expected json: %s", actual)
  }
}
