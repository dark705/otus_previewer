package web

import (
	"context"
	"fmt"
	"github.com/dark705/otus_previewer/internal/dispatcher"
	"net/http"
	"time"

	"github.com/dark705/otus_previewer/internal/helpers"
	"github.com/sirupsen/logrus"
)

type Server struct {
	c  Config
	l  *logrus.Logger
	sd *dispatcher.StorageDispatcher
	ws *http.Server
}

type Config struct {
	HttpListen string
}

func NewServer(conf Config, log *logrus.Logger, sd *dispatcher.StorageDispatcher) Server {

	return Server{
		c:  conf,
		l:  log,
		sd: sd,
		ws: &http.Server{Addr: conf.HttpListen, Handler: logRequest(ServeHTTP, log, sd)},
	}
}

func (s *Server) RunServer() {
	go func() {
		s.l.Infoln("Start HTTP server: ", s.c.HttpListen)
		err := s.ws.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			helpers.FailOnError(err, "Fail start HTTP Server")
		}
	}()

}

func (s *Server) Shutdown() {
	s.l.Infoln("Shutdown HTTP server... ")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := s.ws.Shutdown(ctx)
	if err != nil {
		s.l.Errorln("Fail Shutdown HTTP server")
		return
	}
	s.l.Infoln("Success Shutdown HTTP server")
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Previewer")
	//_, _ = w.Write([]byte("Hello world"))
}

//middleware logger
func logRequest(h http.HandlerFunc, l *logrus.Logger, sd *dispatcher.StorageDispatcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l.Infoln(fmt.Sprintf("Income request: %s %s %s", r.RemoteAddr, r.Method, r.URL))
		p, err := ParseUrl(r.URL)
		l.Debugln("parser:", p, err)

		uniqId := GenUniqIdForUrl(r.URL)
		l.Infoln(fmt.Sprintf("Generate uniqId: %s for Url: %s", uniqId, r.URL.Path))
		l.Debugln("Storage Dispatcher usage:", sd.Usage())

		if sd.Exist(uniqId) {
			l.Infoln(fmt.Sprintf("Content for uniqId: %s, found in cache", uniqId))
			cont, _ := sd.Get(uniqId)
			w.Write(cont)
			return
		}
		l.Infoln(fmt.Sprintf("Content for uniqId: %s, not found in cache", uniqId))
		cont, err := GetContext("https://"+p.RequestUrl, r.Header, nil)
		if err != nil {
			cont, err = GetContext("http://"+p.RequestUrl, r.Header, nil)
			if err != nil {
				l.Warnln(err.Error())
				//TODO check for error
				w.Write([]byte(err.Error()))
				return
			}
		}

		//TODO check for error
		w.Write(cont)
		sd.Add(uniqId, cont)
		l.Debugln("Storage Dispatcher usage:", sd.Usage())

		h(w, r)
	}
}
