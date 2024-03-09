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
	// templates
)

func main() {
	g := gin.Default()

	tmpls := template.Must(template.ParseFS(templatesFS, "templates/*.gohtml"))
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

	log.Print("Listening...")
	http.ListenAndServe(":80", g)
}

func indexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index", nil)
}

func InputHandler(c *gin.Context) {
	field := controls.TextField{
		Value: "Lorem ipsum dolor sit amet",
	}
	var control controls.Control = &controls.FormControl[controls.TextField]{
		URL:          c.Request.URL.Path,
		Label:        "Input",
		Field:        field,
		ReadTemplate: "inputRead",
		EditTemplate: "inputEdit",
		SaveFunc: func(tx *controls.FormControl[controls.TextField]) error {
			log.Printf("Saving: %s", tx.Field.Value)
			return nil
		},
	}
	control.GinHandler(c)
}

func TextareaHandler(c *gin.Context) {
	field := controls.TextField{
		Value: "Lorem ipsum dolor sit amet, eum eligendi petentium temporibus te, et erant volumus erroribus duo. Id duo choro nullam philosophia.",
	}
	var control controls.Control = &controls.FormControl[controls.TextField]{
		URL:          c.Request.URL.Path,
		Label:        "Textarea",
		Field:        field,
		ReadTemplate: "textareaRead",
		EditTemplate: "textareaEdit",
		SaveFunc: func(tx *controls.FormControl[controls.TextField]) error {
			log.Printf("Saving: %s", tx.Field.Value)
			return nil
		},
	}
	control.GinHandler(c)
}
