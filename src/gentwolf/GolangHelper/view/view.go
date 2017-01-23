// 模板类
// 扩展了html/template，支持layout功能，增加缓存功能，定时检测模板根目录中的状态文件，判断是否需要清除缓存。

package view

import (
	"html/template"
	"io"
	"os"
	"strings"
	"time"
)

var (
	tpls map[string]*template.Template

	tplPath string
	tplExt  string
)

//在模板根目录中检测是否有“status.txt”文件，如果有，则清除缓存。
func checkStatus() {
	for _ = range time.Tick(10 * time.Second) {
		status := tplPath + "status.txt"
		if _, err := os.Stat(status); err == nil {
			for k, _ := range tpls {
				delete(tpls, k)
			}
			os.Remove(status)
		}
	}
}

//设置模板根目录，及后缀。
func SetConfig(path, ext string) {
	if !strings.HasSuffix(tplPath, "/") {
		tplPath += "/"
	}

	tplPath = path
	tplExt = ext

	tpls = make(map[string]*template.Template)

	go checkStatus()
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
