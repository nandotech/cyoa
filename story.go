package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTemplate))
}

var defaultHandlerTemplate = `<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-type" content="text/html;charset=UTF-8">
<style>
body {
font-family: -apple-system, BlinkMacSystemFont, 'Segoe WPC', 'Segoe UI', 'HelveticaNeue-Light', 'Ubuntu', 'Droid Sans', sans-serif;font-size: 14px;line-height: 1.6;
}
</style>
</head>
<body>
<h1>{{.Title}}</h1>
<hr />
{{range .Paragraphs}}
<p>{{.}}</p>
{{end}}
<ul>
{{range .Options}}
<li><a href="/{{.Arc}}">{{.Text}}</a></li>
{{end}}
</ul>
</body>
</html>`

// NewHandler function to handle requests
func NewHandler(s Story) http.Handler {
	return handler{s}
}

type handler struct {
	s Story
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}

	path = path[1:]

	if chapter, ok := h.s[path]; ok {
		err := tpl.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)

}

// JsonStory func returns Story from io.Reader
func JsonStory(r io.Reader) (Story, error) {
	d := json.NewDecoder(r)
	var story Story
	if err := d.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

// Story mapping to our Chapters
type Story map[string]Chapter

// Chapter type from json file
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

// Option type from json
type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}
