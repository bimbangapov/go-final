package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	ErrorRepeatRules      = errors.New("incorrect repeat rules")
	ErrorInvalidLetter    = errors.New("invalid letter in rule")
	ErrorInvalidDaysCount = errors.New("invalid count of days in rule")
	dayOfWeek             = map[int]bool{
		1: true,
		2: true,
		3: true,
		4: true,
		5: true,
		6: true,
		7: true,
	}
	dayOfMonth = map[int]bool{
		1:  true,
		2:  true,
		3:  true,
		4:  true,
		5:  true,
		6:  true,
		7:  true,
		8:  true,
		9:  true,
		10: true,
		11: true,
		12: true,
		13: true,
		14: true,
		15: true,
		16: true,
		17: true,
		18: true,
		19: true,
		20: true,
		21: true,
		22: true,
		23: true,
		24: true,
		25: true,
		26: true,
		27: true,
		28: true,
		29: true,
		30: true,
		31: true,
		-1: true,
		-2: true,
	}
	lastMonthDay = map[int]int{
		1:  31,
		2:  28,
		22: 29,
		3:  31,
		4:  30,
		5:  31,
		6:  30,
		7:  31,
		8:  31,
		9:  30,
		10: 31,
		11: 30,
		12: 31,
	}
)

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	dateStart, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", err
	}

	var (
		nextAvailableDate time.Time
	)

	repeatSlice, err := commonCheckRules(repeat)
	if err != nil {
		return "", err
	}

	switch repeatSlice[0] {
	case "y", "d":
		nextAvailableDate, err = basicCheckRules(repeatSlice, dateStart, now)
		if err != nil {
			return "", err
		}
		return nextAvailableDate.Format("20060102"), nil
	case "w":
		nextAvailableDate, err = checkRulesForWeek(repeatSlice, dateStart, now)
		if err != nil {
			return "", err
		}
		return nextAvailableDate.Format("20060102"), nil
	case "m":
		if len(repeatSlice) == 2 {
			nextAvailableDate, err = checkRulesForMonthBasic(repeatSlice, dateStart, now)
			if err != nil {
				return "", ErrorRepeatRules
			}
			return nextAvailableDate.Format("20060102"), nil
		}
		if len(repeatSlice) == 3 {
			var (
				daysSlc        = strings.Split(repeatSlice[1], ",")
				monthSlc       = strings.Split(repeatSlice[2], ",")
				candidateDates []time.Time
				candidateDate  time.Time
			)

			for _, i := range monthSlc {
				month, err := strconv.Atoi(i)
				if err != nil {
					return "", err
				}
				candidateDate, err = checkRulesForMonthExtended(daysSlc, dateStart, now, time.Month(month))
				if err != nil {
					return "", err
				}
				candidateDates = append(candidateDates, candidateDate)
			}
			candidateDate = candidateDates[0]
			for _, item := range candidateDates {
				if item.Before(candidateDate) {
					candidateDate = item
				}
			}

			return candidateDate.Format("20060102"), nil

		}
	}
	return "", ErrorRepeatRules
}
func checkRulesForMonthExtended(repeatSlice []string, dateStart, now time.Time, month time.Month) (time.Time, error) {
	if dateStart.Before(now) {
		dateStart = now
	}
	var (
		tmpYear                = dateStart.Year()
		tmpMonth               = dateStart.Month()
		tmpDay                 = dateStart.Day()
		candidateDatesInternal []time.Time
		candidateDateInternal  time.Time
	)

	for _, i := range repeatSlice {
		day, err := strconv.Atoi(i)
		if err != nil {
			return time.Time{}, err
		}
		tmpDay = day
		tmpMonth = month
		targetDate := time.Date(tmpYear, tmpMonth, tmpDay, 0, 0, 0, 0, time.UTC)
		if DateBefore(targetDate, dateStart) {
			targetDate = targetDate.AddDate(1, 0, 0)
		}
		candidateDatesInternal = append(candidateDatesInternal, targetDate)
	}

	candidateDateInternal = candidateDatesInternal[0]
	for _, item := range candidateDatesInternal {
		if DateBefore(item, dateStart) {
			candidateDateInternal = item
		}
	}

	return candidateDateInternal, nil
}

