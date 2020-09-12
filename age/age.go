package age

import "time"

func GetYears(now, birthDate time.Time) int {
	if now.Before(birthDate) {
		return -1
	}
	tomorrow := now.AddDate(0, 0, 1)
	age := tomorrow.Year() - birthDate.Year()

	if tomorrow.Sub(birthDate.AddDate(age, 0, 0)) > 0 {
		return age
	}
	return age - 1
}

func GetDays(now, birthDate time.Time) int {
	if now.Before(birthDate) {
		return 0
	}

	duration := now.Sub(birthDate)

	days := duration.Hours() / 24.0

	return int(days)
}
