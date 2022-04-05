// Copyright 2022 Qian Qiao
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
	"testing"
)

var available = []string{
	"en",
	"en-us",
	"zh",
	"zh-TW",
}

func TestEqualOrLessSpecific(t *testing.T) {
	if !equalOrLessSpecific("zh", "zh") {
		t.Error("equalOrLessSpecific(\"zh\", \"zh\") should be true")
	}

	if !equalOrLessSpecific("zh", "zh-hk") {
		t.Error("equalOrLessSpecific(\"zh\", \"zh-hk\") should be true")
	}

	if equalOrLessSpecific("zh", "zh1123") {
		t.Error("equalOrLessSpecific(\"zh\", \"zh1123\") should be false")
	}

	if equalOrLessSpecific("zh-tw", "zh-hk") {
		t.Error("equalOrLessSpecific(\"zh-tw\", \"zh-hk\") should be false")
	}

	if equalOrLessSpecific("en", "zh-hk") {
		t.Error("equalOrLessSpecific(\"en\", \"zh-hk\") should be false")
	}
}

func TestDetermineLocale(t *testing.T) {
	if DetermineLocale("zh", []string{}) != "" {
		t.Error("DetermineLocale(\"zh\", []string{}) should return \"\"")
	}

	if DetermineLocale("es", available) != "" {
		t.Error("DetermineLocale(\"es\", available) should return empty")
	}

	if DetermineLocale("zh-TW", available) != "zh-TW" {
		t.Error("DetermineLocale(\"zh-TW\", available) should return \"zh-TW\"")
	}
}

func TestDetermineLocaleWithDefault(t *testing.T) {
	if DetermineLocaleWithDefault("zh", []string{}) != "" {
		t.Error("DetermineLocaleWithDefault(\"zh\", []string{}) should return \"\"")
	}

	if DetermineLocaleWithDefault("es", available) != available[0] {
		t.Error("DetermineLocaleWithDefault(\"es\", available) should return available[0]")
	}
}
