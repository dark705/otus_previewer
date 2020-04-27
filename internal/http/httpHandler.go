package http

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dark705/otus_previewer/internal/dispatcher"
	"github.com/dark705/otus_previewer/internal/image"
	"github.com/sirupsen/logrus"
)

func handlerRequest(logger *logrus.Logger, imageDispatcher *dispatcher.ImageDispatcher, imageLimit int) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Server", "Previewer")

		//Check and parse request params
		logger.Infoln(fmt.Sprintf("income request: %s %s %s", request.RemoteAddr, request.Method, request.URL))
		parsedURL, err := ParseURL(request.URL)
		if err != nil {
			logger.Warnln(err)
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}

		//Generate uniq id for request, witch will be used for save image
		uniqID := GenUniqIDForURL(request.URL)
		logger.Infoln(fmt.Sprintf("generate uniq reqId: %s for Url: %s", uniqID, request.URL.Path))

		if handleCached(logger, imageDispatcher, responseWriter, uniqID) {
			return
		}
		handleNoCached(logger, imageDispatcher, imageLimit, parsedURL, uniqID, responseWriter, request)
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
		logger.Infoln(fmt.Sprintf("image for uniqID: %s, found in cache", uniqID))
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
	parsedURL URLParams,
	uniqID string,
	responseWriter http.ResponseWriter,
	request *http.Request) {
	logger.Infoln(fmt.Sprintf("image for uniq reqId: %s, not found in cache, need to dowload", uniqID))
	response, err := makeRequest(parsedURL.RequestURL, request.Header, nil)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadGateway)
		return
	}
	//If remote server response not StatusOk, proxy response to client with status, headers and body
	if response.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		_ = response.Body.Close()
		if err != nil {
			logger.Errorln(err)
		}
		logger.Warnln(fmt.Sprintf("remote server for url: %s return status: %d ", parsedURL.RequestURL, response.StatusCode))
		for h, v := range response.Header {
			responseWriter.Header().Set(h, v[0])
		}
		responseWriter.WriteHeader(response.StatusCode)
		_, err = responseWriter.Write(bodyBytes)
		if err != nil {
			logger.Errorln(err)
		}
		return
	}

	//Status Ok, read response as image
	im, err := image.ReadImageAsByte(response.Body, imageLimit)
	_ = response.Body.Close()
	if err != nil {
		logger.Warnln(err)
		http.Error(responseWriter, err.Error(), http.StatusBadGateway)
		return
	}

	//Downloaded image as byte, make convert
	logger.Infoln(fmt.Sprintf("success download image for uniq reqId: %s,", uniqID))
	convertedImage, err := image.Resize(im, image.ResizeConfig{
		Action: parsedURL.Service,
		Width:  parsedURL.Width,
		Height: parsedURL.Height,
	})
	if err != nil {
		logger.Errorln(err)
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	//send to client
	_, err = responseWriter.Write(convertedImage)
	if err != nil {
		logger.Errorln(err)
	}

	//save to cache
	err = imageDispatcher.Add(uniqID, convertedImage)
	if err != nil {
		logger.Errorln(err)
	}
}
