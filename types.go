package i18n

import "github.com/yaou-li/go-i18n/language"

type I18nOptsInterface interface {
	SetEnableLangs(shortcuts string)
	SetTargetLang(shortcut string)
	SetSrcLang(shortcut string)
	SetSplitter(splitter string)
	SetLanguageDir(dir string)
	SetFileType(fileType string)
	SetEnableNamespace(enable bool)
	IsEnabled(shortcut string) bool
}

type I18nInterface interface {
	GetLang() language.I18nLang
	GetShortcut() string
	UpdateLang(shortcut string)
	Trans(key string) string
	Transf(key string, a ...interface{}) string
}

type I18nLoader interface {
	load() error
	reload() error
	reset()
	get(key string) string
	merge(data *I18nDict)
	mergeWithNameSpace(data *I18nDict)
	getWithNamespace(key string, namespace string) string
	ReadAllPath(dir string, s []string) ([]string, error)
	getDict(lang language.I18nLang) (dict, error)
}

type I18nParser interface {
	parse(fpath string) (*I18nDict, error)
}
