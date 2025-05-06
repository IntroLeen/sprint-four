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

		log.Println("invalid data format: expected 3 elements")
		return 0, "", 0, errors.New("invalid data format: expected 3 elements")
	}

	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {

		return 0, "", 0, fmt.Errorf("failed to convert number of steps to number: %w", err)
	}

	activity := strings.TrimSpace(parts[1])
	duration, err := time.ParseDuration(strings.TrimSpace(parts[2]))
	if err != nil {

		return 0, "", 0, fmt.Errorf("failed to convert time: %w", err)
	}
	if steps <= 0 || duration <= 0 {

		log.Println("Incorrect input parameters for running")
		return 0, "", 0, errors.New("incorrect input parameters for running")
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

		log.Println("неизвестный тип тренировкиg")
		return "", errors.New("неизвестный тип тренировкиg")
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
	if steps <= 0 || weight <= 0 || duration <= 0 || height <= 0 {

		log.Println("Incorrect input parameters for running")
		return 0, errors.New("incorrect input parameters for running")
	}
	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {

		log.Println("Incorrect input parameters for walking")
		return 0, errors.New("incorrect input parameters for walking")
	}
	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes) / minInH
	calories *= walkingCaloriesCoefficient

	return calories, nil
}
