package entity

type User struct {
	ID        string `db:"id"`
	Name      string `db:"name"`
	Password  string `db:"password"`
	Email     string `db:"email"`
	IsAdmin   bool   `db:"is_admin"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

type DateFreeTime struct {
	ID        int         `db:"id"`
	UserID    string      `db:"user_id"`
	Year      string      `db:"year"`
	Month     string      `db:"month"`
	Day       string      `db:"day"`
	CreatedAt string      `db:"created_at"`
	UpdatedAt string      `db:"updated_at"`
	FreeTimes []*FreeTime `db:"free_times"`
}

type FreeTime struct {
	ID             int    `db:"id"`
	DateFreeTimeID int    `db:"date_free_time_id"`
	StartHour      int    `db:"start_hour"`
	StartMinute    int    `db:"start_minute"`
	EndHour        int    `db:"end_hour"`
	EndMinute      int    `db:"end_minute"`
	CreatedAt      string `db:"created_at"`
	UpdatedAt      string `db:"updated_at"`
}

type SharedUser struct {
	ID           string `db:"id"`
	UserID       string `db:"user_id"`
	SharedUserID string `db:"shared_user_id"`
	CreatedAt    string `db:"created_at"`
	UpdatedAt    string `db:"updated_at"`
}
