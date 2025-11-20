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
	stepLength = 0.65

	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	parts := strings.Split(data, ",")

	if len(parts) != 2 {
		return 0, 0, errors.New("Неверный формат")
	}

	stepsStr := parts[0]
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, errors.New("неверный формат шагов" + err.Error())
	}

	if steps <= 0 {
		return 0, 0, errors.New("шагов должно быть больше 0")
	}

	durationStr := parts[1]
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, errors.New("неверный формат продолжительности" + err.Error())
	}
	if duration <= 0 {
		return 0, 0, errors.New("неверный формат продолжительности")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {

	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		fmt.Println("Ошибка:", err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distanceMeters := float64(steps) * stepLength

	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(
		steps,
		weight,
		height,
		duration,
	)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories,
	)
}
