package i18n

import "fmt"

type dict map[string]string

type dictWithNamespace map[string]dict

type I18nDict struct {
	Lang      string `json:"language"`
	Namespace string `json:"namespace,omitempty"`
	Dict      dict   `json:"dict"`
}

func (d *I18nDict) Merge(ndict *I18nDict) error {
	if d.Lang != ndict.Lang {
		return fmt.Errorf("Failed to merge two dict, language mismatching, source: %v, target: %v", ndict.Lang, d.Lang)
	}
	for key, val := range ndict.Dict {
		if _, ok := d.Dict[key]; !ok {
			d.Dict[key] = val
		}
	}
	return nil
}

func (d *I18nDict) Overwrite(ndict *I18nDict) error {
	if d.Lang != ndict.Lang {
		return fmt.Errorf("Failed to merge two dict, language mismatching, source: %v, target: %v", ndict.Lang, d.Lang)
	}
	for key, val := range ndict.Dict {
		d.Dict[key] = val
	}
	return nil
}

func (d *I18nDict) Clone() *I18nDict {
	nd := &I18nDict{
		Lang:      d.Lang,
		Namespace: d.Namespace,
		Dict:      make(dict),
	}
	for k, v := range d.Dict {
		nd.Dict[k] = v
	}
	return nd
}
