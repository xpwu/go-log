# go-log

并发安全的Log格式化输出，最终的输出使用系统的log

### usage
使用 WithCtx() 生成新的logger  

```
func F1() {
  ctx := context.TODO()
  ctx, logger := WithCtx(ctx)
  logger.PushPrefix("F1")
  logger.Info("first") // output: <file:line [INFO] >F1 first
  
  F2(ctx)
  
  // F2() 添加的Prefix并不会影响到本层函数
  logger.Info("fifth") // output: <file:line [INFO] >F1 fifth
}

func F2(ctx context.Context) {
  ctx, logger := WithCtx(ctx)
  logger.PushPrefix("F2")
  logger.Info("second") // output: <file:line [INFO] >F1 F2 second
  
  logger.PushPrefix("F3")
  logger.Info("third")   // output: <file:line [INFO] >F1 F2 F3 third
  
  logger.PopPrefix()
  logger.Info("fourth")   // output: <file:line [INFO] >F1 F2 fourth
}

```
