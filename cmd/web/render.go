package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type templateData struct {
	Data map[string]any
}

func (app *application) render(c *gin.Context, t string, td *templateData) {
	var tmpl *template.Template

	if app.cfg.UseCache {
		if templateFromMap, ok := app.templateMap[t]; ok {
			tmpl = templateFromMap
		}
	}
	if tmpl == nil {
		newTemplate, err := app.buildTemplateFromDisc(t)
		if err != nil {
			log.Println("Error building template:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		log.Println("building template from disk")
		tmpl = newTemplate
	}

	if td == nil {
		td = &templateData{}
	}

	c.Writer.Header().Set("Content-Type", "text/html")
	if err := tmpl.ExecuteTemplate(c.Writer, t, td); err != nil {
		log.Println("Error executing template:", err)
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (app *application) buildTemplateFromDisc(t string) (*template.Template, error) {
	templateSlice := []string{
		"./templates/base.layout.gohtml",
		"./templates/partials/header.partials.gohtml",
		"./templates/partials/footer.partials.gohtml",
		fmt.Sprintf("./templates/%s", t),
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		return nil, err
	}
	app.templateMap[t] = tmpl
	return tmpl, nil
}
