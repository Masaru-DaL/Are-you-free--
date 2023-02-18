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

type DateFreeTime struct {
	ID        int    `db:"id"`
	UserID    int    `db:"user_id"`
	Year      int    `db:"year"`
	Month     int    `db:"month"`
	Day       int    `db:"day"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

type StartFreeTime struct {
	ID             int    `db:"id"`
	DateFreeTimeID int    `db:"date_free_time_id"`
	Hour           int    `db:"start_time_hour"`
	Minute         int    `db:"start_time_minute"`
	CreatedAt      string `db:"created_at"`
	UpdatedAt      string `db:"updated_at"`
}

type EndFreeTime struct {
	ID             int    `db:"id"`
	DateFreeTimeID int    `db:"date_free_time_id"`
	Hour           int    `db:"end_time_hour"`
	Minute         int    `db:"start_time_minute"`
	CreatedAt      string `db:"created_at"`
	UpdatedAt      string `db:"updated_at"`
}

// type FreeTime struct {
// 	ID              int    `db:"id"`
// 	UserID          int    `db:"user_id"`
// 	StartTimeID     int    `db:"start_time_id"`
// 	EndTimeID       int    `db:"end_time_id"`
// 	Year            int    `db:"year"`
// 	Month           int    `db:"month"`
// 	Day             int    `db:"day"`
// 	StartTimeHour   int    `db:"start_time_hour"`
// 	StartTimeMinute int    `db:"start_time_minute"`
// 	EndTimeHour     int    `db:"end_time_hour"`
// 	EndTimeMinute   int    `db:"end_time_minute"`
// 	CreatedAt       string `db:"created_at"`
// 	UpdatedAt       string `db:"updated_at"`
// }
