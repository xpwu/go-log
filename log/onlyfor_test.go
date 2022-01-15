package log

import (
  "github.com/stretchr/testify/assert"
  "strings"
  "testing"
)

type writerForTest struct {
  expected string
  t        *testing.T
}

func (w *writerForTest) Write(p []byte) (n int, err error) {
  //w.t.Log(string(p))

  strs := strings.SplitN(string(p), "] ", 2)
  a := assert.New(w.t)

  if a.Equalf(2, len(strs), "not find '] ' sep") {
    a.Equal(w.expected, strs[1])
  }

  return len(p), nil
}

type writerForTestAll writerForTest

func (w *writerForTestAll) Write(p []byte) (n int, err error) {
  //w.t.Log(string(p))

  //strs := strings.SplitN(string(p), "] ", 2)
  //a := assert.New(w.t)
  //
  //if a.Equalf(2, len(strs), "not find '] ' sep") {
  //  a.Equal(w.expected, strs[1])
  //}

  return len(p), nil
}
