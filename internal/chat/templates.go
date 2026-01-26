package chat

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/acakp/dumbchat/internal/web"
)

func ParseTemplatesCmd() parsedTemplates {
	layoutTmpl, err := template.ParseFiles("internal/web/templates/layout.html")
	if err != nil {
		return parsedTemplates{err, nil, nil, nil}
	}
	return ParseTemplates(layoutTmpl)
}

func ParseTemplates(t *template.Template) parsedTemplates {
	var ret parsedTemplates
	// chatTmpl := template.New("chat")
	// chatTmpl = t
	// _, err := chatTmpl.Parse(web.ChatHTML)
	_, err := t.Parse(web.ChatHTML)
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

func showAllMessages(w http.ResponseWriter, db *sql.DB, msgTmpl *template.Template, msv MessageView) error {
	msgs, err := getMessages(db)
	if err != nil {
		return err
	}
	for _, msg := range msgs {
		msv.Msg = msg
		err = msgTmpl.ExecuteTemplate(w, "msg", msv)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
