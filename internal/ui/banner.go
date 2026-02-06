package ui

import _ "embed"

//go:embed assets/banner.txt
var banner string

func Banner() string {
	return banner
}
