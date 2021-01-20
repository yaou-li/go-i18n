package i18n

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func ParserFactory(opts *I18nOpts) I18nParser {
	switch opts.fileType {
	case "json":
		return NewJsonParser(opts)
	default:
		return NewJsonParser(opts)
	}
}

func NewJsonParser(opts *I18nOpts) *JsonParser {
	return &JsonParser{opts}
}

type JsonParser struct {
	opts *I18nOpts
}

func (jp *JsonParser) parse(fpath string) (*I18nDict, error) {
	var (
		err  error
		dict I18nDict
	)
	fp, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	bytes, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal((bytes), &dict)
	if err != nil {
		return nil, err
	}
	return &dict, nil
}
