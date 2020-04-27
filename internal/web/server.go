package web

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dark705/otus_previewer/internal/dispatcher"
	"github.com/dark705/otus_previewer/internal/image"

	"github.com/dark705/otus_previewer/internal/helpers"
	"github.com/sirupsen/logrus"
)

type Server struct {
	config           Config
	logger           *logrus.Logger
	httpServer       *http.Server
	imgageDispatcher *dispatcher.ImageDispatcher
}

type Config struct {
	HTTPListen       string
	ImageMaxFileSize int
}

func NewServer(conf Config, log *logrus.Logger, imageDispatcher *dispatcher.ImageDispatcher) Server {
	return Server{
		config:           conf,
		logger:           log,
		httpServer:       &http.Server{Addr: conf.HTTPListen, Handler: handlerRequest(log, imageDispatcher, conf.ImageMaxFileSize)},
		imgageDispatcher: imageDispatcher,
	}
}

func (s *Server) RunServer() {
	go func() {
		s.logger.Infoln("start HTTP server:", s.config.HTTPListen)
		err := s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			helpers.FailOnError(err, "fail start HTTP Server")
		}
	}()
}

func (s *Server) Shutdown() {
	s.logger.Infoln("shutdown HTTP server... ")
	ctx, chancel := context.WithTimeout(context.Background(), time.Second*10)
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		s.logger.Errorln("fail Shutdown HTTP server")
		chancel()
		return
	}
	s.logger.Infoln("success Shutdown HTTP server")
	chancel()
}

func handlerRequest(l *logrus.Logger, imDis *dispatcher.ImageDispatcher, imageLimit int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "Previewer")

		//Check and parse request params
		l.Infoln(fmt.Sprintf("income request: %s %s %s", r.RemoteAddr, r.Method, r.URL))
		p, err := ParseURL(r.URL)
		if err != nil {
			l.Warnln(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//Generate uniq id for request, witch will be used for save image
		uniqID := GenUniqIDForURL(r.URL)
		l.Infoln(fmt.Sprintf("generate uniq reqId: %s for Url: %s", uniqID, r.URL.Path))

		//Image found in cache
		cachedImage, err := imDis.Get(uniqID)
		if err != nil {
			l.Errorln(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if cachedImage != nil {
			l.Infoln(fmt.Sprintf("image for uniqID: %s, found in cache", uniqID))
			_, err = w.Write(cachedImage)
			if err != nil {
				l.Errorln(err)
			}
			return
		}

		//Image not fount in cache, need download
		l.Infoln(fmt.Sprintf("image for uniq reqId: %s, not found in cache, need to dowload", uniqID))
		//first try https
		resp, err := makeRequest("https://", p.RequestURL, r.Header, nil)
		if err != nil {
			l.Warnln(err)
			//if some error, try http
			resp, err = makeRequest("http://", p.RequestURL, r.Header, nil)
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
			l.Warnln(fmt.Sprintf("remote server for url: %s return status: %d ", p.RequestURL, resp.StatusCode))
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
		l.Infoln(fmt.Sprintf("success download image for uniq reqId: %s,", uniqID))
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
		err = imDis.Add(uniqID, convertedImage)
		if err != nil {
			l.Errorln(err)
		}
	}
}