func checkRulesForMonthBasic(repeatSlice []string, dateStart, now time.Time) (time.Time, error) {
	//1 получаем массив необходимых дней
	dayOfMonthSlc := strings.Split(repeatSlice[1], ",")

	if DateBefore(dateStart, now) {
		dateStart = now
	}
	//2 если дни не заполнены то ошибка
	if len(dayOfMonthSlc) == 0 {
		return time.Time{}, ErrorRepeatRules
	}

	//1 для каждого дня в массиве находим дату и сохраняем ее в переменную
	//В цикле каждую найденную дату сравниваем с предыдущей и оставляем в переменной только самую наиближайшую
	var (
		nextAvailableDay, earliestCandidate time.Time
		candidateDate                       []time.Time
	)
	for _, i := range dayOfMonthSlc {
		needDay, err := strconv.Atoi(i)
		if err != nil {
			return time.Time{}, err
		}

		if !dayOfMonth[needDay] {
			return time.Time{}, ErrorRepeatRules
		}

		tmpMonth := dateStart.Month()
		tmpYear := dateStart.Year()
		tmpDay := dateStart.Day()

		// проверяем на последний и предпоследний день месяца
		if needDay == -1 {
			if dateStart.Day() == lastMonthDay[int(dateStart.Month())] {
				tmpMonth += +1
			}

			//учитываем вискосный год
			if isleapYear(tmpYear) && tmpMonth == time.February {
				tmpDay = lastMonthDay[22]
			} else {
				tmpDay = lastMonthDay[int(tmpMonth)]
			}

			if dateStart.Month() == time.December {
				tmpYear += 1
			}
			nextDate := time.Date(tmpYear, tmpMonth, tmpDay, 0, 0, 0, 0, time.UTC)
			nextAvailableDay = getEarliestDate(nextAvailableDay, nextDate)
			candidateDate = append(candidateDate, nextAvailableDay)
			continue
		}
		if needDay == -2 {
			if dateStart.Day() == lastMonthDay[int(dateStart.Month())]-1 {
				tmpMonth += +1
			}
			if isleapYear(tmpYear) && tmpMonth == time.February {
				tmpDay = lastMonthDay[22] - 1
			} else {
				tmpDay = lastMonthDay[int(tmpMonth)] - 1
			}

			if dateStart.Month() == time.December {
				tmpYear += 1
			}
			nextDate := time.Date(tmpYear, tmpMonth, tmpDay, 0, 0, 0, 0, time.UTC)
			nextAvailableDay = getEarliestDate(nextAvailableDay, nextDate)
			candidateDate = append(candidateDate, nextAvailableDay)
			continue
		}

		//находим след месяц у которого есть нужный день
		if needDay <= dateStart.Day() || needDay >= lastMonthDay[int(dateStart.Month())] {
			tmpMonth += +1
		}

		//увеличиваем год если декабрь
		if dateStart.Month() == time.December {
			tmpYear += 1
		}

		nextDate := time.Date(tmpYear, tmpMonth, needDay, 0, 0, 0, 0, time.UTC)
		nextAvailableDay = getEarliestDate(nextAvailableDay, nextDate)
		//кладем найденную дату в массив кандидатов
		candidateDate = append(candidateDate, nextAvailableDay)
	}

	// находим в массиве кандидатов самую раннюю дату
	earliestCandidate = candidateDate[0]
	for _, item := range candidateDate {
		if DateBefore(item, earliestCandidate) {
			earliestCandidate = item
		}
	}
	return earliestCandidate, nil
}

func isleapYear(year int) bool {
	return year%400 == 0 || (year%4 == 0 && year%10 != 0)
}

func getEarliestDate(nextAvailableDay, dateStart time.Time) time.Time {
	if nextAvailableDay.IsZero() {
		nextAvailableDay = dateStart
		return nextAvailableDay
	}
	if DateBefore(dateStart, nextAvailableDay) {
		nextAvailableDay = dateStart
	}
	return nextAvailableDay
}

func dateOfDayOfTheWeek(needDayWeek int, date, now time.Time) (time.Time, error) {
	if DateBefore(date, now) {
		date = now
	}

	dayFrom := int(date.Weekday())
	if dayFrom == 0 {
		dayFrom = 7
	}
	daysAhead := needDayWeek - dayFrom
	if daysAhead <= 0 {
		daysAhead += 7
	}
	for {
		date = date.AddDate(0, 0, daysAhead)
		if afterNow(date, now) {
			break
		}
	}
	return date, nil
}

func afterNow(date, now time.Time) bool {
	if DateAfter(date, now) {
		return true
	}
	return false
}

func commonCheckRules(repeat string) ([]string, error) {
	repeatSlice := strings.Split(repeat, " ")
	if len(repeatSlice) == 0 {
		return nil, ErrorRepeatRules
	}

	if repeatSlice[0] != "d" && repeatSlice[0] != "y" && repeatSlice[0] != "m" && repeatSlice[0] != "w" {
		return nil, ErrorInvalidLetter
	}

	if repeatSlice[0] == "y" && len(repeatSlice) > 1 {
		return nil, ErrorRepeatRules
	}

	if (repeatSlice[0] == "d" || repeatSlice[0] == "w" || repeatSlice[0] == "m") && len(repeatSlice) == 1 {
		return nil, ErrorRepeatRules
	}

	return repeatSlice, nil
}

func basicCheckRules(repeatSlice []string, date, now time.Time) (time.Time, error) {
	var countDays int

	if repeatSlice[0] == "y" {
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
		return date, nil
	} else {
		value, err := strconv.Atoi(repeatSlice[1])
		if err != nil {
			return time.Time{}, err
		}
		countDays = value
	}
	if countDays < 1 || countDays > 400 {
		return time.Time{}, ErrorInvalidDaysCount
	}

	for {
		date = date.AddDate(0, 0, countDays)
		if afterNow(date, now) {
			break
		}
	}

	return date, nil
}

func checkRulesForWeek(repeatSlice []string, dateStart, now time.Time) (time.Time, error) {
	inputSliceWeek := strings.Split(repeatSlice[1], ",")
	if len(inputSliceWeek) > 7 {
		return time.Time{}, ErrorRepeatRules
	}
	var nextAvailableDay time.Time

	for item := range inputSliceWeek {
		needDayWeek, err := strconv.Atoi(inputSliceWeek[item])
		if err != nil {
			return time.Time{}, err
		}

		if !dayOfWeek[needDayWeek] {
			return time.Time{}, ErrorRepeatRules
		}

		if dayOfWeek[needDayWeek] {
			if nextAvailableDay.IsZero() {
				nextAvailableDay, err = dateOfDayOfTheWeek(needDayWeek, dateStart, now)
			} else {
				tempraryAvailableDay, err := dateOfDayOfTheWeek(needDayWeek, dateStart, now)
				if err != nil {
					return time.Time{}, err
				}
				if DateBefore(tempraryAvailableDay, nextAvailableDay) {
					nextAvailableDay = tempraryAvailableDay
				}
			}
		}
	}
	return nextAvailableDay, nil
}
