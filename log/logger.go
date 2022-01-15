package log

import (
  "github.com/xpwu/go-log/log/level"
  sysLog "log"
  "strings"
  "unsafe"
)

/*
 Logger 中存储pre 设计过两种方式
  1、使用一个slice存储所有的pre节点。
  优点：可以使用底层的slice方法
  缺点：1）、每一个Logger必须有独立的slice，不能两个Logger共用一个，所以每次child生成的Logger在第一次Push操作时必须全部copy一次。
      如果没有独立使用，在两个goroutine同时push时，或者一个goroutine pop后再push会出现数据错误；(如果只是pop操作没有影响)
      2）、Logger不能直接值copy, 因为copy后就是共用的同一个slice，必须加地址判断方法
  2、使用一个反向的链表存储所有的节点。
  缺点：需要实现slice的一些方法(就是拼接所有的前缀)，因为是反向，所以每个节点必须记录在整个[]byte中的位置，同时每一个节点的值一旦写入
    就不可再次修改，如果修改，可能出现显示错误
  优点：不用复制；可以值copy

  目前选用的方案2

并发：
同一个Logger不支持并发，所有在并发时，应该使用Child新生成一个后再push或者pop。为了安全，在不确定调用方是否是并发时，使用child新
生成一个总是安全的。如果使用WithCtx()，则在每次调用此函数时都是新生成一个

*Logger,值copy与child的区别：
值copy同时拷贝了eve的指向，能够pop到原对象的第一个push点，push时只与自己有关；
child 重新指定eve，pop时只能pop到当前点，parent push的数据无法pop，push时只与自己有关；
*Logger 是完全相同的两个Logger, push与pop完全一样，对一个做了操作，另一个同时受影响。

常用*Logger与child方式，在明确是同一个goroutine时，应该使用*Logger，在不同的goroutine或者不确定是否是同一个goroutine时，应该使用child


pop:
如果是在一个函数中使用WithCtx/Child获取的Logger，在此函数中Push的数据，退出此函数时，因为作用域的原因
自动释放掉，也不会影响到parent Logger.
如果*Logger传递的方式，则根据逻辑的具体需要，而决定是否调用pop，在函数中使用时，可以借用defer。

*/

type node struct {
  value   string
  bytePos int
  pre     *node
}

func (n *node) append(value string) *node {
  if n == nil {
    return &node{
      value: value,
    }
  }

  return &node{
    value:   value,
    bytePos: n.bytePos + len(n.value),
    pre:     n,
  }
}

func (n *node) String() string {
  if n == nil {
    return ""
  }

  r := make([]byte, n.bytePos+len(n.value))
  l := n
  // 防止value的值写入后再次被修改，出现位置错误
  lastPos := l.bytePos + len(n.value)
  for l != nil {
    t := r[l.bytePos:lastPos]
    copy(t, l.value)
    lastPos = l.bytePos
    l = l.pre
  }

  return *(*string)(unsafe.Pointer(&r))
}

func (n *node) Len() int {
  if n == nil {
    return 0
  }
  return n.bytePos+len(n.value)
}

func (n *node) WriteToBuilder(b *strings.Builder) {
  b.Grow(n.Len())
  if n == nil {
    return
  }
  n.pre.WriteToBuilder(b)
  b.WriteString(n.value)
}

type Logger struct {
  prefix          *node
  eve             *node
  skipCallerDepth int
}

func (l *Logger) Info(messages ...interface{}) {
  log(level.INFO, 1+l.skipCallerDepth, l.prefix, messages...)
}

func (l *Logger) Debug(messages ...interface{}) {
  log(level.DEBUG, 1+l.skipCallerDepth, l.prefix, messages...)
}

func (l *Logger) Warning(messages ...interface{}) {
  log(level.WARNING, 1+l.skipCallerDepth, l.prefix, messages...)
}

func (l *Logger) Error(messages ...interface{}) {
  log(level.ERROR, 1+l.skipCallerDepth, l.prefix, messages...)
}

func (l *Logger) Fatal(messages ...interface{}) {
  if configValue.Level > level.FATAL {
    return
  }

  if pt := panicTrace(); pt != "" {
    messages = append(messages, "\n"+pt)
  }

  log(level.FATAL, 1+l.skipCallerDepth, l.prefix, messages...)
}

// 不是并发安全的，如果要在另一个go程中使用同一个logger的prefix  可以使用Child() 新生成一个
func (l *Logger) PushPrefix(prefix string) {
  l.prefix = l.prefix.append(prefix + " ")
}

// 只能pop由自己Push的prefix
func (l *Logger) PopPrefix() {
  if l.prefix == l.eve {
    return
  }

  l.prefix = l.prefix.pre
}

func (l *Logger) Prefix() *node {
  return l.prefix
}

func (l *Logger) Child() *Logger {
  return &Logger{
    prefix:          l.prefix,
    eve:             l.prefix,
    skipCallerDepth: l.skipCallerDepth,
  }
}

func (l *Logger) AddSkipCallerDepth(depth int) *Logger {
  return &Logger{
    prefix:          l.prefix,
    eve:             l.prefix,
    skipCallerDepth: l.skipCallerDepth + depth,
  }
}

func NewLogger() *Logger {
  return &Logger{}
}

func NewSysLog(logger *Logger, l level.Level) *sysLog.Logger {
  return sysLog.New(&logWrapper{logger: logger, level: l}, "", 0)
}

type logWrapper struct {
  logger *Logger
  level  level.Level
}

func (l *logWrapper) Write(p []byte) (n int, err error) {
  // todo: fatal

  log(l.level, 3+l.logger.skipCallerDepth, l.logger.prefix, string(p))
  return len(p), nil
}

/*
Do not copy a non-zero Logger.
type Logger struct {
 addr   *Logger
 prefix []interface{}
 first  int // 记录属于此Logger push的第一个prefix的位置
 alone  bool
}

func (l *Logger) checkAlone() {
 // copy and modify from strings/builder.go
 if l.addr == nil {
   l.addr = l
 } else if l.addr != l {
   panic("log: illegal use of non-zero Logger copied by value")
 }

 if l.alone {
  return
 }

 l.alone = true
 n := make([]interface{}, len(l.prefix), cap(l.prefix))
 copy(n, l.prefix)
 l.prefix = n
}

// 不是并发安全的，如果要在另一个go程中使用同一个logger的prefix  可以使用Child() 新生成一个
func (l *Logger) PushPrefix(prefix string) {
  l.checkAlone()
  l.prefix = append(l.prefix, prefix, " ")
}

// 只能pop由自己Push的prefix
func (l *Logger) PopPrefix() {
  if !l.alone || len(l.prefix) <= l.first {
   return
  }

  // push 一共加入了两个元素
  l.prefix = l.prefix[:len(l.prefix)-2]
}

func (l *Logger) Child() *Logger {
  return &Logger{
   prefix: l.prefix,
   alone:  false,
   first:  len(l.prefix),
  }
}

// 其他接口注意l.prefix 传给 log方法时的处理

*/
