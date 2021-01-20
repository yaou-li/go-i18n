package i18n

import (
	"strings"
)

/**
* extract reader
**/
type I18nReader interface {
	ReadAllFile() error
}

type reader struct {
	opts   *I18nOpts
	dicts  map[string]*I18nDict
	parser I18nParser
}

func NewReader(opts *I18nOpts, dicts map[string]*I18nDict) I18nReader {
	return &reader{
		opts:   opts,
		dicts:  dicts,
		parser: ParserFactory(opts),
	}
}

func (r *reader) ReadAllFile() error {
	var s []string
	files, err := ReadAllPath(r.opts.dir, s, r.opts.fileType)
	if err != nil {
		return err
	}
	for _, fpath := range files {
		key := GetNamespace(strings.TrimRight(fpath, "."+strings.ToLower(r.opts.fileType)), r.opts.dir, r.opts.splitter)
		if _, ok := r.dicts[key]; !ok {
			r.dicts[key] = &I18nDict{}
		}
		r.dicts[key], err = r.parser.parse(fpath)
		if err != nil {
			return err
		}
	}

	return nil
}
