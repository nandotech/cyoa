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
<title>Choose Your Own Adventure!</title>
<body>
<section class="page">
  <h1>{{.Title}}</h1>
  {{range .Paragraphs}}
	<p>{{.}}</p>
  {{end}}
  {{if .Options}}
	<ul>
	{{range .Options}}
	  <li><a href="/{{.Arc}}">{{.Text}}</a></li>
	{{end}}
	</ul>
  {{else}}
	<h3>The End</h3>
  {{end}}
</section>
<style>
  body {
	font-family: helvetica, arial;
  }
  h1 {
	text-align:center;
	position:relative;
  }
  .page {
	width: 80%;
	max-width: 500px;
	margin: auto;
	margin-top: 40px;
	margin-bottom: 40px;
	padding: 80px;
	background: #FFFCF6;
	border: 1px solid #eee;
	box-shadow: 0 10px 6px -6px #777;
  }
  ul {
	border-top: 1px dotted #ccc;
	padding: 10px 0 0 0;
	-webkit-padding-start: 0;
  }
  li {
	padding-top: 10px;
  }
  a,
  a:visited {
	text-decoration: none;
	color: #6295b5;
  }
  a:active,
  a:hover {
	color: #7792a2;
  }
  p {
	text-indent: 1em;
  }
</style>
</body>
</html>`

func WithTemplate(t *template.Template) HandlerOption {
	return func(h *handler) {
		h.t = t
	}
}

// func WithDatabase(username, password string) HandlerOption {

// }
type HandlerOption func(h *handler)

// NewHandler function to handle requests
func NewHandler(s Story, opts ...HandlerOption) http.Handler {
	h := handler{s, tpl}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

type handler struct {
	s Story
	t *template.Template
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}

	path = path[1:]

	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)
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
