-- name: CreateUser :one
  INSERT INTO users ("user_name", "email", "password_hash", "bio")
  VALUES($1,$2,$3,$4)
  RETURNING id;

-- name: GetUserById :one
 Select id, user_name, email, password_hash, bio, created_at, updated_at
 FROM users
 WHERE id = $1;

-- name: GetUserByEmail :one
 Select id, user_name, email, password_hash, bio, created_at, updated_at
 FROM users
 WHERE email = $1;