package postgres

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"sf-news/pkg/storage"
)

// Хранилище данных
type Storage struct {
	dbPool *pgxpool.Pool
}

// Конструктор объекта хранилища
func New(constr string) (*Storage, error) {
	dbPool, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	return &Storage{dbPool: dbPool}, nil
}

// PushPosts новости в хранилище
func (s *Storage) PushPosts(posts []storage.Post) error {
	for _, post := range posts {
		_, err := s.dbPool.Exec(context.Background(), `
			INSERT INTO posts (pub_time, link, title, content)
			VALUES ($1, $2, $3, $4);`,
			post.PubTime,
			post.Link,
			post.Title,
			post.Content,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// PopPosts получает самые свежие новости из хранилища
func (s *Storage) PopPosts(count int) ([]storage.Post, error) {
	rows, err := s.dbPool.Query(context.Background(),
		`SELECT p.id, p.pub_time, p.link, p.title, p.content
		FROM posts p
		ORDER BY p.pub_time DESC
		LIMIT $1;`,
		count,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []storage.Post
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(&p.ID, &p.PubTime, &p.Link, &p.Title, &p.Content)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, rows.Err()
}
