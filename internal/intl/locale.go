package intl

import (
	"golang.org/x/text/language"
)

var supported []language.Tag
var matcher = language.NewMatcher(supported)

func AddLocale(locals ...interface{}) {
	for _, item := range locals {
		switch local := item.(type) {
		case language.Tag:
			supported = append(supported, local)
		case string:
			supported = append(supported, language.Make(local))
		}
	}
	matcher = language.NewMatcher(supported)
}

func GetLocale(locale string) language.Tag {
	return language.Make(locale)
}

func GuessLocale(locale string) language.Tag {
	tag, _, _ := matcher.Match(language.Make(locale))
	return tag
}

func Locals() []language.Tag {
	return supported
}
