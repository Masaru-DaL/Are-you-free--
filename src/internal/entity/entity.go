package entity

type User struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	Password  string `db:"password"`
	Email     string `db:"email"`
	IsAdmin   bool   `db:"is_admin"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}
