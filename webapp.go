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

package webapp

import (
	"html/template"
	"net/http"
)

// PolyserveURLs are the URLs we should proxy to polymer serve
var PolyserveURLs = []string{
	"/node_modules/",
	"/src/",
}

var templateCache = make(map[string]*template.Template)

// GetTemplate loads the template from the given path.
// The funtion caches the loaded template so that the same template would not
// be parsed over and over again unless skipCache is set to true.
//
// Please note this method panics if template.ParseFiles failes in any way.
func GetTemplate(path string, skipCache bool) *template.Template {
	tmpl, has := templateCache[path]
	if !has || !skipCache {
		tmpl = template.Must(template.ParseFiles(path))
		templateCache[path] = tmpl
	}
	return tmpl
}

// HSTSHandler takes a normal HTTP handler and adds the capability of sending
// HSTS headers.
func HSTSHandler(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security",
			"max-age=63072000; includeSubDomains; preload")
		f(w, r)
	})
}
