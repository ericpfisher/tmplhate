package core

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

// flags
var TmplhateFuncs = template.FuncMap{
	"mul": func(a int, b int) int {
		return a * b
	},
}

// interface for tmplhate
type Tmplhater interface {
	LoadVars(io.Reader)
	ReadTemplate(io.Reader)
	Write(io.Writer) error
	ValidateTemplate() error
}

// impl of tmplhate interface
type Tmplhate struct {
	String string
	Tmpl   *template.Template
	Vars   map[string]any `yaml:"vars,omitempty,flow"`
}

// START Global helpers
// helper: check if input is stdin
func FromStdin() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

// helper: in = reader, out = bytes
func Read(r io.Reader) []byte {
	var b bytes.Buffer
	_, err := io.Copy(&b, r)
	if err != nil {
		log.Fatalf("Unable to read: %v", err)
	}
	return b.Bytes()
}

// END Global helpers

// START Tmplhate methods
// method to get template from input source
func (t *Tmplhate) GetReader(location string) io.Reader {
	var r io.Reader
	if location != "" {
		splitLocation := strings.SplitN(location, "://", 2)
		proto := splitLocation[0]
		switch proto {
		case "http":
			resp, err := http.Get(location)
			if err != nil {
				log.Fatalf("Unable to open http: %v", err)
			}
			r = resp.Body
		default:
			if strings.HasPrefix(location, "~/") {
				dirname, _ := os.UserHomeDir()
				location = filepath.Join(dirname, location[2:])
			} else if strings.HasPrefix(location, "../") {
				curdir, _ := os.Getwd()
				location = filepath.Join(curdir, location[3:])
			}
			f, err := os.Open(location)
			if err != nil {
				log.Fatalf("Unable to open file: %v", err)
			}
			r = f
		}
		return r
	}
	return nil
}

// method to load env vars into object
func (t *Tmplhate) LoadEnvVars() {
	env := make(map[string]any)
	for _, e := range os.Environ() {
		p := strings.SplitN(e, "=", 2)
		env[strings.ToUpper(p[0])] = p[1]
	}

	t.Vars["Env"] = env
}

// method to load variables into object
func (t *Tmplhate) LoadVars(r io.Reader) {
	t.Vars = make(map[string]any)
	if r != nil {
		if closer, ok := r.(io.ReadCloser); ok {
			defer closer.Close()
		}
		bytes := Read(r)
		caser := cases.Title(language.English)
		valuesMap := make(map[string]any)
		err := yaml.Unmarshal(bytes, &valuesMap)
		if err != nil {
			log.Fatalf("Unable to load vars: %v", err)
		}
		for k, v := range valuesMap {
			valuesMap[caser.String(k)] = v
		}
		t.Vars = valuesMap
	}
	t.LoadEnvVars()
}

// method to get template from reader and load
func (t *Tmplhate) LoadTemplate(r io.Reader) {
	if closer, ok := r.(io.ReadCloser); ok {
		defer closer.Close()
	}
	tbytes := Read(r)
	t.String = string(tbytes)
	t.ValidateTemplate()
	t.Tmpl = template.Must(template.New("tmplhate").Funcs(TmplhateFuncs).Parse(t.String))
}

// helper: validate template
func (t *Tmplhate) ValidateTemplate() {
	_, err := template.New("validate").Funcs(TmplhateFuncs).Parse(t.String)
	if err != nil {
		log.Fatalf("Invalid template: %v", err)
	}
}

// method to output rendered template to io.Writer
func (t *Tmplhate) WriteTemplate(w io.Writer) {
	err := t.Tmpl.Execute(w, t.Vars)
	if err != nil {
		log.Fatalf("Unable to write: %v", err)
	}
}

// END Tmplhate methods
