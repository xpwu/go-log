package level

type Level int

const (
  DEBUG Level = iota
  INFO
  WARNING
  ERROR
  FATAL
)

var lString = [...]string{"DEBUG", "INFO", "WARNING", "ERROR", "FATAL"}

func (l Level) String() string {
  return lString[l]
}

//func AllLevelTips() string {
//  return "0:DEBUG; 1:INFO; 2:WARNING; 3:ERROR; 4:FATAL"
//}

type LStr string

func (l LStr)Level() Level {
  switch l {
  case "DEBUG":
    return DEBUG
  case "INFO":
    return INFO
  case "WARNING":
    return WARNING
  case "ERROR":
    return ERROR
  case "FATAL":
    return FATAL
  default:
    return DEBUG
  }
}
