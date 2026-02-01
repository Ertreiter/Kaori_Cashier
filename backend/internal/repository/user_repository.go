package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/kaori/backend/internal/model"
)

// GetByEmail finds a user by email
func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, email, password_hash, name, role, store_id, pin, is_active, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_active = true
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
		&user.StoreID, &user.PIN, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// GetByID finds a user by ID
func (r *UserRepository) GetByID(id uuid.UUID) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, email, password_hash, name, role, store_id, pin, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
		&user.StoreID, &user.PIN, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// Create creates a new user
func (r *UserRepository) Create(user *model.User) error {
	query := `
		INSERT INTO users (email, password_hash, name, role, store_id, pin, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		user.Email, user.PasswordHash, user.Name, user.Role, user.StoreID, user.PIN, user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

// List returns all users for a store
func (r *UserRepository) List(storeID *uuid.UUID) ([]model.User, error) {
	var users []model.User
	var query string
	var args []interface{}

	if storeID != nil {
		query = `
			SELECT id, email, password_hash, name, role, store_id, pin, is_active, created_at, updated_at
			FROM users
			WHERE store_id = $1 OR store_id IS NULL
			ORDER BY created_at DESC
		`
		args = append(args, storeID)
	} else {
		query = `
			SELECT id, email, password_hash, name, role, store_id, pin, is_active, created_at, updated_at
			FROM users
			ORDER BY created_at DESC
		`
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		if err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
			&user.StoreID, &user.PIN, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// Update updates a user
func (r *UserRepository) Update(user *model.User) error {
	query := `
		UPDATE users
		SET name = $2, role = $3, store_id = $4, pin = $5, is_active = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	_, err := r.db.Exec(query, user.ID, user.Name, user.Role, user.StoreID, user.PIN, user.IsActive)
	return err
}

// Delete deletes a user (soft delete by setting is_active = false)
func (r *UserRepository) Delete(id uuid.UUID) error {
	query := `UPDATE users SET is_active = false WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// SaveRefreshToken saves a refresh token
func (r *UserRepository) SaveRefreshToken(token *model.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.db.QueryRow(query, token.UserID, token.Token, token.ExpiresAt).
		Scan(&token.ID, &token.CreatedAt)
}

// GetRefreshToken gets a refresh token
func (r *UserRepository) GetRefreshToken(token string) (*model.RefreshToken, error) {
	rt := &model.RefreshToken{}
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM refresh_tokens
		WHERE token = $1 AND expires_at > CURRENT_TIMESTAMP
	`
	err := r.db.QueryRow(query, token).Scan(
		&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return rt, nil
}

// DeleteRefreshToken deletes a refresh token
func (r *UserRepository) DeleteRefreshToken(token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := r.db.Exec(query, token)
	return err
}

// DeleteUserRefreshTokens deletes all refresh tokens for a user
func (r *UserRepository) DeleteUserRefreshTokens(userID uuid.UUID) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}
