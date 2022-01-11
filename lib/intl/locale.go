package intl

import (
	"fmt"
	"github.com/getevo/evo-ng"
	"golang.org/x/text/language"
)

var defaultLocale = "en-US"
var supported = []language.Tag{language.Make(defaultLocale)}
var matcher = language.NewMatcher(supported)

type TranslationMap struct {
	Tag     language.Tag
	Entries map[string]Entry
}
type Entry struct {
	Singular string
	Plural   struct {
		Zero string
		One  string
		Two  string
		Few  string
	}
}

func SetDefaultLocale(locale interface{}) {
	switch item := locale.(type) {
	case language.Tag:
		supported[0] = item
	case string:
		supported[0] = language.Make(item)
	default:
		evo.Panic("invalid locale type")
	}
	matcher = language.NewMatcher(supported)
}

func AddLocale(locals ...interface{}) {
	for _, item := range locals {
		switch local := item.(type) {
		case language.Tag:
			supported = append(supported, local)
		case string:
			supported = append(supported, language.Make(local))
		default:
			evo.Panic("invalid locale type")
		}
	}
	matcher = language.NewMatcher(supported)
}

func GetLocale(locale string) (language.Tag, error) {
	for _, item := range supported {
		if locale == item.String() {
			return item, nil
		}
	}
	return language.Make(locale), fmt.Errorf("invalid locale")
}

func GuessLocale(locale string) language.Tag {
	tag, _, _ := matcher.Match(language.Make(locale))
	return tag
}

func Locales() []language.Tag {
	return supported
}
