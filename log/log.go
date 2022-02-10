package log

import (
  "fmt"
  "github.com/xpwu/go-config/configs"
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
  Level level.Level `conf:"level,0:DEBUG; 1:INFO; 2:WARNING; 3:ERROR; 4:FATAL"`
}

var configValue = &config{
  Level: level.DEBUG,
}

func init() {
  sysLog.SetFlags(sysLog.Ldate | sysLog.Lmicroseconds)
  configs.Unmarshal(configValue)
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

type BuilderWriter interface {
  WriteToBuilder(b *strings.Builder)
}

type LenGetter interface {
  Len() int
}

type LazyMsg func() string

func (l LazyMsg) String() string {
  return l()
}

func log(l level.Level, skipCallerDepth int, msg interface{}, messages ...interface{}) {
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

  // file:line [level] msg+messages

  lineB := make([]byte, 0)
  itoa(&lineB, line, -1)

  levelStr := fmt.Sprintf(" [%v] ", l)

  msgLen := 0
  switch m := msg.(type) {
  case LenGetter:
    msgLen = m.Len()
  case string:
    msgLen = len(m)
  case []byte:
    msgLen = len(m)
  case *string:
    msgLen = len(*m)
  }

  // todo: 支持BuilderWriter接口
  messagesStr := fmt.Sprint(messages...)

  length := len(file) + 1 + len(lineB) + len(levelStr) + msgLen + len(messagesStr)

  builder := &strings.Builder{}
  builder.Grow(length)

  builder.WriteString(file)
  builder.WriteByte(':')
  builder.Write(lineB)
  builder.WriteString(levelStr)

  if m,ok := msg.(BuilderWriter); ok {
    m.WriteToBuilder(builder)
  } else {
    builder.WriteString(fmt.Sprint(msg))
  }

  builder.WriteString(messagesStr)

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
