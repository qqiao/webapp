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

// +build appengine

package webapp

import (
	"fmt"
	"io"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

var isDev = appengine.IsDevAppServer()

func initPolyserveProxy(mux *http.ServeMux, URL string) error {
	for _, path := range PolyserveURLs {
		mux.HandleFunc(path, makeProxy(URL))
	}

	return nil
}

func makeProxy(URL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		client := urlfetch.Client(ctx)

		url := fmt.Sprintf("%s%s", URL, r.URL.Path)

		req, err := http.NewRequest(r.Method, url, nil)
		if nil != err {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header.Set("user-agent", r.Header.Get("User-Agent"))

		resp, err := client.Do(req)
		if nil != err {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for k, vv := range resp.Header {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}

		io.Copy(w, resp.Body)
	}
}
