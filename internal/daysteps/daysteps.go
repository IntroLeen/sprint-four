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
		err := errors.New("неверный формат данных: ожидается 2 элемента")
		log.Println(err)
		return 0, 0, err
	}
	//тут я уже отчаялся
	if strings.HasPrefix(data, " ") || strings.HasSuffix(data, " ") {
		err := errors.New("данные содержат пробелы в начале или конце")
		log.Println(err)
		return 0, 0, err
	}
	stepsStr := strings.TrimSpace(parts[0])
	if parts[0] != stepsStr { // Проверка на наличие пробелов в начале или конце шага
		err := errors.New("количество шагов содержит лишние пробелы")
		log.Println(err)
		return 0, 0, err
	}
	// преобразуем кол-во шагов и проверяем на ошибку
	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		err = errors.New("не удалось преобразовать кол-во шагов в число: " + err.Error())
		log.Println(err)
		return 0, 0, err
	}

	// проверка шагов
	if steps <= 0 {
		err := errors.New("количество шагов должно быть больше 0")
		log.Println(err)
		return 0, 0, err
	}

	// преобразуем время и проверяем на ошибку
	duration, err := time.ParseDuration(strings.TrimSpace(parts[1]))
	if err != nil {
		err = errors.New("не удалось преобразовать время: " + err.Error())
		log.Println(err)
		return 0, 0, err
	}
	if duration <= 0 {
		err := errors.New("количество шагов должно быть больше 0")
		log.Println(err)
		return 0, 0, err
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	//проверки
	if err != nil {
		log.Println("Ошибка при парсинге данных:", err)
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
		log.Println("Ошибка при расчете калорий:", calErr)
		return ""
	}

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceInKm, calories)

	return result
}
