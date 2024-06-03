package storage

// Post - новостная публикация
type Post struct {
	ID      int    `json:"ID"`
	PubTime int64  `json:"PubTime"`
	Link    string `json:"Link"`
	Title   string `json:"Title"`
	Content string `json:"Content"`
}

// Interface задаёт контракт на работу с БД.
type StorageIface interface {
	// Получение самых новых публикаций из хранилища новостей
	PopPosts(count int) ([]Post, error)
	// Добавление публикаций в хранилище новостей
	PushPosts([]Post) error
}
