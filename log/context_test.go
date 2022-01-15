package log

import (
  "context"
  "testing"
)

func F2(ctx context.Context, w *writerForTest) {
  ctx, logger := WithCtx(ctx)
  logger.PushPrefix("F2")
  w.expected = "F1 F2 second\n"
  logger.Info("second") // output: <file:line [INFO] >F1 F2 second

  logger.PushPrefix("F3")
  w.expected = "F1 F2 F3 third\n"
  logger.Info("third")   // output: <file:line [INFO] >F1 F2 F3 third

  logger.PopPrefix()
  w.expected = "F1 F2 fourth\n"
  logger.Info("fourth")   // output: <file:line [INFO] >F1 F2 fourth
}

func TestWithCtx(t *testing.T) {
  w := &writerForTest{
    expected: "",
    t:        t,
  }
  oldW := Writer()
  defer SetWriter(oldW)

  SetWriter(w)

  ctx := context.TODO()
  ctx, logger := WithCtx(ctx)
  logger.PushPrefix("F1")
  w.expected = "F1 first\n"
  logger.Info("first") // output: <file:line [INFO] >F1 first

  F2(ctx, w)

  // F2() 添加的Prefix并不会影响到本层函数
  w.expected = "F1 fifth\n"
  logger.Info("fifth") // output: <file:line [INFO] >F1 fifth
}
