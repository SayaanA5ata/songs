package dat

import (
	"log"
	"time"
)

func ConvertDateToDB(inputDate string) *time.Time {
	parsedDate, err := time.Parse("02.01.2006", inputDate) // Преобразуем строку в time.Time
	if err != nil {
		log.Fatal("Failed to parse date:", err)
	}
	return &parsedDate
}

func ConvertDateToUser(raw_date time.Time) string {
	outputDate := raw_date.Format("02.01.2006")
	return outputDate
}
