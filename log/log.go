package log

import (
  "fmt"
  config2 "github.com/xpwu/go-config/config"
  "github.com/xpwu/go-log/log/level"
  "io"
  sysLog "log"
  "runtime"
  "strings"
)

//var (
//  level_ = level.DEBUG
//)

type config struct {
  Level level.Level
  Tips  string
}

var configValue = &config{
  Level: level.DEBUG,
  Tips:  level.AllLevelTips(),
}

func init() {
  sysLog.SetFlags(sysLog.Ldate | sysLog.Lmicroseconds)
  config2.Unmarshal(configValue)
}

func SetLevel(l level.Level) {
  configValue.Level = l
}

func SetWriter(writer io.Writer) {
  sysLog.SetOutput(writer)
}

func Writer() io.Writer {
  return sysLog.Writer()
}

func log(l level.Level, skipCallerDepth int, msg interface{}, message ...interface{}) {
  if configValue.Level > l {
    return
  }

  const fileLen = 30
  _, file, line, ok := runtime.Caller(skipCallerDepth + 1)
  if len(file) > fileLen {
    file = "..." + file[len(file)-fileLen:]
  }
  if !ok {
    file = "???"
    line = 0
  }

  builder := strings.Builder{}
  builder.WriteString(file)
  builder.WriteByte(':')
  i := make([]byte, 0)
  itoa(&i, line, -1)
  builder.Write(i)
  builder.WriteString(fmt.Sprintf(" [%v] ", l))
  builder.WriteString(fmt.Sprint(msg))
  builder.WriteString(fmt.Sprint(message...))
  _ = sysLog.Output(0, builder.String())
}

// copy from log/log.go:76
// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
  // Assemble decimal in reverse order.
  var b [20]byte
  bp := len(b) - 1
  for i >= 10 || wid > 1 {
    wid--
    q := i / 10
    b[bp] = byte('0' + i - q*10)
    bp--
    i = q
  }
  // i < 10
  b[bp] = byte('0' + i)
  *buf = append(*buf, b[bp:]...)
}

func Debug(messages ...interface{}) {
  log(level.DEBUG, 1, "", messages...)
}

func Info(messages ...interface{}) {
  log(level.INFO, 1, "", messages...)
}

func Warning(messages ...interface{}) {
  log(level.WARNING, 1, "", messages...)
}

func Error(messages ...interface{}) {
  log(level.ERROR, 1, "", messages...)
}

func Fatal(messages ...interface{}) {
  if configValue.Level > level.FATAL {
    return
  }

  if pt := panicTrace(); pt != "" {
    messages = append(messages, "\n"+pt)

  }

  log(level.FATAL, 1, "", messages...)
}
