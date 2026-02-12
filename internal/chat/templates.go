package chat

import (
	"html/template"

	"github.com/acakp/dumbchat/web"
)

func ParseTemplatesCmd() parsedTemplates {
	layoutTmpl := template.New("layout")
	layoutTmpl, err := layoutTmpl.Parse(web.LayoutHTML)
	if err != nil {
		return parsedTemplates{err, nil, nil, nil}
	}
	return ParseTemplates(layoutTmpl)
}

func ParseTemplates(t *template.Template) parsedTemplates {
	var ret parsedTemplates
	_, err := t.Parse(web.ChatHTML)
	_, err = t.Parse(web.MessageHTML)
	if err != nil {
		return parsedTemplates{err, nil, nil, nil}
	}
	ret.ChatTmpl = t

	messageTmpl := template.New("msg")
	messageTmpl, err = messageTmpl.Parse(web.MessageHTML)
	if err != nil {
		return parsedTemplates{err, nil, nil, nil}
	}
	ret.MessageTmpl = messageTmpl

	loginTmpl := template.New("login")
	loginTmpl, err = loginTmpl.Parse(web.LoginHTML)
	if err != nil {
		return parsedTemplates{err, nil, nil, nil}
	}
	ret.LoginTmpl = loginTmpl

	return ret
}
