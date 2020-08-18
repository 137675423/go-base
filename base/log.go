package base

import (
	"fmt"
	"time"
)

type LogLevel int

const (
	Waring LogLevel = iota
	Info
	Err
)

var ShowLevel = map[LogLevel]string{
	Waring: "Waring",
	Info:   "Info",
	Err:    "Err",
}

//事务步进
type step struct {
	level   LogLevel
	content interface{}
}

func newStep(l LogLevel, content ...interface{}) step {
	return step{l, content}
}

type Logger struct {
	//开始时间
	StartTime time.Time
	//途经路线集合
	Steps []step
}

func NewLogger() *Logger {
	return &Logger{
		time.Now(), nil,
	}
}

func (l *Logger) SaveFile() {
	str := fmt.Sprintf("LOG BEGIN AT [%v] Cost Time %v \n", l.StartTime.Format("2006-01-02 15:04:05"), time.Now().Sub(l.StartTime))
	for k, v := range l.Steps {
		str += fmt.Sprintf("Step:%d | Level:%v | Events:%v\n", k+1, ShowLevel[v.level], v.content)
	}
	str += "LOG END"
	fmt.Println(str)
}

func (l *Logger) Waring(content ...interface{}) {
	l.Steps = append(l.Steps, newStep(Waring, content...))
}
func (l *Logger) Info(content ...interface{}) {
	l.Steps = append(l.Steps, newStep(Info, content...))
}

func (l *Logger) Err(content ...interface{}) {
	l.Steps = append(l.Steps, newStep(Err, content...))
}
