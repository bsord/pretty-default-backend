/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	// FormatHeader name of the header used to extract the format
	FormatHeader = "X-Format"

	// CodeHeader name of the header used as source of the HTTP status code to return
	CodeHeader = "X-Code"

	// ContentType name of the header that defines the format of the reply
	ContentType = "Content-Type"

	// OriginalURI name of the header with the original URL from NGINX
	OriginalURI = "X-Original-URI"

	// Namespace name of the header that contains information about the Ingress namespace
	Namespace = "X-Namespace"

	// IngressName name of the header that contains the matched Ingress
	IngressName = "X-Ingress-Name"

	// ServiceName name of the header that contains the matched Service in the Ingress
	ServiceName = "X-Service-Name"

	// ServicePort name of the header that contains the matched Service port in the Ingress
	ServicePort = "X-Service-Port"

	// RequestId is a unique ID that identifies the request - same as for backend service
	RequestId = "X-Request-ID"

	// ErrFilestemplatePathVar is the name of the environment variable indicating
	// the location on disk of files served by the handler.
	ErrFilestemplatePathVar = "ERROR_FILES_templatePath"
)

func init() {
	prometheus.MustRegister(requestCount)
	prometheus.MustRegister(requestDuration)
}

func main() {
	errFilestemplatePath := "/app/www"
	if os.Getenv(ErrFilestemplatePathVar) != "" {
		errFilestemplatePath = os.Getenv(ErrFilestemplatePathVar)
	}

	http.HandleFunc("/", errorHandler(errFilestemplatePath))

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.ListenAndServe(fmt.Sprintf(":8080"), nil)
}

func errorHandler(templatePath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		var ext string

		if os.Getenv("DEBUG") != "" {
			w.Header().Set(FormatHeader, r.Header.Get(FormatHeader))
			w.Header().Set(CodeHeader, r.Header.Get(CodeHeader))
			w.Header().Set(ContentType, r.Header.Get(ContentType))
			w.Header().Set(OriginalURI, r.Header.Get(OriginalURI))
			w.Header().Set(Namespace, r.Header.Get(Namespace))
			w.Header().Set(IngressName, r.Header.Get(IngressName))
			w.Header().Set(ServiceName, r.Header.Get(ServiceName))
			w.Header().Set(ServicePort, r.Header.Get(ServicePort))
			w.Header().Set(RequestId, r.Header.Get(RequestId))
		}

		// Get request format (html vs json)
		format := r.Header.Get(FormatHeader)
		log.Printf(format)
		if format == "" {
			format = "text/html"
			log.Printf("format not specified. Using %v", format)
		}

		switch format {
		case "text/html; charset=utf-8":
			ext = ".html"
		case "text/html":
			ext = ".html"
		case "application/json":
			ext = ".json"
		default:
			format = "text/html"
			ext = ".html"
		}

		// Set format response header based on request content type
		w.Header().Set(ContentType, format)

		// Get Error code passed in from ingress
		errCode := r.Header.Get(CodeHeader)
		code, err := strconv.Atoi(errCode)
		if err != nil {
			code = 404
			log.Printf("unexpected error reading return code: %v. Using %v", err, code)
		}

		// Set error code response header to header passed in from ingress
		w.WriteHeader(code)

		// Parse template file based on content type
		file := fmt.Sprintf("%v/%v%v", templatePath, "error", ext)
		t := template.Must(template.New(path.Base(file)).Funcs(template.FuncMap{"safeCSS": func(css string) template.CSS { return template.CSS(css) }}).ParseFiles(file))

		// Variable replacement struct
		data := struct {
			BGColor   string
			Branding  string
			ErrorCode int
			Origin    string
		}{os.Getenv("BG_COLOR"), os.Getenv("BRANDING_TEXT"), code, r.Header.Get(OriginalURI)}

		// serve templated html
		err = t.Execute(w, data)
		if err != nil {
			log.Printf("There was a problem processing the template: %v", err)
			// if error serve basic error response
			http.NotFound(w, r)
			return
		}

		duration := time.Now().Sub(start).Seconds()

		proto := strconv.Itoa(r.ProtoMajor)
		proto = fmt.Sprintf("%s.%s", proto, strconv.Itoa(r.ProtoMinor))

		requestCount.WithLabelValues(proto).Inc()
		requestDuration.WithLabelValues(proto).Observe(duration)
	}
}
