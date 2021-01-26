package i18n

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/yaou-li/go-i18n/language"
)

var once sync.Once
var i18nSingleton *i18n
var i18nRuntimeDir string

type I18nOpts struct {
	target          language.I18nLang
	src             language.I18nLang
	langs           []language.I18nLang
	splitter        string
	dir             string
	fileType        string
	enableNamespace bool
}

func NewI18nOpts() *I18nOpts {
	// only target lang and directory needs to be set specifically
	defaultOpts := &I18nOpts{
		src:             language.English,
		splitter:        ".",
		dir:             "./i18n",
		fileType:        "json",
		enableNamespace: false,
	}
	defaultOpts.SetEnableLangs("en,ko,zh,ru,ja")
	return defaultOpts
}

func (opts *I18nOpts) SetEnableLangs(shortcuts string) {
	var langs []string
	if strings.Contains(shortcuts, ",") {
		langs = strings.Split(shortcuts, ",")
	} else { // 如果key不包含点, 则就是直接一级调用
		langs = []string{shortcuts}
	}
	for _, lang := range langs {
		if language.IsSupported(lang) {
			opts.langs = append(opts.langs, language.GetLang(lang))
		}
	}
}

func (opts *I18nOpts) SetTargetLang(shortcut string) {
	if !language.IsSupported(shortcut) {
		panic(fmt.Sprintf("target language: %v is not supported", shortcut))
	}
	opts.target = language.GetLang(shortcut)
}

func (opts *I18nOpts) SetSrcLang(shortcut string) {
	if !language.IsSupported(shortcut) {
		panic(fmt.Sprintf("target language: %v is not supported", shortcut))
	}
	opts.src = language.GetLang(shortcut)
}

func (opts *I18nOpts) SetSplitter(splitter string) {
	opts.splitter = splitter
}

func (opts *I18nOpts) SetLanguageDir(dir string) {
	opts.dir = dir
}

func (opts *I18nOpts) SetFileType(fileType string) {
	opts.fileType = fileType
}

func (opts *I18nOpts) SetEnableNamespace(enable bool) {
	opts.enableNamespace = enable
}

func (opts *I18nOpts) IsEnabled(shortcut string) bool {
	if !language.IsSupported(shortcut) {
		return false
	}
	lang := language.GetLang(shortcut)
	for _, l := range opts.langs {
		if l == lang {
			return true
		}
	}
	return false
}

func (opts *I18nOpts) IsNamespaced() bool {
	return opts.enableNamespace
}

func (opts *I18nOpts) GetDir() string {
	return opts.dir
}

func (opts *I18nOpts) GetSplitter() string {
	return opts.splitter
}

type i18n struct {
	opts   *I18nOpts
	log    *logrus.Logger
	loader *loader
}

func Init(opts *I18nOpts, log *logrus.Logger) {
	once.Do(func() {
		i18nSingleton = &i18n{
			opts:   opts,
			log:    log,
			loader: Newloader(opts, log),
		}
		if dir, err := os.Getwd(); err != nil {
			log.Error("Failed to get runtime folder.")
		} else {
			i18nRuntimeDir = dir
		}
		if err := i18nSingleton.loader.load(); err != nil {
			log.Errorf("Failed to load trans data, error: %v", err)
		}
	})
}

func (i18n *i18n) GetLang() language.I18nLang {
	return i18n.opts.target
}

func (i18n *i18n) GetShortcut() string {
	return i18n.opts.target.Shortcut()
}

func (i18n *i18n) UpdateLang(shortcut string) {
	if i18n.opts.IsEnabled(shortcut) {
		i18n.opts.SetTargetLang(shortcut)
	}
}

func (i18n *i18n) trans(key string) string {
	if i18n.opts.enableNamespace {
		if _, fpath, _, ok := runtime.Caller(2); !ok {
			i18n.log.Errorf("Failed to get caller of trans function, key: %v", key)
			return i18n.loader.get(key)
		} else {
			namespace := i18n.opts.target.Shortcut() + "." + GetNamespace(filepath.Dir(fpath), i18nRuntimeDir, i18n.opts.splitter)
			// fmt.Println("-----------")
			// fmt.Println(filepath.Dir(fpath), i18nRuntimeDir, namespace)
			// fmt.Println("-----------")
			return i18n.loader.getWithNamespace(key, namespace)
		}
	} else {
		return i18n.loader.get(key)
	}
}

func GetDicts() map[language.I18nLang]dict {
	if i18nSingleton == nil {
		return make(map[language.I18nLang]dict)
	} else {
		return i18nSingleton.loader.dicts
	}
}

func Trans(key string) string {
	if i18nSingleton == nil {
		panic("i18n is not initialized.")
	} else {
		return i18nSingleton.trans(key)
	}
}

func Transf(key string, a ...interface{}) string {
	if i18nSingleton == nil {
		panic("i18n is not initialized.")
	} else {
		return fmt.Sprintf(i18nSingleton.trans(key), a...)
	}
}

func UpdateLang(shortcut string) {
	i18nSingleton.UpdateLang(shortcut)
}

func GetLang() string {
	return i18nSingleton.opts.target.Shortcut()
}
