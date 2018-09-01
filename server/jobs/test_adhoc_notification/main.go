package main

import (
	"encoding/json"
	"letstalk/server/core/notifications"
	"letstalk/server/data"
	"letstalk/server/utility"

	"github.com/namsral/flag"
	"github.com/romana/rlog"
)

var (
	recipient      *int
	message        *string
	title          *string
	thumbnail      *string
	templatePath   *string
	templateParams *string
)

func main() {
	db, err := utility.GetDB()
	if err != nil {
		panic(err)
	}

	flag.IntVar(recipient, "recipient", 0, "Enter the userId of the recipient")
	flag.StringVar(message, "message", "", "Enter the message to send")
	flag.StringVar(title, "title", "", "Enter the title")
	flag.StringVar(thumbnail, "thumbnail", "", "Enter the thumbnail path")
	flag.StringVar(templatePath, "templatePath", "", "Enter the template path")
	flag.StringVar(templateParams, "templateParams", "", "Enter the template params")
	flag.Parse()
	var params map[string]string
	if err := json.Unmarshal([]byte(*templateParams), &params); err != nil {
		panic(err)
	}

	rlog.Infof(
		`Sending notification:
		\trecipient:%d
		\tmessage:%s
		\ttitle:%s
		\tthumbnail:%s
		\ttemplate:%s
		\tparams:%v`, recipient, message, title, thumbnail, templatePath, params)

	if err = notifications.CreateAdHocNotification(
		db,
		data.TUserID(*recipient),
		*title,
		*message,
		thumbnail,
		*templatePath,
		params,
	); err != nil {
		panic(err)
	}
}
