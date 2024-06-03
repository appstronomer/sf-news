package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sf-news/pkg/storage"
	"sf-news/pkg/storage/memdb"
	"testing"
	"unicode/utf8"
)

func TestAPI_registerStatic(t *testing.T) {
	// Меняется только желаемый текст и url
	type fields struct {
		wantText string
		url      string
	}
	// Тесты
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "static file",
			fields: fields{
				wantText: "<!doctype html><meta charset=utf-8><title>Hello, world!</title>",
				url:      "/",
			},
		},
		{
			name: "static file",
			fields: fields{
				wantText: "V2hvcyB5b3VyIGRhZGR5PyE=",
				url:      "/test-route.txt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Настройки
			// Подготовка API
			api := New("./testdata", memdb.New())

			// Отправка запроса
			req := httptest.NewRequest(http.MethodGet, tt.fields.url, nil)
			rr := httptest.NewRecorder()
			api.router.ServeHTTP(rr, req)

			// Получение ответа и десерализация
			if !(rr.Code == http.StatusOK) {
				t.Fatalf("http response status: got %d; want %d", rr.Code, http.StatusOK)
			}
			b, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Fatalf("http response body read error %v", err)
			}
			if !utf8.Valid(b) {
				t.Fatalf("http response body is invalid utf8: %v", b)
			}
			gotText := string(b)

			// Проверка контента
			if gotText != tt.fields.wantText {
				t.Fatalf("http response text got \"%s\"; want \"%s\"", gotText, tt.fields.wantText)
			}
		})
	}

}

func TestAPI_handlePosts(t *testing.T) {
	// Настройки
	sendLen := 3
	wantLen := sendLen - 1
	// Запрашивается ответов меньше, чем есть в хранилище: [:wantLen]
	wantPosts := helpMakePostSl(sendLen, true)[:wantLen]

	// Заолнение данными
	db := memdb.New()
	db.PushPosts(helpMakePostSl(sendLen, false))
	api := New("./testdata", db)

	// Отправка запроса
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/news/%d", wantLen), nil)
	rr := httptest.NewRecorder()
	api.router.ServeHTTP(rr, req)

	// Получение ответа и десерализация
	if !(rr.Code == http.StatusOK) {
		t.Fatalf("http response status: got %d; want %d", rr.Code, http.StatusOK)
	}
	b, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("http response body read error %v", err)
	}
	var gotPosts []storage.Post
	err = json.Unmarshal(b, &gotPosts)
	if err != nil {
		t.Fatalf("http response json decode error: %v", err)
	}

	// Проверка длины массива
	if len(gotPosts) != wantLen {
		t.Fatalf("posts count got %v, want %v", len(gotPosts), wantLen)
	}
	// Проверка комплектности массива
	for i := 0; i < wantLen; i++ {
		if !reflect.DeepEqual(gotPosts[i], wantPosts[i]) {
			t.Errorf("http response item idx=%d got %#v; want %#v", i, gotPosts[i], wantPosts[i])
		}
	}
}

// Создаёт воспроизводимую последовательность постов
func helpMakePostSl(count int, setId bool) []storage.Post {
	if count <= 0 {
		panic(fmt.Sprintf("helpMakePostSl(%d) - negative count got", count))
	}
	posts := make([]storage.Post, 0, count)
	for i, id := count, 1; i > 0; i, id = i-1, id+1 {
		post := storage.Post{
			PubTime: 1717060660 + int64(1000*i),
			Link:    fmt.Sprintf("https://fake.com/%d", i),
			Title:   fmt.Sprintf("Post %d title", i),
			Content: fmt.Sprintf("Post %d content", i),
		}
		if setId {
			post.ID = id
		}
		posts = append(posts, post)
	}
	return posts
}
