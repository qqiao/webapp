// Copyright 2017 Qian Qiao
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build local

package webapp

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

var isDev = true

// PolyserveURLs are the URLs we should proxy to polymer serve
var PolyserveURLs = []string{
	"/manifest.json",
	"/service-worker.js",
	"/sw.js",
	"/node_modules/",
	"/src/",
}

// InitPolyserveProxy initializes a proxy for 'polymer serve'
func InitPolyserveProxy(mux *http.ServeMux, URL string) error {
	backend, err := url.Parse(URL)
	if nil != err {
		return err
	}

	for _, path := range PolyserveURLs {
		mux.Handle(path, httputil.NewSingleHostReverseProxy(backend))
	}

	return nil
}
