package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		err := errors.New("неверный формат данных: ожидается 3 элемента")
		log.Println(err)
		return 0, "", 0, err
	}

	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		err = errors.New("не удалось преобразовать кол-во шагов в число: " + err.Error())
		log.Println(err)
		return 0, "", 0, err
	}

	activity := strings.TrimSpace(parts[1])
	duration, err := time.ParseDuration(strings.TrimSpace(parts[2]))
	if err != nil {
		err = errors.New("не удалось преобразовать время: " + err.Error())
		log.Println(err)
		return 0, "", 0, err
	}
	if steps <= 0 || duration <= 0 {
		err := errors.New("некорректные входные параметры для бега")
		log.Println(err)
		return 0, "", 0, err
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	distanceInMeters := float64(steps) * stepLength
	distanceInKm := distanceInMeters / mInKm

	return distanceInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	hours := duration.Hours()
	meanSpeed := dist / hours

	return meanSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	var dist float64 = distance(steps, height)
	var speed float64 = meanSpeed(steps, height, duration)
	var calories float64

	switch activityType {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		err = errors.New("неизвестный тип тренировки")
		log.Println(err)
		return "", err
	}

	if err != nil {
		log.Println(err)
		return "", err
	}

	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityType,
		duration.Hours(),
		dist,
		speed,
		calories)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || duration <= 0 {
		err := errors.New("некорректные входные параметры для бега")
		log.Println(err)
		return 0, err
	}
	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		err := errors.New("некорректные входные параметры для ходьбы")
		log.Println(err)
		return 0, err
	}
	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes) / minInH
	calories *= walkingCaloriesCoefficient

	return calories, nil
}
