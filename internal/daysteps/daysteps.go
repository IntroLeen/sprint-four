package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// проверяем длину слайса

	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		log.Println("invalid data format: expected 2 elements")
		return 0, 0, errors.New("invalid data format: expected 2 elements")
	}
	//тут я уже отчаялся
	if strings.HasPrefix(data, " ") || strings.HasSuffix(data, " ") {

		log.Println("the data contains spaces at the beginning or end")
		return 0, 0, errors.New("the data contains spaces at the beginning or end")
	}
	stepsStr := strings.TrimSpace(parts[0])
	if parts[0] != stepsStr { // Проверка на наличие пробелов в начале или конце шага

		log.Println("the number of steps contains extra spaces")
		return 0, 0, errors.New("the number of steps contains extra spaces")
	}
	// преобразуем кол-во шагов и проверяем на ошибку
	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {

		return 0, 0, fmt.Errorf("failed to convert number of steps to number: %w", err)
	}

	// проверка шагов
	if steps <= 0 {
		log.Println("the number of steps must be greater than 0")
		return 0, 0, errors.New("the number of steps must be greater than 0")
	}

	// преобразуем время и проверяем на ошибку
	duration, err := time.ParseDuration(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert time:  %w", err)
	}
	if duration <= 0 {

		log.Println("the number of steps must be greater than 0")
		return 0, 0, errors.New("the number of steps must be greater than 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	//проверки
	if err != nil {
		log.Printf("Error while parsing data: %v", err)
		return ""
	}
	if steps <= 0 {
		return ""
	}
	// применяем формулы
	distanceInMeters := float64(steps) * stepLength
	distanceInKm := distanceInMeters / mInKm
	calories, calErr := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if calErr != nil {
		log.Printf("Error in calorie calculation: %v", calErr)
		return ""
	}

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceInKm, calories)

	return result
}
