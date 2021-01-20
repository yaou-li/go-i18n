package i18n

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

/**
* extract writer
* flush() write all the cached data into json file
**/
type I18nWriter interface {
	Append(namespace string, key string) error
	Flush() error
	WriteJSON(namespace string, dict *I18nDict) error
}

func NewWriter(opts *I18nOpts, dicts map[string]*I18nDict) I18nWriter {
	return &writer{
		opts:   opts,
		odicts: dicts,
		ndicts: make(map[string]*I18nDict),
	}
}

type writer struct {
	sync.Mutex
	opts   *I18nOpts
	odicts map[string]*I18nDict
	ndicts map[string]*I18nDict
}

func (w *writer) Append(namespace string, key string) error {
	if namespace == "" {
		if w.opts.enableNamespace {
			return nil
		}
		namespace = "index"
	}
	w.Lock()
	defer w.Unlock()
	if _, ok := w.ndicts[namespace]; !ok {
		w.ndicts[namespace] = &I18nDict{
			Dict: make(dict),
		}
		w.ndicts[namespace].Dict[key] = ""
	} else if _, ok := w.ndicts[namespace].Dict[key]; !ok {
		w.ndicts[namespace].Dict[key] = ""
	}
	return nil
}

func (w *writer) Flush() error {
	w.Lock()
	defer w.Unlock()
	for _, lang := range w.opts.langs {
		for namespace, dict := range w.ndicts {
			ndict := dict.Clone()
			namespace = strings.Join([]string{lang.Shortcut(), namespace}, ".")
			ndict.Lang = lang.Shortcut()
			if _, ok := w.odicts[namespace]; ok {
				ndict.Overwrite(w.odicts[namespace])
			}
			if w.opts.enableNamespace {
				ndict.Namespace = namespace
			}
			switch strings.ToUpper(w.opts.fileType) {
			case "JSON":
				if err := w.WriteJSON(namespace, ndict); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (w *writer) WriteJSON(namespace string, dict *I18nDict) error {
	p := path.Join(strings.Split(namespace, w.opts.splitter)...)
	p = p + ".json"
	dir := path.Dir(p)
	// make sure the directory already exists
	if err := os.MkdirAll(path.Join(w.opts.dir, dir), os.ModePerm); err != nil {
		return err
	}
	if data, err := json.MarshalIndent(dict, "", "    "); err != nil {
		return err
	} else {
		return ioutil.WriteFile(path.Join(w.opts.dir, p), data, os.ModePerm)
	}
}
