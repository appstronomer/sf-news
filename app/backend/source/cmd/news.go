package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"sf-news/pkg/api"
	"sf-news/pkg/output"
	"sf-news/pkg/rss"
	"sf-news/pkg/storage"
	"sf-news/pkg/storage/postgres"
)

// Сообщение - мануал
const helpMessage string = `news PATH_CONFIG PATH_WWW DB_CONSTR
PATH_CONFIG - path to application config.json file
PATH_WWW - path to static web-interface files
DB_CONSTR - postgres connection config string`

func main() {
	// Создание логгера
	out := output.Make(os.Stdout, os.Stderr)

	// Проверка аргументов запуска
	if len(os.Args) != 4 {
		out.Err(fmt.Sprintf("ERROR: invalid arguments passed to application: %v\n\nSYNOPSIS\n%v", os.Args, helpMessage))
		os.Exit(1)
	}

	// Получение конфигурации
	cfg, err := readConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Подключение к БД
	// Пример:"postgres://user:password@postgres:5432/sf"
	db, err := postgres.New(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	// Создание роутера
	api := api.New(os.Args[2], db)

	// Запуск парсеров
	chPosts := parseInit(out, cfg)

	// Аггрегация новостей в БД
	go dbLoop(out, db, chPosts)

	err = http.ListenAndServe(":80", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}

// КОНФИГУРАЦИЯ ПРИЛОЖЕНИЯ
// Контейнер под конфиг приложения
type config struct {
	RssUrls       []string `json:"rss"`
	RequestPeriod uint32   `json:"request_period"`
}

// Чтение и валидация конфига
func readConfig(path string) (config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return config{}, err
	}
	var cfg config
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		return config{}, err
	}
	for _, urlItem := range cfg.RssUrls {
		_, err := url.ParseRequestURI(urlItem)
		if err != nil {
			return config{}, err
		}
	}
	return cfg, nil
}

// ЧТЕНИЕ НОВОСТЕЙ
// Инициализация всех горутин, осуществляющих чтение новостей
func parseInit(out output.Output, cfg config) <-chan []storage.Post {
	chPosts := make(chan []storage.Post)
	delay := time.Duration(cfg.RequestPeriod) * time.Minute
	for _, url := range cfg.RssUrls {
		go parseLoop(out, chPosts, url, delay)
	}
	return chPosts
}

// Цкл чтения новостей конкретной горутиной
func parseLoop(out output.Output, chPosts chan<- []storage.Post, url string, delay time.Duration) {
	for {
		posts, err := parseUrl(url)
		if err != nil {
			out.Err(err)
			continue
		}
		chPosts <- posts
		time.Sleep(delay)
	}
}

// Итерация получения и десереализации новостей
func parseUrl(url string) ([]storage.Post, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	posts, err := rss.Parse(res.Body)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// ЗАПИСЬ В БД
// Цикл записи новостей в БД
func dbLoop(out output.Output, db *postgres.Storage, chPosts <-chan []storage.Post) {
	for posts := range chPosts {
		err := db.PushPosts(posts)
		if err != nil {
			out.Err(err)
		}
	}
}
