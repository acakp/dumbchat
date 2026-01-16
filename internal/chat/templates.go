package chat

import (
	"text/template"
)

type parsedTemplates struct {
	Err         error
	ChatTmpl    *template.Template
	MessageTmpl *template.Template
	LoginTmpl   *template.Template
}

func ParseTemplates() parsedTemplates {
	var ret parsedTemplates
	chatTmpl, err := template.ParseFiles("internal/web/templates/layout.html")
	if err != nil {
		return parsedTemplates{err, nil, nil, nil}
	}
	ret.ChatTmpl = chatTmpl

	msgTmpl, err := template.ParseFiles("internal/web/templates/message.html")
	if err != nil {
		return parsedTemplates{err, nil, nil, nil}
	}
	ret.MessageTmpl = msgTmpl

	loginTmpl, err := template.ParseFiles("internal/web/templates/login.html")
	if err != nil {
		return parsedTemplates{err, nil, nil, nil}
	}
	ret.LoginTmpl = loginTmpl

	return ret
}
