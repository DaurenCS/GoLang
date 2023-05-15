package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	BotToken         = "5912845667:AAG1ny8xgd3UN4lIJqc863-6WsEwrxwOnKc"
	botApi           = "https://api.telegram.org/bot"
	UnsplahAccessKey = "om46TtSTHE-l_wlCEfQnfr1Jo4vhk3JQwi93aHPRsrw"
)

type UnsplashPhoto struct {
	Urls struct {
		Regular string `json:"regular"`
	} `json:"urls"`
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func GetRandomPhoto() string {
	url := "https://api.unsplash.com/photos/random?client_id=" + UnsplahAccessKey
	data, err := http.Get(url)
	CheckError(err)
	defer data.Body.Close()

	var photo UnsplashPhoto
	err = json.NewDecoder(data.Body).Decode(&photo)
	CheckError(err)
	return photo.Urls.Regular
}

func main() {

	bot, err := tgbotapi.NewBotAPI(BotToken)
	CheckError(err)
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var counter int

	ch := make(chan bool)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)
 
	go func() {
		for update := range updates {
			if update.Message != nil {
				if update.Message.IsCommand() && strings.ToLower(update.Message.Command()) == "image" || strings.ToLower(update.Message.Text) == "image" {

					mu.Lock()
					counter++
					fmt.Println("Counter:", counter)
					mu.Unlock()
					wg.Add(1)

					go func() {
						defer wg.Done()

						msg := tgbotapi.NewMessage(update.Message.Chat.ID, GetRandomPhoto())
						msg.ReplyToMessageID = update.Message.MessageID
						bot.Send(msg)
						wg.Wait()
						ch <- true
					}()
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This commanddd not found!")
					msg.ReplyToMessageID = update.Message.MessageID
					bot.Send(msg)
				}
			}
		}
	}()
	for {
		select {
		case <-ch:
			wg.Wait()
		}
	}
}
