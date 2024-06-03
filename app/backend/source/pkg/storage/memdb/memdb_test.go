package memdb

import (
	"fmt"
	"reflect"
	"sf-news/pkg/storage"
	"testing"
)

func TestStorage_PushPosts(t *testing.T) {
	// Подготовка данных
	const itemCount int = 4
	const itemDup int = 1
	posts := helpMakePostSl(itemCount, false)
	// Добавляется статья-дубликат, которая должна быть отсеяна хранилищем
	posts = append(posts, helpMakePost(itemDup, false))

	// Добавление статей в хранилище
	var store storage.StorageIface = New()
	err := store.PushPosts(posts)
	if err != nil {
		t.Fatalf("Storage.PushPosts() error = %v", err)
	}

	// Получение внутреннего слайса для проверки
	memdbStorage := store.(*Storage)
	sliceInner := memdbStorage.sl

	// Проверка комплектности списка
	if len(sliceInner) != itemCount {
		t.Fatalf("len(Storage.sl) got %v; want %v", len(sliceInner), itemCount)
	}

	// Проверка отсутствия статьи-дубликата, наличия ID у всех статей
	wantId := 0
	for _, gotItem := range sliceInner {
		wantId += 1
		wantItem := helpMakePost(wantId, true)
		if !reflect.DeepEqual(gotItem, wantItem) {
			t.Errorf("Storage.PopPosts()[%v] got %#v; want %#v", wantId, gotItem, wantItem)
		}
	}

}

func TestStorage_PopPosts(t *testing.T) {
	// Подготовка данных
	const itemCount int = 4
	posts := helpMakePostSl(itemCount, false)

	// Добавление статей в хранилище
	var store storage.StorageIface = New()
	err := store.PushPosts(posts)
	if err != nil {
		t.Fatalf("Storage.PushPosts() error = %v", err)
	}

	// Получение отсортированного списка статей
	news, err := store.PopPosts(itemCount + 1)
	if err != nil {
		t.Fatalf("Storage.PopPosts() error = %v", err)
	}

	// Проверка комплектности списка
	if len(news) != itemCount {
		t.Fatalf("len(Storage.PopPosts()) got %v; want %v", len(news), itemCount)
	}

	// Проверка сортировки
	gotIdx := 0
	wantId := itemCount
	for gotIdx < len(news) {
		wantItem := helpMakePost(wantId, true)
		gotItem := news[gotIdx]
		if !reflect.DeepEqual(gotItem, wantItem) {
			t.Errorf("Storage.PopPosts()[%v] got %#v; want %#v", gotIdx, gotItem, wantItem)
		}
		gotIdx += 1
		wantId -= 1
	}
}

// Воспроизводимо создаёт статью на основе одного чила
// с возможностью не задавать ID статей
func helpMakePost(idx int, setId bool) storage.Post {
	post := storage.Post{
		PubTime: int64(idx * 100),
		Link:    fmt.Sprintf("https://test.com/post/%v", idx),
		Title:   fmt.Sprintf("Post %v", idx),
		Content: fmt.Sprintf("Content %v", idx),
	}
	if setId {
		post.ID = idx
	}
	return post
}

// Воспроизводимо создаёт список статей на основе требуемого
// количества с возможностью не задавать ID статей
func helpMakePostSl(itemCount int, setId bool) []storage.Post {
	posts := make([]storage.Post, 0, itemCount)
	for i := 1; i <= itemCount; i++ {
		posts = append(posts, helpMakePost(i, setId))
	}
	return posts
}
