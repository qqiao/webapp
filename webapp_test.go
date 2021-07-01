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
	if DetermineLocale("zh", []string{}) != "" {
		t.Error("DetermineLocale(\"zh\", []string{}) should return \"\"")
	}

	if DetermineLocale("es", available) != "" {
		t.Error("DetermineLocale(\"es\", available) should return empty")
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
