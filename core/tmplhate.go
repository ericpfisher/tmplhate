package core

import (
	"bytes"
	"fmt"
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
	"lower": func(s string) string {
		return cases.Lower(language.English).String(s)
	},
	"upper": func(s string) string {
		return cases.Upper(language.English).String(s)
	},
	"mul": func(a int, b int) int {
		return a * b
	},
}

// interface for tmplhate
type Tmplhater interface {
	LoadVars(io.Reader, bool, string)
	ReadTemplate(io.Reader)
	Write(io.Writer) error
	ValidateTemplate() error
}

// impl of tmplhate interface
type Tmplhate struct {
	Caser         cases.Caser
	Language      language.Tag
	NormalizeVars bool
	String        string
	Tmpl          *template.Template
	Vars          map[string]any `yaml:"vars,omitempty,flow"`
	VarsCase      string
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

// method to set a cases.Caser
func (t *Tmplhate) LoadCaser() {
	switch t.VarsCase {
	case "lower":
		t.Caser = cases.Lower(t.Language)
	case "upper":
		t.Caser = cases.Upper(t.Language)
	case "title":
		t.Caser = cases.Title(t.Language)
	default:
		log.Fatalf("Caser type '%s' is invalid.", t.VarsCase)
	}
}

// method to load env vars into object
func (t *Tmplhate) LoadEnvVars() {
	var es string
	env := make(map[string]any)
	for _, e := range os.Environ() {
		var ks string
		p := strings.SplitN(e, "=", 2)
		es = "env" // default Caser is 'lower'
		ks = p[0]
		env[ks] = p[1]
		if t.NormalizeVars {
			nks := t.Caser.String(p[0])
			env[nks] = p[1]
		}
	}

	if t.NormalizeVars {
		nes := t.Caser.String("env")
		t.Vars[nes] = env
	}

	t.Vars[es] = env
}

// method to load variables into object
func (t *Tmplhate) LoadVars(r io.Reader) {
	if r != nil {
		if closer, ok := r.(io.ReadCloser); ok {
			defer closer.Close()
		}
		bytes := Read(r)
		valuesMap := make(map[string]any)
		err := yaml.Unmarshal(bytes, &valuesMap)
		if err != nil {
			log.Fatalf("Unable to load vars: %v", err)
		}
		for k, v := range valuesMap {
			var ks string
			if t.NormalizeVars {
				ks = t.Caser.String(k)
			} else {
				ks = k
			}
			valuesMap[ks] = v
		}
		t.Vars = valuesMap
	}
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

func (t *Tmplhate) Init(tmplLocation string, varsLocation string, dontNormalize bool, varsCase string) {
	// create map for vars
	t.Vars = make(map[string]any)
	// set attrs of object:
	//   Normalize is opposite of dontNormalize value
	t.NormalizeVars, t.VarsCase = !dontNormalize, varsCase
	// read from location if one is provided
	// else, try to read from stdin
	if tmplLocation != "" {
		t.LoadTemplate(t.GetReader(tmplLocation))
	} else {
		fileInfo, _ := os.Stdin.Stat()
		if fileInfo.Mode()&os.ModeCharDevice == 0 {
			t.LoadTemplate(os.Stdin)
		} else {
			fmt.Println("Use 'tmplhate --help' for usage.")
			os.Exit(1)
		}
	}
	t.LoadCaser()
	t.LoadVars(t.GetReader(varsLocation))
	t.LoadEnvVars()
}

// END Tmplhate methods
