package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTemplate))
}

var defaultHandlerTemplate = `
<!DOCTYPE html>
    <html>
    <head>
        <meta http-equiv="Content-type" content="text/html;charset=UTF-8">
		<style>
		body { font-family: "Segoe WPC", "Segoe UI", "SFUIText-Light", "HelveticaNeue-Light", sans-serif, "Droid Sans Fallback"; font-size: 14px; padding: 0 12px; line-height: 22px; word-wrap: break-word; } body.scrollBeyondLastLine { margin-bottom: calc(100vh - 22px); } body.showEditorSelection .code-line { position: relative; } body.showEditorSelection .code-active-line:before, body.showEditorSelection .code-line:hover:before { content: ""; display: block; position: absolute; top: 0; left: -12px; height: 100%; } body.showEditorSelection li.code-active-line:before, body.showEditorSelection li.code-line:hover:before { left: -30px; } .vscode-light.showEditorSelection .code-active-line:before { border-left: 3px solid rgba(0, 0, 0, 0.15); } .vscode-light.showEditorSelection .code-line:hover:before { border-left: 3px solid rgba(0, 0, 0, 0.40); } .vscode-dark.showEditorSelection .code-active-line:before { border-left: 3px solid rgba(255, 255, 255, 0.4); } .vscode-dark.showEditorSelection .code-line:hover:before { border-left: 3px solid rgba(255, 255, 255, 0.60); } .vscode-high-contrast.showEditorSelection .code-active-line:before { border-left: 3px solid rgba(255, 160, 0, 0.7); } .vscode-high-contrast.showEditorSelection .code-line:hover:before { border-left: 3px solid rgba(255, 160, 0, 1); } img { max-width: 100%; max-height: 100%; } a { color: #4080D0; text-decoration: none; } a:focus, input:focus, select:focus, textarea:focus { outline: 1px solid -webkit-focus-ring-color; outline-offset: -1px; } hr { border: 0; height: 2px; border-bottom: 2px solid; } h1 { padding-bottom: 0.3em; line-height: 1.2; border-bottom-width: 1px; border-bottom-style: solid; } h1, h2, h3 { font-weight: normal; } h1 code, h2 code, h3 code, h4 code, h5 code, h6 code { font-size: inherit; line-height: auto; } a:hover { color: #4080D0; text-decoration: underline; } table { border-collapse: collapse; } table > thead > tr > th { text-align: left; border-bottom: 1px solid; } table > thead > tr > th, table > thead > tr > td, table > tbody > tr > th, table > tbody > tr > td { padding: 5px 10px; } table > tbody > tr + tr > td { border-top: 1px solid; } blockquote { margin: 0 7px 0 5px; padding: 0 16px 0 10px; border-left: 5px solid; } code { font-family: Menlo, Monaco, Consolas, "Droid Sans Mono", "Courier New", monospace, "Droid Sans Fallback"; font-size: 14px; line-height: 19px; } body.wordWrap pre { white-space: pre-wrap; } .mac code { font-size: 12px; line-height: 18px; } code > div { padding: 16px; border-radius: 3px; overflow: auto; } /** Theming */ .vscode-light { color: rgb(30, 30, 30); } .vscode-dark { color: #DDD; } .vscode-high-contrast { color: white; } .vscode-light code { color: #A31515; } .vscode-dark code { color: #D7BA7D; } .vscode-light code > div { background-color: rgba(220, 220, 220, 0.4); } .vscode-dark code > div { background-color: rgba(10, 10, 10, 0.4); } .vscode-high-contrast code > div { background-color: rgb(0, 0, 0); } .vscode-high-contrast h1 { border-color: rgb(0, 0, 0); } .vscode-light table > thead > tr > th { border-color: rgba(0, 0, 0, 0.69); } .vscode-dark table > thead > tr > th { border-color: rgba(255, 255, 255, 0.69); } .vscode-light h1, .vscode-light hr, .vscode-light table > tbody > tr + tr > td { border-color: rgba(0, 0, 0, 0.18); } .vscode-dark h1, .vscode-dark hr, .vscode-dark table > tbody > tr + tr > td { border-color: rgba(255, 255, 255, 0.18); } .vscode-light blockquote, .vscode-dark blockquote { background: rgba(127, 127, 127, 0.1); border-color: rgba(0, 122, 204, 0.5); } .vscode-high-contrast blockquote { background: transparent; border-color: #fff; }


/* Tomorrow Theme */ /* http://jmblog.github.com/color-themes-for-google-code-highlightjs */ /* Original theme - https://github.com/chriskempson/tomorrow-theme */ /* Tomorrow Comment */ .hljs-comment, .hljs-quote { color: #8e908c; } /* Tomorrow Red */ .hljs-variable, .hljs-template-variable, .hljs-tag, .hljs-name, .hljs-selector-id, .hljs-selector-class, .hljs-regexp, .hljs-deletion { color: #c82829; } /* Tomorrow Orange */ .hljs-number, .hljs-built_in, .hljs-builtin-name, .hljs-literal, .hljs-type, .hljs-params, .hljs-meta, .hljs-link { color: #f5871f; } /* Tomorrow Yellow */ .hljs-attribute { color: #eab700; } /* Tomorrow Green */ .hljs-string, .hljs-symbol, .hljs-bullet, .hljs-addition { color: #718c00; } /* Tomorrow Blue */ .hljs-title, .hljs-section { color: #4271ae; } /* Tomorrow Purple */ .hljs-keyword, .hljs-selector-tag { color: #8959a8; } .hljs { display: block; overflow-x: auto; color: #4d4d4c; padding: 0.5em; } .hljs-emphasis { font-style: italic; } .hljs-strong { font-weight: bold; }


ul.contains-task-list { padding-left: 0; } ul ul.contains-task-list { padding-left: 40px; } .task-list-item { list-style-type: none; } .task-list-item-checkbox { vertical-align: middle; }

        
            body {
                font-family: -apple-system, BlinkMacSystemFont, 'Segoe WPC', 'Segoe UI', 'HelveticaNeue-Light', 'Ubuntu', 'Droid Sans', sans-serif;
                font-size: 14px;
                line-height: 1.6;
            }
        </style>
    </head>
    <body>
        <h1>{{.Title}}</h1>
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
	err := tpl.Execute(w, h.s["intro"])
	if err != nil {
		panic(err)
	}
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
