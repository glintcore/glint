package server

import (
	"fmt"
	"log"
	"net/http"
)

type serverLog struct {
	logger   *log.Logger
	pid      uint64
	writePid bool
}

func (l *serverLog) log(format string, v ...interface{}) {

	s := fmt.Sprintf(format, v...)

	var pidPrefix string
	if l.writePid {
		pidPrefix = fmt.Sprintf("[%v] ", l.pid)
	}

	if l.logger == nil {
		log.Printf("%s%s", pidPrefix, s)
	} else {
		l.logger.Printf("%s%s", pidPrefix, s)
	}
}

func (l *serverLog) logRequest(r *http.Request, statusCode int) {

	/*
		var host string
		var err error
		if host, _, err = net.SplitHostPort(r.RemoteAddr); err != nil {
			log.Print(err)
			host = "-"
		}
	*/

	/*
		var status string
		if statusCode != 0 {
			status = fmt.Sprintf(" %v", statusCode)
		}
	*/

	//l.log("%s \"%s %s\"%s", host, r.Method, r.URL, status)
	//l.log("%s \"%s %s\"", host, r.Method, r.URL)
	//l.log("%s \"%s %s\"", host, r.Method, r.RequestURI)
	l.log("%s \"%s %s\"", r.RemoteAddr, r.Method, r.RequestURI)
}
