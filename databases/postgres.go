package databases

import (
	"context"
	"database/sql"
	"log"

	"github.com/adrisongomez/project-go/models"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, error := sql.Open("postgres", url)
	if error != nil {
		return nil, error
	}
	return &PostgresRepository{db}, nil
}

func (repo *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	_, error := repo.db.ExecContext(
		ctx,
		"INSERT INTO users (id, email, password) VALUES ($1, $2, $3)",
		user.Id,
		user.Email,
		user.Password,
	)
	return error
}

func (repo *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := repo.db.QueryContext(
		ctx,
		"SELECT id, email, password FROM users WHERE email = $1",
		email,
	)
	defer handleCloseCursor(rows)
	if err != nil {
		return nil, err
	}
	return mapFromRowsToUser(rows)
}

func (repo *PostgresRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	rows, err := repo.db.QueryContext(
		ctx,
		"SELECT id, email, password FROM users WHERE id = $1",
		id,
	)
	defer handleCloseCursor(rows)
	if err != nil {
		return nil, err
	}
	return mapFromRowsToUser(rows)
}

func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}

func (repo *PostgresRepository) InsertPost(ctx context.Context, post *models.Post) error {
	_, err := repo.db.ExecContext(
		ctx,
		"INSERT INTO posts (id, post_content, user_id) VALUES ($1, $2, $3)",
		post.Id,
		post.PostContent,
		post.UserId,
	)
	return err
}

func (repo *PostgresRepository) GetPostById(ctx context.Context, id string) (*models.Post, error) {
	rows, err := repo.db.QueryContext(
		ctx,
		"SELECT id, user_id, post_content, created_at FROM posts WHERE id = $1",
		id,
	)
	defer handleCloseCursor(rows)
	if err != nil {
		return nil, err
	}
	return mapFromRowsToPost(rows)
}

func (repo *PostgresRepository) UpdatePost(ctx context.Context, post *models.Post) error {
	_, err := repo.db.ExecContext(
		ctx,
		"UPDATE posts SET post_content = $1 WHERE id = $2 AND user_id = $3",
		post.PostContent,
		post.Id,
		post.UserId,
	)
	return err

}

func (repo *PostgresRepository) DeletePost(ctx context.Context, id, userId string) error {
	_, err := repo.db.ExecContext(
		ctx,
		"DELETE FROM posts WHERE id = $1 AND user_id = $2",
		id,
		userId,
	)
	return err
}

func (repo *PostgresRepository) ListPost(ctx context.Context, page uint64) ([]*models.Post, error) {
	rows, err := repo.db.QueryContext(
		ctx,
		"SELECT id, post_content, user_id, created_at FROM posts LIMIT $1 OFFSET $2",
		2,
		page*2,
	)

	if err != nil {
		return nil, err
	}
	defer handleCloseCursor(rows)

	var posts []*models.Post

	for rows.Next() {
		var post = models.Post{}

		if err = rows.Scan(&post.Id, &post.UserId, &post.PostContent, &post.CreatedAt); err == nil {
			posts = append(posts, &post)
		} 
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func handleCloseCursor(rows *sql.Rows) {
	err := rows.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func mapFromRowsToUser(rows *sql.Rows) (*models.User, error) {
	user := models.User{}
	for rows.Next() {
		if err := rows.Scan(&user.Id, &user.Email, &user.Password); err != nil {
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}

func mapFromRowsToPost(rows *sql.Rows) (*models.Post, error) {
	post := models.Post{}
	for rows.Next() {
		if err := rows.Scan(&post.Id, &post.UserId, &post.PostContent, &post.CreatedAt); err != nil {
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &post, nil
}
