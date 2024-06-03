package memdb

import (
	"slices"
	"sort"
	"sync"

	"sf-news/pkg/storage"
)

// Хранилище данных
type Storage struct {
	mu       sync.Mutex
	idx      int
	isSorted bool
	sl       []storage.Post
}

// Конструктор объекта хранилища
func New() *Storage {
	return &Storage{isSorted: true}
}

// PushPosts новости в хранилище
func (s *Storage) PushPosts(posts []storage.Post) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, post := range posts {
		idxFound := slices.IndexFunc(s.sl, func(p storage.Post) bool {
			return p.Link == post.Link
		})
		if idxFound == -1 {
			s.idx += 1
			post.ID = s.idx
			s.sl = append(s.sl, post)
		}
	}
	s.isSorted = false
	return nil
}

// PopPosts получает самые свежие новости из хранилища
func (s *Storage) PopPosts(count int) ([]storage.Post, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.isSorted {
		s.sort()
	}
	if count >= len(s.sl) {
		return s.sl, nil
	}
	return s.sl[:count], nil
}

func (s *Storage) sort() {
	sort.Slice(s.sl, func(i, j int) bool {
		return s.sl[i].PubTime > s.sl[j].PubTime
	})
	s.isSorted = true
}
