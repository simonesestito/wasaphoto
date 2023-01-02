package main

import (
	"github.com/simonesestito/wasaphoto/service/storage"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
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

// mustWriteToStorage checks if the specified Storage works, and it's read-write.
//
// This test was added to be sure this feature works,
// even in special Docker situations (someone can forget to mount a volume).
//
// As the name "must" says, in case of an error, it'll be logged as Fatal
// and the execution will stop.
func mustWriteToStorage(storageFs storage.Storage, logger logrus.FieldLogger) {
	logger.Debugln("performing write Storage test")

	var testData = []byte{0xFF, 0x00, 0x01}
	const testFileName = "test.hex"

	locationUrl, err := storageFs.SaveFile(testFileName, testData)
	if err != nil {
		cwd, _ := os.Getwd()
		fsPath := cwd + "/" + storageFs.GetRoot()
		logger.WithError(err).Fatalln("unable to write files: make sure to mount a writable volume to", fsPath)
	}

	logger.Debugln("successfully written test data with locationUrl =", locationUrl)

	if err = storageFs.DeleteFile(testFileName); err != nil {
		logger.WithError(err).Fatalln("unable to delete previously written test file, available at", locationUrl)
	}
}
