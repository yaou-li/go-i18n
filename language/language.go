package language

import "strings"

type I18nLang uint16

const (
	Afrikaans = iota + 1
	Amharic
	Arabic
	ModernStandardArabic
	Azerbaijani
	Bulgarian
	Bengali
	Catalan
	Czech
	Danish
	German
	Greek
	English
	AmericanEnglish
	BritishEnglish
	Spanish
	EuropeanSpanish
	LatinAmericanSpanish
	Estonian
	Persian
	Finnish
	Filipino
	French
	CanadianFrench
	Gujarati
	Hebrew
	Hindi
	Croatian
	Hungarian
	Armenian
	Indonesian
	Icelandic
	Italian
	Japanese
	Georgian
	Kazakh
	Khmer
	Kannada
	Korean
	Kirghiz
	Lao
	Lithuanian
	Latvian
	Macedonian
	Malayalam
	Mongolian
	Marathi
	Malay
	Burmese
	Nepali
	Dutch
	Norwegian
	Punjabi
	Polish
	Portuguese
	BrazilianPortuguese
	EuropeanPortuguese
	Romanian
	Russian
	Sinhala
	Slovak
	Slovenian
	Albanian
	Serbian
	SerbianLatin
	Swedish
	Swahili
	Tamil
	Telugu
	Thai
	Turkish
	Ukrainian
	Urdu
	Uzbek
	Vietnamese
	Chinese
	SimplifiedChinese
	TraditionalChinese
	Zulu
)

func (lang I18nLang) Shortcut() string {
	switch lang {
	case English:
		return "en"
	case Chinese:
		return "zh"
	case Korean:
		return "ko"
	case Russian:
		return "ru"
	case Japanese:
		return "ja"
	default:
		return "unsupported language"
	}
}

var validLangMap = map[string]I18nLang{
	"en": English,
	"zh": Chinese,
	"ko": Korean,
	"ru": Russian,
	"ja": Japanese,
}

var validLang = []I18nLang{
	English,
	Chinese,
	Korean,
	Russian,
	Japanese,
}

func LangMap() map[string]I18nLang {
	return validLangMap
}

func GetLang(shortcut string) I18nLang {
	shortcut = strings.ToLower(shortcut)
	if lang, ok := validLangMap[shortcut]; !ok {
		return 0
	} else {
		return lang
	}
}

func Langs() []I18nLang {
	return validLang
}

func IsSupported(shortcut string) bool {
	shortcut = strings.ToLower(shortcut)
	if _, ok := validLangMap[shortcut]; ok {
		return true
	}
	return false
}
