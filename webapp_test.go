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
		t.Errorf("equalOrLessSpecific(\"zh\", \"zh\") should be true")
	}

	if !equalOrLessSpecific("zh", "zh-hk") {
		t.Errorf("equalOrLessSpecific(\"zh\", \"zh-hk\") should be true")
	}

	if equalOrLessSpecific("zh-tw", "zh-hk") {
		t.Errorf("equalOrLessSpecific(\"zh-tw\", \"zh-hk\") should be false")
	}

	if equalOrLessSpecific("en", "zh-hk") {
		t.Errorf("equalOrLessSpecific(\"en\", \"zh-hk\") should be false")
	}
}

func TestDetermineLocale(t *testing.T) {
	if "" != DetermineLocale("zh", []string{}) {
		t.Errorf("DetermineLocale(\"zh\", []string{}) should return empty")
	}

	if "" != DetermineLocale("es", available) {
		t.Errorf("DetermineLocale(\"es\", available) should return empty")
	}
}
