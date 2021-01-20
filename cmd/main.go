package main

import (
	"flag"

	"github.com/yaou-li/go-i18n"
	"github.com/yaou-li/go-i18n/cmd/extract"
)

var (
	extra     = flag.Bool("extract", false, "extract from source file")
	src       = flag.String("src", ".", "set the golang src directory")
	clean     = flag.Bool("clean", false, "clear all data")
	namespace = flag.Bool("namespace", false, "use namespace mode")
)

func main() {
	flag.Parse()
	if *extra {
		opts := i18n.NewI18nOpts()
		opts.SetLanguageDir("./i18n")
		opts.SetTargetLang("zh")
		opts.SetEnableNamespace(*namespace)

		ex := extract.NewExtractor(opts)
		ex.Extract(*src, *clean)
	}
}
