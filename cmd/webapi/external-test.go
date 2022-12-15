package main

import (
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
)

// mustPerformDnsQuery checks if DNS queries work.
//
// This test was added to be sure this feature works,
// even in special Docker situations (solved with the flag -tags netgo back then).
//
// As the name "must" says, in case of an error, it'll be logged as Fatal
// and the execution will stop.
func mustPerformDnsQuery(host string, logger logrus.FieldLogger) {
	logger.Debugf("testing DNS queries using %s\n", host)

	dnsResponse, err := net.LookupIP(host)
	if err != nil {
		logger.WithError(err).Fatalln("error testing DNS")
	}

	if len(dnsResponse) == 0 {
		logger.Fatalf("DNS response for %s is empty\n", host)
	}

	logger.Debugf("DNS response for %s contains %d results, first one is %s\n", host, len(dnsResponse), dnsResponse[0].String())
}

// mustPerformHttpsRequest checks if HTTPS requests work.
//
// This test was added to be sure this feature works,
// even in special Docker situations (solved with the flag -tags netgo back then).
//
// As the name "must" says, in case of an error, it'll be logged as Fatal
// and the execution will stop.
func mustPerformHttpsRequest(logger logrus.FieldLogger) {
	logger.Debugln("performing tests to ensure third-party services will be reachable")
	const host, httpsUrlString = "api.tinypng.com", "https://api.tinypng.com/"
	mustPerformDnsQuery(host, logger)

	logger.Debugf("testing HTTPS requests with %s\n", httpsUrlString)
	response, err := http.DefaultClient.Get(httpsUrlString)
	if err != nil {
		logger.WithError(err).Fatalln("error performing test HTTPS request")
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		logger.Fatalf("HTTPS test response has status %s\n", response.Status)
	}

	logger.Debugf("HTTPS test response has status %s\n", response.Status)
}
