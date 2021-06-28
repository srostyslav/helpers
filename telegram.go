package helpers

import (
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type adminBot struct {
	bot         *tgbotapi.BotAPI
	chatID      int64
	projectName string
}

var AdminBot = &adminBot{}

func (a *adminBot) Init(chatID int64, projectName, token string) (err error) {
	a.chatID, a.projectName = chatID, projectName

	if a.bot, err = tgbotapi.NewBotAPI(token); err != nil {
		ErrorLogger.Println("Cannot init bot", err)
	}

	return err
}

func (a *adminBot) SendDocument(filePath, text string) {
	if a.bot == nil {
		return
	}

	msg := tgbotapi.NewDocumentUpload(a.chatID, filePath)
	msg.Caption = text
	if _, err := a.bot.Send(msg); err != nil {
		ErrorLogger.Println("NewDocumentUpload", err)
	}
}

func (a *adminBot) sendText(txt string) {
	msg := tgbotapi.NewMessage(a.chatID, a.projectName+" "+txt)
	if _, err := a.bot.Send(msg); err != nil {
		ErrorLogger.Println("bot Send", err)
	}
}

func (a *adminBot) writeFile(path string, data []byte) error {
	file := &File{Name: path}
	file.Create()

	return file.Write(data)
}

func (a *adminBot) SendError(data interface{}, txt string) {
	if a.bot == nil {
		return
	}

	if data == nil {
		a.sendText(txt)
	} else {
		path := "errors.txt"
		if body, err := json.Marshal(&data); err != nil {
			ErrorLogger.Println("cannot unmarshal error", err)
		} else if err := a.writeFile(path, body); err != nil {
			ErrorLogger.Println("write error", err)
		} else {
			a.SendDocument(path, txt)
		}
	}
}

func (a *adminBot) SetWebhook(url string) (err error) {
	_, err = a.bot.SetWebhook(tgbotapi.NewWebhookWithCert(url, nil))
	return err
}
