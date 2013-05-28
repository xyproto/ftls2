package main

import (
	"fmt"
	"time"

	"github.com/xyproto/norwegiantime"
)

type WorkDayAndLocation struct {
	dayoftheweek time.Weekday
	fromHour     int
	uptoHour     int
	location     string
}

type PersonPlan struct {
	who      string
	workdays []*WorkDayAndLocation
}

type PeriodPlan struct {
	year        int
	fromMonth   int
	uptoMonth   int
	personPlans []*PersonPlan
}

func NewPersonPlan(who string) *PersonPlan {
	var pp PersonPlan
	pp.who = who
	return &pp
}

func (pp *PersonPlan) AddWorkday(dayoftheweek time.Weekday, fromHour, uptoHour int, location string) {
	newday := &WorkDayAndLocation{dayoftheweek, fromHour, uptoHour, location}
	pp.workdays = append(pp.workdays, newday)
}

func (pp *PersonPlan) String() string {
	cal, err := norwegiantime.NewCalendar("nb_NO", true)
	if err != nil {
		panic("No calendar available for nb_NO")
	}
	s := "User: " + pp.who + "\n"
	s += "-----------------------------------------------\n"
	for _, day := range pp.workdays {
		s += "\n"
		s += "\t" + day.dayoftheweek.String() + " (" + cal.DayName(day.dayoftheweek) + ")\n"
		s += fmt.Sprintf("\tFrom this hour: \t%d\n", day.fromHour)
		s += fmt.Sprintf("\tUp to this hour:\t%d\n", day.uptoHour)
		s += fmt.Sprintf("\tAt this location:\t%s\n", day.location)
	}
	return s
}

func NewPeriodPlan(year, fromMonth, uptoMonth int) *PeriodPlan {
	var pps []*PersonPlan
	return &PeriodPlan{year, fromMonth, uptoMonth, pps}
}

func (pp *PeriodPlan) AddPersonPlan(persplan *PersonPlan) {
	pp.personPlans = append(pp.personPlans, persplan)
}

func (pp *PeriodPlan) ForAllWeekdays(fn func(string, time.Weekday, int, string) string) string {
	s := ""
	for day := 0; day < 7; day++ {
		for hour := 8; hour < 21; hour++ {
			for _, persplan := range pp.personPlans {
				for _, personday := range persplan.workdays {
					if personday.dayoftheweek == time.Weekday(day) {
						if (hour >= personday.fromHour) && (hour < personday.uptoHour) {
							s += fn(persplan.who, time.Weekday(day), hour, personday.location)
						}
					}
				}
			}
		}
	}
	return s
}

func infoline(who string, weekday time.Weekday, hour int, location string) string {
	return fmt.Sprintf("%s on %s hour that starts at %d at %s\n", who, weekday, hour, location)
}

func (pp *PeriodPlan) String() string {
	s := fmt.Sprintf("From %d, month %d\n", pp.year, pp.fromMonth)
	s += fmt.Sprintf("Up to %d, month %d\n", pp.year, pp.uptoMonth)
	s += pp.ForAllWeekdays(infoline)
	return s
}

// TODO: Broken, fix
// TODO: Return an HourInfo struct or something instead
func (pp *PeriodPlan) HourInfo(t time.Time) string {
	s := ""
	for dayoftheweek := 0; dayoftheweek < 7; dayoftheweek++ {
		for _, persplan := range pp.personPlans {
			for _, personday := range persplan.workdays {
				// Right day of the week?
				if personday.dayoftheweek == t.Weekday() {
					s += fmt.Sprintf("%s on %s, at the hour that starts at %d, at %s\n", persplan.who, t.Weekday(), t.Hour(), personday.location)
				}
			}
		}
	}
	return s
}

// TODO: Broken, fix
// Find info for a given hour for all given period plans
// Note that if period plans are overlapping, this function will only return info for the first hour it finds!
// TODO: Separate day and day of the week more clearly
func HourInfo(allPlans []*PeriodPlan, t time.Time) string {
	fmt.Printf("year %d, month %d, day of the month %d, hour %d\n", t.Year(), t.Month(), t.Day(), t.Hour())
	for _, pp := range allPlans {
		if (pp.year == t.Year()) && (pp.fromMonth <= int(t.Month())) && (int(t.Month()) < pp.uptoMonth) {
			return pp.HourInfo(t)
		}
	}
	// TODO: Introduce err
	return "No hour info found"
}

func main() {
	ppAlexander := NewPersonPlan("Alexander")
	ppAlexander.AddWorkday(time.Monday, 8, 15, "KNH")     // monday, from 8, up to 15
	ppAlexander.AddWorkday(time.Wednesday, 12, 17, "KOH") // wednesday, from 12, up to 17

	fmt.Println(ppAlexander.String())

	ppBob := NewPersonPlan("Bob")
	ppBob.AddWorkday(time.Tuesday, 9, 11, "KOH")  // monday, from 9, up to 11
	ppBob.AddWorkday(time.Thursday, 8, 10, "KNH") // wednesday, from 8, up to 10

	fmt.Println(ppBob.String())

	periodplan := NewPeriodPlan(2013, 1, 8)
	periodplan.AddPersonPlan(ppAlexander)
	periodplan.AddPersonPlan(ppBob)

	fmt.Println(periodplan.String())

	fmt.Println("Hour info:")
	fmt.Println(periodplan.HourInfo(time.Date(2013, 3, 1, 12, 0, 0, 0, time.UTC)))

	var allPlans []*PeriodPlan
	allPlans = append(allPlans, periodplan)

	fmt.Println("Hour info for all plans:")
	fmt.Println(HourInfo(allPlans, time.Date(2013, 3, 3, 15, 0, 0, 0, time.UTC))) // y m d h
}
