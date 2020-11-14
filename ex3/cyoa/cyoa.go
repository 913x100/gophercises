package cyoa

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type StoryArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

type handler struct {
	story    Story
	template *template.Template
}

type Story map[string]StoryArc

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

func LoadJson(file string) Story {
	jsonFile, err := os.Open(file)
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()

	jsonByte, _ := ioutil.ReadAll(jsonFile)

	story := make(Story)
	err = json.Unmarshal(jsonByte, &story)
	if err != nil {
		log.Fatalln(err)
	}

	return story
}

var defaultHandlerTmpl = `
<!DOCTYPE html>
<head>
    <meta charset="utf-8">
    <title>Choose Your Own Adventure</title>
</head>
<body>
    <h1>{{.Title}}</h1>
        {{range .Story}}
            <p>{{.}}</p>
        {{end}}
    <ul>
        {{range .Options}}
            <li> <a href="/{{.Arc}}">{{.Text}}</a></li>
        {{end}}
    </ul>
</body>`

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}

	if chapter, ok := h.story[path[1:]]; ok {
		err := h.template.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went Wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter Not Found.", http.StatusNotFound)

}

func NewHandler(story Story, tmpl *template.Template) http.Handler {
	if tmpl == nil {
		tmpl = tpl
	}
	return handler{
		story:    story,
		template: tmpl,
	}
}
