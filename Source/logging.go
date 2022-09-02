package discordgoself

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

const (
	LogError int = iota
	LogWarning
	LogInformational
	LogDebug
)

var Logger func(msgL, caller int, format string, a ...interface{})

func msglog(msgL, caller int, format string, a ...interface{}) {

	if Logger != nil {
		Logger(msgL, caller, format, a...)
	} else {

		pc, file, line, _ := runtime.Caller(caller)

		files := strings.Split(file, "/")
		file = files[len(files)-1]

		name := runtime.FuncForPC(pc).Name()
		fns := strings.Split(name, ".")
		name = fns[len(fns)-1]

		msg := fmt.Sprintf(format, a...)

		log.Printf("[DG%d] %s:%d:%s() %s\n", msgL, file, line, name, msg)
	}
}

func (s *Session) log(msgL int, format string, a ...interface{}) {

	if msgL > s.LogLevel {
		return
	}

	msglog(msgL, 2, format, a...)
}

func (v *VoiceConnection) log(msgL int, format string, a ...interface{}) {

	if msgL > v.LogLevel {
		return
	}

	msglog(msgL, 2, format, a...)
}
