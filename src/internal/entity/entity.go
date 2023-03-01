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
	ID        int         `db:"id"`
	UserID    int         `db:"user_id"`
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
