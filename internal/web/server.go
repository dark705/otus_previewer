package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dark705/otus_previewer/internal/helpers"
	"github.com/sirupsen/logrus"
)

type Server struct {
	c  Config
	l  *logrus.Logger
	ws *http.Server
}

type Config struct {
	HttpListen string
}

func NewServer(conf Config, log *logrus.Logger) Server {

	return Server{
		c:  conf,
		l:  log,
		ws: &http.Server{Addr: conf.HttpListen, Handler: logRequest(ServeHTTP, log)},
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
func logRequest(h http.HandlerFunc, l *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l.Infoln(fmt.Sprintf("Income request: %s %s %s", r.RemoteAddr, r.Method, r.URL))
		p, err := ParseUrl(r.URL)
		l.Debugln("parser:", p, err)

		uniqId := GenUniqIdForUrl(r.URL)

		l.Infoln(fmt.Sprintf("Generate uniqId: %s for Url: %s", uniqId, r.URL.Path))

		resp, err := GetContext("https://"+p.RequestUrl, r.Header, nil)
		if err != nil {
			resp, err = GetContext("http://"+p.RequestUrl, r.Header, nil)
			if err != nil {
				l.Warnln(err.Error())
				//TODO check for error
				w.Write([]byte(err.Error()))

				return
			}
		}
		//TODO check for error
		w.Write(resp)

		h(w, r)
	}
}