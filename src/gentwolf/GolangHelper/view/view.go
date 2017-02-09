package view

import (
	"html/template"
	"io"
	"strings"
)

var (
	tpls map[string]*template.Template

	tplPath string
	tplExt  string
)

//设置模板根目录，及后缀。
func SetConfig(path, ext string) {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	tplPath = path
	tplExt = ext

	tpls = make(map[string]*template.Template)
}

//输出模板
func Render(wr io.Writer, name string, data interface{}, files ...string) error {
	tpl, bl := tpls[name]
	if !bl {
		var err error

		for k, v := range files {
			files[k] = tplPath + v + tplExt
		}

		tpl, err = template.ParseFiles(files...)
		if err != nil {
			return err
		}
		tpls[name] = tpl
	}

	return tpl.Execute(wr, data)
}

//输出含有layout的模板，
func RenderWithLayout(wr io.Writer, name string, data interface{}, layout string, files ...string) error {
	files = append(files, layout)
	return Render(wr, name, data, files...)
}
