package http

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dark705/otus_previewer/internal/dispatcher"
	"github.com/dark705/otus_previewer/internal/image"
	"github.com/sirupsen/logrus"
)

func handlerRequest(logger *logrus.Logger, imageDispatcher *dispatcher.ImageDispatcher, imageLimit int, imageGetTimeout int) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Server", "Previewer")
		logger.Infof("income request: %s %s %s", request.RemoteAddr, request.Method, request.URL)
		parsedURL, err := ParseURL(request.URL)
		if err != nil {
			logger.Warnln(err)
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}

		uniqID := GenUniqIDForURL(request.URL)
		logger.Infof("generate uniq reqId: %s for Url: %s", uniqID, request.URL.Path)

		if handleCached(logger, imageDispatcher, responseWriter, uniqID) {
			return
		}
		handleNoCached(logger, imageDispatcher, imageLimit, imageGetTimeout, parsedURL, uniqID, responseWriter, request)
	}
}

func handleCached(logger *logrus.Logger,
	imageDispatcher *dispatcher.ImageDispatcher,
	responseWriter http.ResponseWriter,
	uniqID string) bool {
	cachedImage, err := imageDispatcher.Get(uniqID)
	if err != nil {
		logger.Errorln(err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return true
	}
	if cachedImage != nil {
		logger.Infof("image for uniqID: %s, found in cache", uniqID)
		_, err = responseWriter.Write(cachedImage)
		if err != nil {
			logger.Errorln(err)
		}
		return true
	}
	return false
}

func handleNoCached(logger *logrus.Logger,
	imageDispatcher *dispatcher.ImageDispatcher,
	imageLimit int,
	imageGetTimeout int,
	parsedURL URLParams,
	uniqID string,
	responseWriter http.ResponseWriter,
	request *http.Request) {
	logger.Infof("image for uniq reqId: %s, not found in cache, need to download", uniqID)
	httpClientContext, httpClientCancel := context.WithTimeout(context.Background(), time.Second*time.Duration(imageGetTimeout))
	defer httpClientCancel()
	remoteResponse, err := makeRequest(httpClientContext, parsedURL.RequestURL, request.Header, nil)
	if err != nil {
		logger.Warnf("error on make request to remote server: %s", err)
		http.Error(responseWriter, err.Error(), http.StatusBadGateway)
		return
	}

	if remoteResponse.StatusCode != http.StatusOK {
		proxyErrorRemoteResponse(logger, remoteResponse, responseWriter, parsedURL)
		return
	}

	handleRemoteResponse(logger, remoteResponse, responseWriter, imageDispatcher, uniqID, imageLimit, parsedURL)
}

func proxyErrorRemoteResponse(logger *logrus.Logger,
	remoteResponse *http.Response,
	responseWriter http.ResponseWriter,
	parsedURL URLParams) {
	bodyBytes, err := ioutil.ReadAll(remoteResponse.Body)
	_ = remoteResponse.Body.Close()
	if err != nil {
		logger.Errorln(err)
	}
	logger.Warnf("remote server for url: %s return status: %d ", parsedURL.RequestURL, remoteResponse.StatusCode)
	for h, v := range remoteResponse.Header {
		responseWriter.Header().Set(h, v[0])
	}
	responseWriter.WriteHeader(remoteResponse.StatusCode)
	_, err = responseWriter.Write(bodyBytes)
	if err != nil {
		logger.Errorln(err)
	}
}

func handleRemoteResponse(logger *logrus.Logger,
	response *http.Response,
	responseWriter http.ResponseWriter,
	imageDispatcher *dispatcher.ImageDispatcher,
	uniqID string,
	imageLimit int,
	parsedURL URLParams) {
	remoteImage, err := image.ReadImageAsByte(response.Body, imageLimit)
	_ = response.Body.Close()
	if err != nil {
		logger.Warnln(err)
		http.Error(responseWriter, err.Error(), http.StatusBadGateway)
		return
	}

	logger.Infof("success download image for uniq reqId: %s,", uniqID)
	convertedImage, err := image.Resize(remoteImage, image.ResizeConfig{
		Action: parsedURL.Service,
		Width:  parsedURL.Width,
		Height: parsedURL.Height,
	})
	if err != nil {
		logger.Errorln(err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = responseWriter.Write(convertedImage)
	if err != nil {
		logger.Errorln(err)
	}

	err = imageDispatcher.Add(uniqID, convertedImage)
	if err != nil {
		logger.Errorln(err)
	}
}
