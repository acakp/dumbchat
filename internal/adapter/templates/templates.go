package templates

import (
	"fmt"
	"html/template"

	"github.com/acakp/dumbchat/web"
)

type ParsedTemplates struct {
	Err         error
	ChatTmpl    *template.Template
	MessageTmpl *template.Template
	LoginTmpl   *template.Template
}

func ParseTemplatesCmd() ParsedTemplates {
	layoutTmpl := template.New("layout")
	layoutTmpl, err := layoutTmpl.Parse(web.LayoutHTML)
	if err != nil {
		err = fmt.Errorf("error parsing templates in ParseTemplatesCmd: %w", err)
		return ParsedTemplates{err, nil, nil, nil}
	}
	return ParseTemplates(layoutTmpl)
}

func ParseTemplates(t *template.Template) ParsedTemplates {
	var ret ParsedTemplates
	_, err := t.Parse(web.ChatHTML)
	_, err = t.Parse(web.MessageHTML)
	if err != nil {
		err = fmt.Errorf("error parsing templates: %w", err)
		return ParsedTemplates{err, nil, nil, nil}
	}
	ret.ChatTmpl = t

	messageTmpl := template.New("msg")
	messageTmpl, err = messageTmpl.Parse(web.MessageHTML)
	if err != nil {
		err = fmt.Errorf("error parsing message template: %w", err)
		return ParsedTemplates{err, nil, nil, nil}
	}
	ret.MessageTmpl = messageTmpl

	loginTmpl := template.New("login")
	loginTmpl, err = loginTmpl.Parse(web.LoginHTML)
	if err != nil {
		err = fmt.Errorf("error parsing login template: %w", err)
		return ParsedTemplates{err, nil, nil, nil}
	}
	ret.LoginTmpl = loginTmpl

	return ret
}
