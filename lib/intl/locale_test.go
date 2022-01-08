package intl

import (
	"fmt"
	"testing"
)

func TestLocale(t *testing.T) {
	AddLocale("en-GB", "en-US", "it-IT", "fa-IR")
	fmt.Println(GetLocale("en"))
	fmt.Println(GuessLocale("en"))

}
