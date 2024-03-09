package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"htmx.samples.dev/controls"

	"github.com/gin-gonic/gin"
)

var (
	//go:embed static/*
	staticFS embed.FS
	//go:embed templates/*
	templatesFS embed.FS
)

func main() {
	g := gin.Default()

	funcMap := template.FuncMap{
		"stringInSlice": stringInSlice,
	}

	tmpls, err := template.New("").Funcs(funcMap).ParseFS(templatesFS, "templates/*.gohtml")
	if err != nil {
		log.Fatalf("Unable to parse templates: %v", err)
	}
	g.SetHTMLTemplate(tmpls)

	static, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatalf("Unable to load static files: %v", err)
	}
	g.StaticFS("/static", http.FS(static))

	g.GET("/", indexHandler)
	g.GET("/input", InputHandler)
	g.POST("/input", InputHandler)
	g.GET("/textarea", TextareaHandler)
	g.POST("/textarea", TextareaHandler)
	g.GET("/choice", ChoiceHandler)
	g.POST("/choice", ChoiceHandler)
	g.GET("/multi-choice", MultiChoiceHandler)
	g.POST("/multi-choice", MultiChoiceHandler)

	log.Print("Listening...")
	http.ListenAndServe(":80", g)
}

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index", nil)
}

func InputHandler(c *gin.Context) {
	control := &controls.FormControl[controls.TextField]{
		URL:   c.Request.URL.Path,
		Label: "Input",
		Field: controls.TextField{
			Value: "Lorem ipsum dolor sit amet",
		},
		ReadTemplate: "inputRead",
		EditTemplate: "inputEdit",
		SaveFunc: func(tx *controls.FormControl[controls.TextField]) error {
			log.Printf("Saving: %s", tx.Field.Value)
			return nil
		},
	}
	if c.Request.Method == http.MethodPost {
		control.Field.Value = ""
	}
	control.GinHandler(c)
}

func TextareaHandler(c *gin.Context) {
	control := &controls.FormControl[controls.TextField]{
		URL:   c.Request.URL.Path,
		Label: "Textarea",
		Field: controls.TextField{
			Value: "Lorem ipsum dolor sit amet, eum eligendi petentium temporibus te, et erant volumus erroribus duo. Id duo choro nullam philosophia.",
		},
		ReadTemplate: "textareaRead",
		EditTemplate: "textareaEdit",
		SaveFunc: func(tx *controls.FormControl[controls.TextField]) error {
			log.Printf("Saving: %s", tx.Field.Value)
			return nil
		},
	}
	if c.Request.Method == http.MethodPost {
		control.Field.Value = ""
	}
	control.GinHandler(c)
}

func ChoiceHandler(c *gin.Context) {
	control := &controls.FormControl[controls.ChoiceField]{
		URL:   c.Request.URL.Path,
		Label: "Choice",
		Field: controls.ChoiceField{
			Value:   "Option 1",
			Choices: []string{"Option 1", "Option 2", "Option 3", "Option 4", "Option 5"},
		},
		ReadTemplate: "choiceRead",
		EditTemplate: "choiceEdit",
		SaveFunc: func(tx *controls.FormControl[controls.ChoiceField]) error {
			log.Printf("Saving: %s", tx.Field.Value)
			return nil
		},
	}
	if c.Request.Method == http.MethodPost {
		control.Field.Value = ""
	}
	control.GinHandler(c)
}

func MultiChoiceHandler(c *gin.Context) {
	control := &controls.FormControl[controls.MultiChoiceField]{
		URL:   c.Request.URL.Path,
		Label: "Multi Choice",
		Field: controls.MultiChoiceField{
			Values:  []string{"Option 1", "Option 2"},
			Choices: []string{"Option 1", "Option 2", "Option 3", "Option 4", "Option 5"},
		},
		ReadTemplate: "multiChoiceRead",
		EditTemplate: "multiChoiceEdit",
		SaveFunc: func(tx *controls.FormControl[controls.MultiChoiceField]) error {
			log.Printf("Saving: %s", tx.Field.Values)
			return nil
		},
	}
	if c.Request.Method == http.MethodPost {
		control.Field.Values = []string{}
	}
	control.GinHandler(c)
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
