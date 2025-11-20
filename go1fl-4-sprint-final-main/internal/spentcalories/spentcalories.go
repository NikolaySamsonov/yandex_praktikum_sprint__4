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
		return 0, "", 0, errors.New("Неверный формат")
	}

	stepsStr := parts[0]
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, errors.New("Неверное значение шагов" + err.Error())
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("шагов должно быть больше 0")
	}
	activity := parts[1]

	durationStr := parts[2]
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, errors.New("Неверная продолжительность" + err.Error())
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("Продолжительность должна быть больше 0")
	}
	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient

	distanceMeters := float64(steps) * stepLength

	distanceKm := distanceMeters / mInKm

	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height)

	hours := duration.Hours()
	speed := dist / hours

	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	dist := distance(steps, height)             // км
	speed := meanSpeed(steps, height, duration) // км/ч

	var calories float64

	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}

	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}

	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity,
		duration.Hours(),
		dist,
		speed,
		calories,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 {
		return 0, errors.New("Шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("Масса должна быть больше 0")
	}
	if height <= 0 {
		return 0, errors.New("Рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("Продолжительность должна быть больше 0")
	}

	// 2. Рассчитываем среднюю скорость (км/ч)
	speed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	calories := (weight * speed * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 {
		return 0, errors.New("Шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("Масса должна быть больше 0")
	}
	if height <= 0 {
		return 0, errors.New("Рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("Продолжительность должна быть больше 0")
	}

	speed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	baseCalories := (weight * speed * durationInMinutes) / minInH

	finalCalories := baseCalories * walkingCaloriesCoefficient

	return finalCalories, nil
}
