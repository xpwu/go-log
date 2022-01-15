package log

import (
  "github.com/stretchr/testify/assert"
  "github.com/xpwu/go-log/log/level"
  sysLog "log"
  "strings"
  "testing"
)

type writerForTest struct {
  result string
  t      *testing.T
}

func (w *writerForTest) Write(p []byte) (n int, err error) {
  strs := strings.SplitN(string(p), "] ", 2)
  if assert.Equalf(w.t, 2, len(strs), "not find '] ' sep") {
    assert.Equal(w.t, w.result, strs[1])
  }

  return len(p), nil
}

func func1(l *sysLog.Logger) {
  func2(l)
}

func func2(l *sysLog.Logger) {
  func3(l)
}

func func3(l *sysLog.Logger) {
  l.Print("real log 3 place")
}

func func4(l *sysLog.Logger) {
  l.Print("real log 4 place, ", "real log 4.1 place")
}

func TestNewSysLog(t *testing.T) {
  l := NewSysLog(NewLogger(), level.INFO)
  func1(l)
  func4(l)
}

func TestLogger_PushPrefix(t *testing.T) {
  l := NewLogger()
  w := &writerForTest{
    result: "",
    t:      t,
  }
  SetWriter(w)

  w.result = "info 1\n"
  l.Info("info 1")

  l.PushPrefix("this is push 1. ")
  w.result = "this is push 1.  info 1\n"
  l.Info("info 1")

  l.PushPrefix("this is push 2")
  w.result = "this is push 1.  this is push 2 info 2 info 2.1\n"
  l.Info("info 2 ", "info 2.1")

  l.PopPrefix()
  w.result = "this is push 1.  info 3\n"
  l.Info("info 3")
  l.PopPrefix()

  w.result = "info 4\n"
  l.Info("info 4")
}

func TestLogger_PushPrefixGo(t *testing.T) {
  l := NewLogger()
  w := &writerForTest{
    result: "",
    t:      t,
  }
  SetWriter(w)

  l.PushPrefix("this is push 1. ")
  w.result = "this is push 1.  info 1\n"
  l.Info("info 1")

  l.PushPrefix("this is push 2")
  w.result = "this is push 1.  this is push 2 info 2 info 2.1\n"
  l.Info("info 2 ", "info 2.1")

  ch := make(chan int)

  go func() {
    l := l.Child()
    l.PopPrefix()

    l.PushPrefix("this is push go")
    w.result = "this is push 1.  this is push 2 this is push go info 2.go info 2.1\n"
    l.Info("info 2.go ", "info 2.1")

    l.PopPrefix()
    l.PopPrefix()
    l.PopPrefix()

    w.result = "this is push 1.  this is push 2 info 2.go info 2.1\n"
    l.Info("info 2.go ", "info 2.1")

    l.PushPrefix("this is push go")
    w.result = "this is push 1.  this is push 2 this is push go info 2.go1 info 2.1\n"
    l.Info("info 2.go1 ", "info 2.1")

    close(ch)
  }()

  <-ch

  l.PopPrefix()
  w.result = "this is push 1.  info 3\n"
  l.Info("info 3")
  l.PopPrefix()

  w.result = "info 4\n"
  l.Info("info 4")
}
