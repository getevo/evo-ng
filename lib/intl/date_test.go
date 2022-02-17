package intl_test

import (
	"github.com/getevo/evo-ng/lib/intl"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	assert.Equal(t, "2006/06/08", intl.Date("2006/06/08").Time().Format("2006/01/02"))
	assert.Equal(t, "2006/06/08", intl.Date("2006", "06", "08").Time().Format("2006/01/02"))
	assert.Equal(t, "2006/06/08", intl.Date(2006, 6, 8).Time().Format("2006/01/02"))

	assert.Equal(t, "2006/06/08 10:45:30", intl.Date("2006/06/08 10:45:30").Time().Format("2006/01/02 15:04:05"))
	assert.Equal(t, "2006/06/08 10:45:30", intl.Date(2006, 6, 8, 10, 45, 30).Time().Format("2006/01/02 15:04:05"))
	assert.Equal(t, int64(1149763530000000020), intl.Date(2006, 6, 8, 10, 45, 30, 20).Time().UnixNano())

	assert.Equal(t, intl.Date().Time(), time.Now())

	var l, _ = time.LoadLocation("Asia/Tehran")

	assert.Equal(t, "2006/06/08 07:15:30", intl.Date(2006, 6, 8, 10, 45, 30, 20, l).Time().UTC().Format("2006/01/02 15:04:05"))
	assert.Equal(t, "2006/06/08 07:15:30", intl.Date(2006, 6, 8, 10, 45, 30, 20, *l).Time().UTC().Format("2006/01/02 15:04:05"))

	//t.Error() // to indicate test failed
	intl.AddLocale("en-GB", "en-US", "it-IT", "fa-IR", "es-ES", "ru-RU")
	var locale = intl.GuessLocale("it")

	assert.Equal(t, locale.String(), "it-IT")

	assert.Equal(t, "Giovedì 06 Giu, 2006", intl.Date(2006, 6, 8).Format("Monday 01 Jan, 2006", "it-IT"))
	assert.Equal(t, "Giovedì 06 Giu, 2006", intl.Date(2006, 6, 8).Format("Monday 01 Jan, 2006", "it"))
	assert.Equal(t, "Giovedì 06 Giu, 2006", intl.Date(2006, 6, 8).Format("Monday 01 Jan, 2006", locale))
	assert.Equal(t, "Четверг 06 Июн, 2006", intl.Date(2006, 6, 8).Format("Monday 01 Jan, 2006", "ru"))

	var base = intl.Date(2006, 6, 8, 10, 45, 30, 20)

	assert.Equal(t, "2006/06/04", base.SetDay(4).Format("2006/01/02"))
	assert.Equal(t, "2022/06/08", base.SetYear(2022).Format("2006/01/02"))
	assert.Equal(t, "2006/02/08", base.SetMonth(2).Format("2006/01/02"))
	assert.Equal(t, "2006/06/08 01:45:30", base.SetHour(1).Format("2006/01/02 15:04:05"))
	assert.Equal(t, "2006/06/08 10:02:30", base.SetMinute(2).Format("2006/01/02 15:04:05"))
	assert.Equal(t, "2006/06/08 10:45:03", base.SetSecond(3).Format("2006/01/02 15:04:05"))

	var diff time.Duration
	var cases = map[string][]string{
		"2007/06/08": []string{"next year", "1 year after", "+1 year"},
		"2006/07/08": []string{"next month", "1 month after", "+1 month"},
		"2006/06/15": []string{"next week", "1 week after", "+1 week"},
		"2006/06/09": []string{"next day", "1 day after", "+1 day", "tomorrow"},

		"2005/06/08": []string{"prev year", "previous year", "past year", "1 year before", "-1 year"},
		"2006/05/08": []string{"prev month", "previous month", "past month", "1 month before", "-1 month"},
		"2006/06/01": []string{"prev week", "previous week", "past week", "1 week before", "-1 week"},
		"2006/06/07": []string{"prev day", "previous day", "past day", "1 day before", "-1 day", "yesterday"},
	}
	for expected, list := range cases {
		for _, item := range list {
			calculated, err := base.Calculate(item)
			assert.NoError(t, err)
			assert.Equal(t, expected, calculated.Format("2006/01/02"), item)

			diff, err = base.DiffExpr(item)
			assert.Equal(t, expected, base.Add(diff).Format("2006/01/02"), item)
		}
	}

	cases = map[string][]string{
		"2007/06/08 00:00:00": []string{"next year midnight", "1 year after midnight", "+1 year midnight"},
		"2006/07/08 00:00:00": []string{"next month midnight", "1 month after midnight", "+1 month midnight"},
		"2006/06/15 00:00:00": []string{"next week midnight", "1 week after midnight", "+1 week midnight"},
		"2006/06/09 00:00:00": []string{"next day midnight", "1 day after midnight", "+1 day midnight", "tomorrow midnight"},

		"2005/06/08 00:00:00": []string{"prev year midnight", "previous year midnight", "past year midnight", "1 year before midnight", "-1 year midnight"},
		"2006/05/08 00:00:00": []string{"prev month midnight", "previous month midnight", "past month midnight", "1 month before midnight", "-1 month midnight"},
		"2006/06/01 00:00:00": []string{"prev week midnight", "previous week midnight", "past week midnight", "1 week before midnight", "-1 week midnight"},
		"2006/06/07 00:00:00": []string{"prev day midnight", "previous day midnight", "past day midnight", "1 day before midnight", "-1 day midnight", "yesterday midnight"},
	}

	for expected, list := range cases {
		for _, item := range list {
			//Test Calculate
			calculated, err := base.Calculate(item)
			assert.NoError(t, err)
			assert.Equal(t, expected, calculated.Format("2006/01/02 15:04:05"), item)

			//Test DiffExpr
			diff, err = base.DiffExpr(item)
			assert.NoError(t, err)
			assert.Equal(t, expected, base.Add(diff).Format("2006/01/02 15:04:05"), item)
		}
	}

	assert.Equal(t, true, base.After("2006/06/09") == false)
	assert.Equal(t, true, base.Before("2006/06/07") == false)

}
