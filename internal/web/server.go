package web

import (
	"context"
	"fmt"

	"net/http"
	"time"

	"github.com/dark705/otus_previewer/internal/dispatcher"
	"github.com/dark705/otus_previewer/internal/image"

	"github.com/dark705/otus_previewer/internal/helpers"
	"github.com/sirupsen/logrus"
)

type Server struct {
	c  Config
	l  *logrus.Logger
	sd *dispatcher.ImageDispatcher
	ws *http.Server
}

type Config struct {
	HttpListen       string
	ImageMaxFileSize int
}

func NewServer(conf Config, log *logrus.Logger, sd *dispatcher.ImageDispatcher) Server {

	return Server{
		c:  conf,
		l:  log,
		sd: sd,
		ws: &http.Server{Addr: conf.HttpListen, Handler: handlerRequest(log, sd, conf.ImageMaxFileSize)},
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

func handlerRequest(l *logrus.Logger, imDis *dispatcher.ImageDispatcher, imageLimit int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "Previewer")

		//Check and parse request params
		l.Infoln(fmt.Sprintf("Income request: %s %s %s", r.RemoteAddr, r.Method, r.URL))
		p, err := ParseUrl(r.URL)
		if err != nil {
			l.Warnln(err)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				l.Errorln(err)
			}
			return
		}

		//Generate uniq id for request, witch will be used for save image
		uniqId := GenUniqIdForUrl(r.URL)
		l.Infoln(fmt.Sprintf("Generate uniq reqId: %s for Url: %s", uniqId, r.URL.Path))

		//Image found in cache
		if imDis.Exist(uniqId) {
			l.Infoln(fmt.Sprintf("Image for uniqId: %s, found in cache", uniqId))
			resp, err := imDis.Get(uniqId)
			if err != nil {
				l.Errorln(err)
				resp = []byte("Internal server error")
			}
			_, err = w.Write(resp)
			if err != nil {
				l.Errorln(err)
			}
			return
		}

		//Image not fount in cache, need download
		l.Infoln(fmt.Sprintf("Image for uniq reqId: %s, not found in cache, need to dowload", uniqId))
		var im []byte
		var errHttps, errHttp error
		//first try https
		im, errHttps = GetImageAsBytes("https://", p.RequestUrl, r.Header, nil, imageLimit)
		if errHttps != nil {
			l.Warnln(errHttps.Error())
			//if some error, try http
			im, errHttp = GetImageAsBytes("http://", p.RequestUrl, r.Header, nil, imageLimit)
			if errHttp != nil {
				l.Warnln(errHttp.Error())
				_, err := w.Write([]byte(errHttps.Error() + "\n" + errHttp.Error()))
				if err != nil {
					l.Errorln(err)
				}
				return
			}
		}

		//Downloaded image as byte, make convert
		l.Infoln(fmt.Sprintf("Success download image for uniq reqId: %s,", uniqId))
		convertedImage, err := image.Resize(im, image.ResizeConfig{
			Action: p.Service,
			Width:  p.Width,
			Height: p.Height,
		})
		if err != nil {
			l.Errorln(err.Error())
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				l.Errorln(err)
			}
			return
		}

		//send to client
		_, err = w.Write(convertedImage)
		if err != nil {
			l.Errorln(err)
		}

		//save to cache
		err = imDis.Add(uniqId, convertedImage)
		if err != nil {
			l.Errorln(err)
		}
	}
}
