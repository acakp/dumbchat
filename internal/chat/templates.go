package chat

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"
)

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

func showAllMessages(w http.ResponseWriter, db *sql.DB, msgTmpl *template.Template, msv MessageView) error {
	msgs, err := getMessages(db)
	if err != nil {
		return err
	}
	for _, msg := range msgs {
		msv.Msg = msg
		msv.IsAdmin = true
		fmt.Println(msv.URLs.Delete)
		err = msgTmpl.ExecuteTemplate(w, "msg", msv)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
