package i18n

import (
	"io/ioutil"
	"path"
	"strings"
)

/**
* recursively read all file path that matches the file type
 */
func ReadAllPath(dir string, s []string, fileType string) ([]string, error) {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return s, err
	}
	for _, fi := range fis {
		name := fi.Name()
		if fi.IsDir() {
			fullir := path.Join(dir, name)
			s, err = ReadAllPath(fullir, s, fileType)
			if err != nil {
				return s, err
			}
		} else {
			if !strings.Contains(name, ".") {
				continue
			}
			ts := strings.Split(name, ".")
			ft := ts[len(ts)-1]
			if fileType != ft {
				// fmt.Printf("Unsupported file type: %v, loader file type: %v\n", name, fileType)
				continue
			}
			fullName := path.Join(dir, name)
			s = append(s, fullName)
		}
	}
	return s, nil
}

func GetNamespace(fulld string, parentd string, splitter string) string {
	// normalize the path format
	fulld = cleanPath(fulld)
	parentd = cleanPath(parentd)
	// if does not contain parentd, just return fulld
	if !strings.Contains(fulld, parentd) {
		return convert(fulld, splitter)
	}
	dirs := strings.SplitN(fulld, parentd, 2)
	return convert(path.Join(dirs[1:]...), splitter)
}

func cleanPath(p string) string {
	return strings.ReplaceAll(path.Clean(p), "\\", "/")
}

func convert(p string, splitter string) string {
	p = strings.Trim(p, "/")
	return strings.Join(strings.Split(p, "/"), splitter)
}
