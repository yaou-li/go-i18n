package i18n

import (
	"fmt"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/yaou-li/go-i18n/language"
)

type loader struct {
	sync.Mutex
	opts               *I18nOpts
	log                *logrus.Logger
	parser             I18nParser
	dicts              map[language.I18nLang]dict
	dictsWithNamespace map[language.I18nLang]dictWithNamespace
}

func Newloader(opts *I18nOpts, log *logrus.Logger) *loader {
	return &loader{
		opts:               opts,
		log:                log,
		parser:             ParserFactory(opts),
		dicts:              make(map[language.I18nLang]dict),
		dictsWithNamespace: make(map[language.I18nLang]dictWithNamespace),
	}
}

func (l *loader) load() error {
	var s []string
	// read all files
	files, err := ReadAllPath(l.opts.dir, s, l.opts.fileType)
	if err != nil {
		l.log.Error(err)
		return err
	}
	// loop through each file and save in dicts or dictsWithNamespace accordingly
	l.Lock()
	defer l.Unlock()
	for _, fpath := range files {
		data, err := l.parser.parse(fpath)
		if err != nil {
			l.log.Errorf("Fail to parser file: %v, error: %v", fpath, err)
			return err
		}
		// check if lang is valid
		if !l.opts.IsEnabled(data.Lang) {
			l.log.Errorf("Unsupported language: %v", data.Lang)
			continue
		}
		// store in dicts with namespace if enabled, store in dicts otherwise
		if l.opts.enableNamespace {
			l.mergeWithNameSpace(fpath, data)
		} else {
			l.merge(data)
		}
	}
	return nil
}

func (l *loader) merge(data *I18nDict) {
	lang := language.GetLang(data.Lang)
	if _, ok := l.dicts[lang]; !ok {
		l.dicts[lang] = make(dict)
	}
	for k, v := range data.Dict {
		l.dicts[lang][k] = v
	}
}

func (l *loader) mergeWithNameSpace(fpath string, data *I18nDict) {
	lang := language.GetLang(data.Lang)
	namespace := GetNamespace(strings.TrimRight(fpath, "."+l.opts.fileType), l.opts.dir, l.opts.splitter)
	if namespace != data.Namespace {
		l.log.Error("Failed to load into namespace, namespace unmatched: %v vs %v", namespace, data.Namespace)
		// if namespace is not matched, fallback to general dict
		l.merge(data)
		return
	}
	if _, ok := l.dictsWithNamespace[lang]; !ok {
		l.dictsWithNamespace[lang] = make(dictWithNamespace)
	}
	l.dictsWithNamespace[lang][namespace] = make(dict)
	for k, v := range data.Dict {
		l.dictsWithNamespace[lang][namespace][k] = v
	}
}

func (l *loader) reload() error {
	l.Lock()
	defer l.Unlock()
	l.reset()
	return l.load()
}

func (l *loader) reset() {
	for lang := range l.dicts {
		delete(l.dicts, lang)
	}
	for lang := range l.dictsWithNamespace {
		delete(l.dictsWithNamespace, lang)
	}
}

func (l *loader) get(key string) string {
	if dict, ok := l.dicts[l.opts.target]; !ok {
		l.log.Errorf("Missing translation for lang: %v", l.opts.target.Shortcut())
		return key
	} else {
		if val, ok := dict[key]; !ok || val == "" {
			l.log.Errorf("Missing translation: %v", key)
			return key
		} else {
			return val
		}
	}
}

func (l *loader) getWithNamespace(key string, namespace string) string {
	if dicts, ok := l.dictsWithNamespace[l.opts.target]; !ok {
		l.log.Errorf("Missing translation in namespace mode, lang: %v", l.opts.target.Shortcut())
		// fall back with none namespaced dict
		return l.get(key)
	} else {
		if dict, ok := dicts[namespace]; !ok {
			l.log.Errorf("Missing translation in namespace %v, for %v", namespace, key)
			// fall back with none namespaced dict
			return l.get(key)
		} else {
			if val, ok := dict[key]; !ok || val == "" {
				l.log.Errorf("Missing translation in namespace %v, for %v", namespace, key)
				// fall back with none namespaced dict
				return l.get(key)
			} else {
				return val
			}
		}
	}
}

func (l *loader) getDict(lang language.I18nLang) (dict, error) {
	if dict, ok := l.dicts[lang]; !ok {
		return nil, fmt.Errorf("Unloaded dict with lang :%v", lang.Shortcut())
	} else {
		return dict, nil
	}

}
