package Handlers

import (
	"bytes"
	"github.com/3stadt/presla/src/PreslaTemplates"
	"github.com/labstack/echo"
	"github.com/spf13/afero"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

func (conf *Conf) Md(c echo.Context) error {
	file := c.Param("file")

	if file == "info.md" {
		return conf.showInfo(c)
	}

	file = conf.MarkdownPath + "/" + file

	tpl, err := afero.ReadFile(conf.Fs, file)
	if err != nil {
		c.NoContent(http.StatusNotFound)
		return err
	}
	// Rendering is needed so Code isn't commented automatically
	return render(c, tpl, nil)
}

func (conf *Conf) showInfo(c echo.Context) error {
	var presentations []string

	files := make([]string, 1)
	files, err := afero.Glob(conf.Fs, conf.MarkdownPath+"/*.md")
	if err != nil {
		files[0] = "Error loading presentations: " + err.Error()
	}
	for _, file := range files {
		presentations = append(presentations, strings.TrimSuffix(filepath.Base(file), ".md"))
	}
	tmpDir, err := afero.TempDir(conf.Fs, "", "presla")
	if err != nil {
		tmpDir = "/tmp"
	}
	data := map[string]interface{}{
		"Presentations": presentations,
		"ConfigFile":    conf.ConfigFile,
		"TempDir":       tmpDir,
	}

	tpl, err := Asset("templates/info.md")
	if err != nil {
		return err
	}
	return render(c, tpl, data)
}

func render(c echo.Context, tpl []byte, data map[string]interface{}) error {
	parsedTemplate, err := template.New("default").Parse(string(tpl))
	if err != nil {
		return err
	}
	t := &PreslaTemplates.DefaultTemplate{
		Template: parsedTemplate,
	}

	buf := new(bytes.Buffer)
	err = t.Render(buf, "default", data, c)
	if err != nil {
		return err
	}
	return c.Blob(http.StatusOK, "text/markdown; charset=utf-8", buf.Bytes())
}
