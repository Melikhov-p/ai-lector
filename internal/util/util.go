package util

import (
	"errors"
	"regexp"
	"strings"
)

// isValidEmail проверяет корректность email
func IsValidEmail(email string) error {
	// Простая регулярка для валидации email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	if emailRegex.MatchString(email) {
		return nil
	} else {
		return errors.New("invalid email")
	}
}

// NormalizePhone приводит телефон к виду 79998887766
func NormalizePhone(phone string) (string, error) {
	// убираем всё, кроме цифр
	re := regexp.MustCompile(`\D`)
	digits := re.ReplaceAllString(phone, "")

	// если номер начинается с "8", заменим на "7"
	if strings.HasPrefix(digits, "8") && len(digits) == 11 {
		digits = "7" + digits[1:]
	}

	// если номер начинается с "7" и длина 11 — оставляем как есть
	// если начинается с "9" и длина 10 — добавим "7" в начало
	if strings.HasPrefix(digits, "9") && len(digits) == 10 {
		digits = "7" + digits
	}

	if len(digits) > 11 {
		return "", errors.New("phone is too big")
	}

	return digits, nil
}
