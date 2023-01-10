package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

// GetUrlPrefix elaborates the current URL prefix in order to point to this server.
//
// e.g. if the request is received from the server running as http://localhost:12345/ but
// it detects to be under a reverse proxy, that is adjusted.
func GetUrlPrefix(r *http.Request, logger logrus.FieldLogger) string {
	return fmt.Sprintf("%s://%s", getUrlProtocol(r, logger), r.Host)
}

// getUrlProtocol detects the current protocol this request is using,
// even if behind a reverse proxy.
func getUrlProtocol(r *http.Request, logger logrus.FieldLogger) string {
	// First of all, check if there's a reverse proxy header with such information
	reverseProxyProto, isReverseProxy := r.Header["X-Forwarded-Proto"]
	if isReverseProxy && len(reverseProxyProto) > 1 {
		// There's something strange here, log that
		logger.Warn("Multiple X-Forwarded-Proto headers found, but only the first one is used.")
		for i, headerValue := range reverseProxyProto {
			logger.Warnf("> X-Forwarded-Proto[%d] = %s", i, headerValue)
		}
	}
	if isReverseProxy && len(reverseProxyProto) > 0 {
		return reverseProxyProto[0] // Use this protocol
	}

	// It doesn't seem like we are behind a reverse proxy,
	// so we can detect HTTPS using TLS status.
	if r.TLS == nil {
		return "http"
	} else {
		return "https"
	}
}
