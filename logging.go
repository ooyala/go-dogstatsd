package dogstatsd

import "github.com/blackjack/syslog"
import "fmt"
import "log"
import "os"
import "path"

func init() {
	// Get the program base name
	cmdline := os.Args[0]
	program := path.Base(cmdline)
	syslog.Openlog(program, syslog.LOG_PID, syslog.LOG_USER)
}

type Log_ struct {
	systemlog bool
}

// Intialize Log_ type
func LogStart(slog bool) Log_ {
	return Log_{systemlog: slog}
}

func Logg(slog bool, priority syslog.Priority, msg string) {
	if slog != true {
		log.Printf("%s", msg)
		return
	}

	syslog.Syslog(priority, msg)
}

func Logf(slog bool, priority syslog.Priority, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	if slog != true {
		log.Printf("%s", msg)
		return
	}

	syslog.Syslog(priority, msg)
}

func (l *Log_) Emerg(msg string) {
	Logg(l.systemlog, syslog.LOG_EMERG, msg)
}

func (l *Log_) Emergf(format string, a ...interface{}) {
	Logf(l.systemlog, syslog.LOG_EMERG, format, a...)
}

func (l *Log_) Alert(msg string) {
	Logg(l.systemlog, syslog.LOG_ALERT, msg)
}

func (l *Log_) Alertf(format string, a ...interface{}) {
	Logf(l.systemlog, syslog.LOG_ALERT, format, a...)
}

func (l *Log_) Crit(msg string) {
	Logg(l.systemlog, syslog.LOG_CRIT, msg)
}

func (l *Log_) Critf(format string, a ...interface{}) {
	Logf(l.systemlog, syslog.LOG_CRIT, format, a...)
}

func (l *Log_) Err(msg string) {
	Logg(l.systemlog, syslog.LOG_ERR, msg)
}

func (l *Log_) Errf(format string, a ...interface{}) {
	Logf(l.systemlog, syslog.LOG_ERR, format, a...)
}

func (l *Log_) Warning(msg string) {
	Logg(l.systemlog, syslog.LOG_WARNING, msg)
}

func (l *Log_) Warningf(format string, a ...interface{}) {
	Logf(l.systemlog, syslog.LOG_WARNING, format, a...)
}

func (l *Log_) Notice(msg string) {
	Logg(l.systemlog, syslog.LOG_NOTICE, msg)
}

func (l *Log_) Noticef(format string, a ...interface{}) {
	Logf(l.systemlog, syslog.LOG_NOTICE, format, a...)
}

func (l *Log_) Info(msg string) {
	Logg(l.systemlog, syslog.LOG_INFO, msg)
}

func (l *Log_) Infof(format string, a ...interface{}) {
	Logf(l.systemlog, syslog.LOG_INFO, format, a...)
}

func (l *Log_) Debug(msg string) {
	Logg(l.systemlog, syslog.LOG_DEBUG, msg)
}

func (l *Log_) Debugf(format string, a ...interface{}) {
	Logf(l.systemlog, syslog.LOG_DEBUG, format, a...)
}
