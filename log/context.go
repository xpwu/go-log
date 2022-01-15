package log

import (
  "context"
)

type logKey struct {}

func WithCtx(parent context.Context) (ctx context.Context, logger *Logger) {
  switch value := parent.Value(logKey{}).(type) {
  case *Logger:
    logger = value.Child()
  default:
    logger = NewLogger()
  }

  ctx = context.WithValue(parent, logKey{},logger)
  return
}

func NewContext(parent context.Context, logger *Logger) context.Context {
  if value,ok := parent.Value(logKey{}).(*Logger); ok && value == logger {
    return parent
  }
  return context.WithValue(parent, logKey{}, logger)
}

