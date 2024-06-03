package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"
)

// Конфигурация тестового клиента
type Config struct {
	IsExactCheck bool   `json:"is_exact_check"`
	Url          string `json:"url"`
	PostsWant    []Post `json:"posts_want"`
}

// Post - новостная публикация
type Post struct {
	ID      int    `json:"ID"`
	PubTime int64  `json:"PubTime"`
	Link    string `json:"Link"`
	Title   string `json:"Title"`
	Content string `json:"Content"`
}

func main() {
	// Аргумент приложения
	if len(os.Args) != 2 {
		log.Fatal("please, provide path to config file as the only argument\n")
	}

	// Чтение и десериализации конфигурации
	cfgBytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	var cfg Config
	err = json.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Несколько попыток проверить новости
	for i := 1; i <= 4; i++ {
		time.Sleep(time.Second * 5)
		err = check(cfg)
		if err == nil {
			log.Printf("Try %d: success", i)
			break
		}
		log.Printf("Try %d: failure", i)
		log.Println(err)
	}

	if err != nil {
		log.Println("[FAILURE] all tries failed")
		log.Fatal(err)
	}
	log.Println("[SUCCESS] all checks were successful")
}

func check(cfg Config) error {
	// Получение и десериализации новостей из аггрегатора новостей
	res, err := http.Get(cfg.Url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	postsBytes, err := io.ReadAll(res.Body)
	var postsGot []Post
	err = json.Unmarshal(postsBytes, &postsGot)
	if err != nil {
		return err
	}

	// Сравнение новостей
	if cfg.IsExactCheck {
		if len(postsGot) != len(cfg.PostsWant) {
			return fmt.Errorf("posts count got %v; want %v", len(postsGot), len(cfg.PostsWant))
		}
		for i := 0; i < len(cfg.PostsWant); i++ {
			// ID сложно предсказать, поэтому просто зануляем их
			postsGot[i].ID = 0
			cfg.PostsWant[i].ID = 0
			// Проверка остальных полей
			if !reflect.DeepEqual(postsGot[i], cfg.PostsWant[i]) {
				return fmt.Errorf("posts[%d] got %#v; want %#v", i, postsGot[i], cfg.PostsWant[i])
			}
		}
	}
	return nil
}
