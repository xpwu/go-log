package log

import (
  "github.com/xpwu/go-log/log/level"
  sysLog "log"
  "testing"
)


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
    expected: "",
    t:        t,
  }

  oldW := Writer()
  defer SetWriter(oldW)

  SetWriter(w)

  w.expected = "info 1\n"
  l.Info("info 1")

  l.PushPrefix("this is push 1. ")
  w.expected = "this is push 1.  info 1\n"
  l.Info("info 1")

  l.PushPrefix("this is push 2")
  w.expected = "this is push 1.  this is push 2 info 2 info 2.1\n"
  l.Info("info 2 ", "info 2.1")

  l.PopPrefix()
  w.expected = "this is push 1.  info 3\n"
  l.Info("info 3")
  l.PopPrefix()

  w.expected = "info 4\n"
  l.Info("info 4")
}

func TestLogger_PushPrefixGo(t *testing.T) {
  l := NewLogger()
  w := &writerForTest{
    expected: "",
    t:        t,
  }

  oldW := Writer()
  defer SetWriter(oldW)

  SetWriter(w)

  l.PushPrefix("this is push 1. ")
  w.expected = "this is push 1.  info 1\n"
  l.Info("info 1")

  l.PushPrefix("this is push 2")
  w.expected = "this is push 1.  this is push 2 info 2 info 2.1\n"
  l.Info("info 2 ", "info 2.1")

  ch := make(chan int)

  go func() {
    l := l.Child()
    l.PopPrefix()

    l.PushPrefix("this is push go")
    w.expected = "this is push 1.  this is push 2 this is push go info 2.go info 2.1\n"
    l.Info("info 2.go ", "info 2.1")

    l.PopPrefix()
    l.PopPrefix()
    l.PopPrefix()

    w.expected = "this is push 1.  this is push 2 info 2.go info 2.1\n"
    l.Info("info 2.go ", "info 2.1")

    l.PushPrefix("this is push go")
    w.expected = "this is push 1.  this is push 2 this is push go info 2.go1 info 2.1\n"
    l.Info("info 2.go1 ", "info 2.1")

    close(ch)
  }()

  <-ch

  l.PopPrefix()
  w.expected = "this is push 1.  info 3\n"
  l.Info("info 3")
  l.PopPrefix()

  w.expected = "info 4\n"
  l.Info("info 4")
}
