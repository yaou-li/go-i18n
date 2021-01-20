package main

import (
	"flag"

	"github.com/yaou-li/go-i18n"
)

var (
	src       = flag.String("src", ".", "set the golang src directory")
	clean     = flag.Bool("clean", false, "clear all data")
	namespace = flag.Bool("namespace", false, "use namespace mode")
)

func main() {
	flag.Parse()

	opts := i18n.NewI18nOpts()
	opts.SetLanguageDir("./i18n")
	opts.SetTargetLang("en")
	opts.SetEnableNamespace(*namespace)

	ex := NewExtractor(opts)
	ex.Extract(*src, *clean)
}
