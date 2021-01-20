package extract

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/yaou-li/go-i18n"
)

type I18nExtractor interface {
	Extract(sourced string, clean bool) error
}

type extractor struct {
	opts   *i18n.I18nOpts
	reader i18n.I18nReader
	writer i18n.I18nWriter
	log    *logrus.Logger
}

func NewExtractor(opts *i18n.I18nOpts) I18nExtractor {
	dicts := make(map[string]*i18n.I18nDict)
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		DisableColors: true,
	}
	return &extractor{
		opts:   opts,
		writer: i18n.NewWriter(opts, dicts),
		reader: i18n.NewReader(opts, dicts),
		log:    log,
	}
}

/**
* read all existing trans files and
* recursively extract all i18n Trans/Transf calls from go files
**/
func (ex *extractor) Extract(sourced string, clean bool) error {
	if clean {
		os.RemoveAll(ex.opts.GetDir())
	} else {
		if err := ex.reader.ReadAllFile(); err != nil {
			ex.log.Error(err)
		}
	}

	var s []string
	files, err := i18n.ReadAllPath(sourced, s, "go")
	if err != nil {
		ex.log.Error(err)
		return err
	}
	fset := token.NewFileSet()
	for _, fname := range files {
		node, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
		if err != nil {
			ex.log.Errorf("Error when parsing source file: %v, error: %v", fname, err)
		}
		ast.Inspect(node, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.File:
				fmt.Sprint(x.Name)
			case *ast.CallExpr:
				if callexp, ok := n.(*ast.CallExpr); !ok {
					return false
				} else {
					if fun, ok := callexp.Fun.(*ast.SelectorExpr); ok {
						if pack, ok := fun.X.(*ast.Ident); (fun.Sel.Name == "Trans" || fun.Sel.Name == "Transf") && ok {
							fpath := i18n.GetNamespace(path.Dir(fname), sourced, ex.opts.GetSplitter())
							if len(callexp.Args) == 0 {
								ex.log.Error("Missing translation data")
								return false
							}
							key, ok := callexp.Args[0].(*ast.BasicLit)
							if !ok {
								ex.log.Error("Unable to get key value")
								return false
							}
							if ex.opts.IsNamespaced() {
								ex.writer.Append(fpath, strings.Trim(key.Value, "\""))
							} else {
								ex.writer.Append("index", strings.Trim(key.Value, "\""))
							}
							fmt.Println("------Package name------")
							fmt.Println(pack.Name, fun.Sel.Name, fname, fpath, strings.Trim(key.Value, "\""))
							fmt.Println("------Package name------")
						}
					}
				}
			}
			return true
		})
	}

	if err := ex.writer.Flush(); err != nil {
		ex.log.Errorf("Failed to flush to i18n files, error: %v", err)
	}
	return nil
}
