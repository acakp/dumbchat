package web

import _ "embed"

// //go:embed templates/*
// var TemplateFS embed.FS

//go:embed templates/chat.html
var ChatHTML string

//go:embed templates/message.html
var MessageHTML string

//go:embed templates/login.html
var LoginHTML string

//go:embed templates/layout.html
var LayoutHTML string
