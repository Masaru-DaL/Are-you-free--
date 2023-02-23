package time

import (
	"fmt"
	"src/internal/entity"
	"src/internal/utils/time"
	"testing"
)

func TestCheckInputTime(t *testing.T) {
	// テスト用のフリータイムを作成
	df := entity.DateFreeTime{
		ID:        1,
		UserID:    100,
		Year:      2022,
		Month:     2,
		Day:       23,
		CreatedAt: "2022-02-23 00:00:00",
		UpdatedAt: "2022-02-23 00:00:00",
	}
	ft1 := &entity.FreeTime{
		ID:             1,
		DateFreeTimeID: df.ID,
		StartHour:      10,
		StartMinute:    0,
		EndHour:        12,
		EndMinute:      0,
		CreatedAt:      "2022-02-23 00:00:00",
		UpdatedAt:      "2022-02-23 00:00:00",
	}
	ft2 := &entity.FreeTime{
		ID:             2,
		DateFreeTimeID: df.ID,
		StartHour:      14,
		StartMinute:    0,
		EndHour:        16,
		EndMinute:      0,
		CreatedAt:      "2022-02-23 00:00:00",
		UpdatedAt:      "2022-02-23 00:00:00",
	}
	df.FreeTimes = []*entity.FreeTime{ft1, ft2}

	// テストケース1: 既存のフリータイムと時間が重複している場合
	startHour := 10
	startMinute := 0
	endHour := 11
	endMinute := 0

	expect := false
	result := time.CheckInputTime(startHour, startMinute, endHour, endMinute, &df)
	if result != expect {
		t.Errorf("Expected is '%v', but got '%v'", expect, result)
	}
	fmt.Println(result)

	// テストケース2: 既存のフリータイムと時間が重複していない場合
	startHour = 9
	startMinute = 0
	endHour = 10
	endMinute = 0
	expect = true
	result = time.CheckInputTime(startHour, startMinute, endHour, endMinute, &df)
	if result != expect {
		t.Errorf("Expected is '%v', but got '%v'", expect, result)
	}
	fmt.Println(result)
}
