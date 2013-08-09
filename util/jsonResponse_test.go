package util

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestMarshalStringToJson(t *testing.T) {
  resp := JsonResponse{"test":"val"}
  json := resp.String()
  assert.Equal(t, "{\"test\":\"val\"}", json)
}

func TestMarshalIntToJson(t *testing.T) {
  resp := JsonResponse{"test": 1}
  json := resp.String()
  assert.Equal(t, "{\"test\":1}", json)
}

