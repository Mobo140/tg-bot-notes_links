package main

import (
	"log"

	"github.com/Mobo140/projects/go-pocket-sdk"
	"github.com/Mobo140/projects/tg-bot-notes_links/pkg/config"
	"github.com/Mobo140/projects/tg-bot-notes_links/pkg/repository"
	"github.com/Mobo140/projects/tg-bot-notes_links/pkg/repository/boltdb"
	"github.com/Mobo140/projects/tg-bot-notes_links/pkg/server"
	"github.com/Mobo140/projects/tg-bot-notes_links/pkg/telegram"
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load() //Чтение переменных окружения

	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	godotenv.Load()
	pocketClient, err := pocket.NewClient(cfg.PocketConsumerKey)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, "http://localhost/")

	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, "https://t.me/pocket_140_bot")

	go func() { //Вынесли в отдельную горутину так как при запуске бота запускается в бесконечном цикле считывание обновлений и функция main будет заблокирована
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()
	//тоже является блокирующей операцией так как вызывает метод Listed and Serve. Это позволяет запустить телеграмм бота и сервер в одном приложении
	if err := authorizationServer.Start(); err != nil {
		log.Fatal(err)
	}

}

func initDB() (*bolt.DB, error) {

	db, err := bolt.Open("bot db", 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}

		tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}
