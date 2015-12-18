package dogo

import (
	"fmt"
	"log"
	"os"
)

type MyLoger struct {
	Handler *log.Logger
}

func NewLoger() *MyLoger {
	return &MyLoger{
		Handler: log.New(os.Stdout, "[", log.Ldate|log.Ltime),
	}
}

func (l *MyLoger) D(infos ...interface{}) {
	l.output("DEBUG", infos)
}

func (l *MyLoger) I(infos ...interface{}) {
	l.output("INFOS", infos)
}

func (l *MyLoger) E(infos ...interface{}) {
	l.output("ERROR", infos)
}

func (l *MyLoger) W(infos ...interface{}) {
	l.output("WARNG", infos)
}

func (l *MyLoger) output(tag string, infos ...interface{}) {
	var s string
	for _, v := range infos {
		s += fmt.Sprintf("%v", v)
	}
	l.Handler.Println(fmt.Sprintf(`] [%s] - "%s"`, tag, s))
}
