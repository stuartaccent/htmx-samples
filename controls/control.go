// Package controls provides an interface and a struct for form controls.
package controls

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Control is an interface that defines the methods that a form control should implement.
type Control interface {
	GetURL() string
	GetLabel() string
	GetError() string
	GetReadTemplate() string
	GetEditTemplate() string
	GinHandler(c *gin.Context)
	RenderRead(c *gin.Context)
	RenderEdit(c *gin.Context)
	Save() error
	SetError(err error)
}

// FormControl is a struct that implements the Control interface.
// It represents a form control with a URL, label, field, error message, read and edit templates, and a save function.
type FormControl[F any] struct {
	URL          string
	Label        string
	Field        F
	Error        string
	ReadTemplate string
	EditTemplate string
	SaveFunc     func(*FormControl[F]) error
}

// GetURL returns the URL of the form control.
func (c *FormControl[F]) GetURL() string {
	return c.URL
}

// GetLabel returns the label of the form control.
func (c *FormControl[F]) GetLabel() string {
	return c.Label
}

// GetError returns the error message of the form control.
func (c *FormControl[F]) GetError() string {
	return c.Error
}

// GetReadTemplate returns the read template of the form control.
func (c *FormControl[F]) GetReadTemplate() string {
	return c.ReadTemplate
}

// GetEditTemplate returns the edit template of the form control.
func (c *FormControl[F]) GetEditTemplate() string {
	return c.EditTemplate
}

// GinHandler handles the HTTP requests for the form control.
// It renders the read or edit template based on the HTTP method and query parameters.
func (c *FormControl[F]) GinHandler(g *gin.Context) {
	switch g.Request.Method {
	case http.MethodGet:
		if _, exists := g.GetQuery("edit"); exists {
			c.RenderEdit(g)
			return
		}
		c.RenderRead(g)
	case http.MethodPost:
		if err := g.ShouldBind(&c.Field); err != nil {
			c.SetError(err)
			c.RenderEdit(g)
			return
		}
		if err := c.Save(); err != nil {
			c.SetError(err)
			c.RenderEdit(g)
			return
		}
		c.RenderRead(g)
	}
}

// RenderEdit renders the edit template of the form control.
func (c *FormControl[F]) RenderEdit(g *gin.Context) {
	g.HTML(http.StatusOK, c.GetEditTemplate(), gin.H{"Control": c})
}

// RenderRead renders the read template of the form control.
func (c *FormControl[F]) RenderRead(g *gin.Context) {
	g.HTML(http.StatusOK, c.GetReadTemplate(), gin.H{"Control": c})
}

// Save saves the form control.
// If a save function is provided, it will be used to save the form control.
func (c *FormControl[F]) Save() error {
	if c.SaveFunc != nil {
		return c.SaveFunc(c)
	}
	return nil
}

// SetError sets the error message of the form control.
// If the error is a validation error, it sets the error message to the first validation error tag.
// Otherwise, it sets the error message to the error's message.
func (c *FormControl[F]) SetError(err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		c.Error = ve[0].Tag()
	} else {
		c.Error = err.Error()
	}
}
