package chat

import (
	"database/sql"
	"net/http"
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

func showAllMessages(w http.ResponseWriter, db *sql.DB, msgTmpl *template.Template, isAdmin bool) error {
	msgs, err := getMessages(db)
	if err != nil {
		return err
	}
	for _, msg := range msgs {
		nice := struct {
			Msg     Message
			IsAdmin bool
		}{
			msg,
			isAdmin,
		}
		_ = msgTmpl.ExecuteTemplate(w, "msg", nice)
	}
	return nil
}
