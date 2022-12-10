//go:build webui

package main

import (
	"fmt"
	"github.com/simonesestito/wasaphoto/webui"
	"io/fs"
	"net/http"
	"strings"
)

func registerWebUI(hdl http.Handler) (http.Handler, error) {
	distDirectory, err := fs.Sub(webui.Dist, "dist")
	if err != nil {
		return nil, fmt.Errorf("error embedding WebUI dist/ directory: %w", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/dashboard/") {
			http.StripPrefix("/dashboard/", http.FileServer(http.FS(distDirectory))).ServeHTTP(w, r)
		} else if r.URL.Path == "/" {
			http.Redirect(w, r, "/dashboard/", http.StatusPermanentRedirect)
		} else {
			hdl.ServeHTTP(w, r)
		}
	}), nil
}
