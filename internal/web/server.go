package web

import (
	"context"
	"fmt"
	"io/ioutil"
	"sync"

	"net/http"
	"time"

	"github.com/dark705/otus_previewer/internal/dispatcher"
	"github.com/dark705/otus_previewer/internal/image"

	"github.com/dark705/otus_previewer/internal/helpers"
	"github.com/sirupsen/logrus"
)

type Server struct {
	c       Config
	l       *logrus.Logger
	ws      *http.Server
	imgDisp *dispatcher.ImageDispatcher
}

type Config struct {
	HttpListen       string
	ImageMaxFileSize int
}

func NewServer(conf Config, log *logrus.Logger, imageDispatcher *dispatcher.ImageDispatcher) Server {
	m := sync.Mutex{}
	return Server{
		c:       conf,
		l:       log,
		ws:      &http.Server{Addr: conf.HttpListen, Handler: handlerRequest(log, imageDispatcher, conf.ImageMaxFileSize, m)},
		imgDisp: imageDispatcher,
	}
}

func (s *Server) RunServer() {
	go func() {
		s.l.Infoln("Start HTTP server:", s.c.HttpListen)
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

func handlerRequest(l *logrus.Logger, imDis *dispatcher.ImageDispatcher, imageLimit int, mu sync.Mutex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "Previewer")

		//Check and parse request params
		l.Infoln(fmt.Sprintf("Income request: %s %s %s", r.RemoteAddr, r.Method, r.URL))
		p, err := ParseUrl(r.URL)
		if err != nil {
			l.Warnln(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//Generate uniq id for request, witch will be used for save image
		uniqId := GenUniqIdForUrl(r.URL)
		l.Infoln(fmt.Sprintf("Generate uniq reqId: %s for Url: %s", uniqId, r.URL.Path))

		//Image found in cache
		mu.Lock()
		if imDis.Exist(uniqId) {
			l.Infoln(fmt.Sprintf("Image for uniqId: %s, found in cache", uniqId))
			resp, err := imDis.Get(uniqId)
			mu.Unlock()
			if err != nil {
				l.Errorln(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = w.Write(resp)
			if err != nil {
				l.Errorln(err)
			}
			return
		}
		mu.Unlock()

		//Image not fount in cache, need download
		l.Infoln(fmt.Sprintf("Image for uniq reqId: %s, not found in cache, need to dowload", uniqId))
		//first try https
		resp, err := makeRequest("https://", p.RequestUrl, r.Header, nil, imageLimit)
		if err != nil {
			l.Warnln(err)
			//if some error, try http
			resp, err = makeRequest("http://", p.RequestUrl, r.Header, nil, imageLimit)
			if err != nil {
				l.Warnln(err)
				http.Error(w, err.Error(), http.StatusBadGateway)
				return
			}
		}
		//If remote server response not StatusOk, proxy response to client with status, headers and body
		if resp.StatusCode != http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if err != nil {
				l.Errorln(err)
			}
			l.Warnln(fmt.Sprintf("Remote server for url: %s return status: %d ", p.RequestUrl, resp.StatusCode))
			for h, v := range resp.Header {
				w.Header().Set(h, v[0])
			}
			w.WriteHeader(resp.StatusCode)
			_, err = w.Write(bodyBytes)
			if err != nil {
				l.Errorln(err)
			}
			return
		}

		//Status Ok, read response as image
		im, err := image.ReadImageAsByte(resp.Body, imageLimit)
		_ = resp.Body.Close()
		if err != nil {
			l.Warnln(err)
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		//Downloaded image as byte, make convert
		l.Infoln(fmt.Sprintf("Success download image for uniq reqId: %s,", uniqId))
		convertedImage, err := image.Resize(im, image.ResizeConfig{
			Action: p.Service,
			Width:  p.Width,
			Height: p.Height,
		})
		if err != nil {
			l.Errorln(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//send to client
		_, err = w.Write(convertedImage)
		if err != nil {
			l.Errorln(err)
		}

		//save to cache
		mu.Lock()
		err = imDis.Add(uniqId, convertedImage)
		mu.Unlock()
		if err != nil {
			l.Errorln(err)
		}
	}
}
