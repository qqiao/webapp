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

package webapp // import "github.com/qqiao/webapp"

import (
	"html/template"
	"net/http"
	"strings"
	"sync"
)

var templateCache sync.Map

// DetermineLocale tries to find the best locale for a given input from a list
// of available locales. It applies the following logic:
//    1. If the input can be directly found in list of available locales, it
//       will be returned directly.
//    2. If a less specific value is available, then it will be returned. E.g.
//       if the input is zh-tw but zh is available then zh will be returned
//    3. Otherwise an empty string is returned.
func DetermineLocale(input string, available []string) string {
	candidate := ""
	for _, locale := range available {
		if equalOrLessSpecific(locale, input) && len(locale) > len(candidate) {
			candidate = locale
		}
	}
	return candidate
}

// DetermineLocaleWithDefault tries to find the best locale for a given input
//from a list of available locales. It applies the following logic:
//    1. If the input can be directly found in list of available locales, it
//       will be returned directly.
//    2. If a less specific value is available, then it will be returned. E.g.
//       if the input is zh-tw but zh is available then zh will be returned
//    3. If no locale can be found, the first entry in available will be
//       returned. If available is empty, an empty string will be returned. It
//       is down to the caller of the function to handle such situation.
func DetermineLocaleWithDefault(input string, available []string) string {
	candidate := DetermineLocale(input, available)
	if candidate == "" && len(available) > 0 {
		return available[0]
	}
	return candidate
}

func equalOrLessSpecific(locale string, input string) bool {
	prefix := locale
	if len(input) > len(locale) {
		prefix = prefix + "-"
	}
	return strings.HasPrefix(input, prefix)
}

// GetTemplate loads the template from the given path.
// The funtion caches the loaded template so that the same template would not
// be parsed over and over again unless skipCache is set to true.
//
// Please note this method panics if template.ParseFiles failes in any way.
func GetTemplate(path string, skipCache bool) *template.Template {
	tmpl, has := templateCache.Load(path)
	if !has || skipCache {
		tmpl = template.Must(template.ParseFiles(path))
		templateCache.Store(path, tmpl)
	}
	return tmpl.(*template.Template)
}

// HSTSHandler takes a normal HTTP handler and adds the capability of sending
// HSTS headers.
func HSTSHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security",
			"max-age=63072000; includeSubDomains; preload")
		f(w, r)
	}
}
