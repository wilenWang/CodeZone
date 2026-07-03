package users

import (
	"context"
	"database/sql"
)

type User struct {
	ID          int64   `json:"id"`
	WorkspaceID int64   `json:"workspaceId"`
	Username    string  `json:"username"`
	DisplayName string  `json:"displayName"`
	AvatarURL   *string `json:"avatarUrl"`
	UserType    string  `json:"userType"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (User, error) {
	const query = `
		SELECT id, workspace_id, username, display_name, avatar_url, user_type
		FROM users
		WHERE username = ?
		LIMIT 1`
	return scanUser(r.db.QueryRowContext(ctx, query, username))
}

func (r *Repository) FindByWorkspaceUsername(ctx context.Context, workspaceID int64, username string) (User, error) {
	const query = `
		SELECT id, workspace_id, username, display_name, avatar_url, user_type
		FROM users
		WHERE workspace_id = ? AND username = ?
		LIMIT 1`
	return scanUser(r.db.QueryRowContext(ctx, query, workspaceID, username))
}

func (r *Repository) FindByID(ctx context.Context, id int64) (User, error) {
	const query = `
		SELECT id, workspace_id, username, display_name, avatar_url, user_type
		FROM users
		WHERE id = ?
		LIMIT 1`
	return scanUser(r.db.QueryRowContext(ctx, query, id))
}

func (r *Repository) List(ctx context.Context, workspaceID int64) ([]User, error) {
	const query = `
		SELECT id, workspace_id, username, display_name, avatar_url, user_type
		FROM users
		WHERE workspace_id = ?
		ORDER BY display_name, id`
	rows, err := r.db.QueryContext(ctx, query, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

type userScanner interface {
	Scan(dest ...any) error
}

func scanUser(scanner userScanner) (User, error) {
	var user User
	var avatarURL sql.NullString
	err := scanner.Scan(
		&user.ID,
		&user.WorkspaceID,
		&user.Username,
		&user.DisplayName,
		&avatarURL,
		&user.UserType,
	)
	if err != nil {
		return User{}, err
	}
	if avatarURL.Valid {
		user.AvatarURL = &avatarURL.String
	}
	return user, nil
}
