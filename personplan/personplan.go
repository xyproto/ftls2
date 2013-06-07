package main

import (
	"fmt"
	"time"

	"github.com/xyproto/moskus"
)

const (
	FIRSTHOUR = 8
	LASTHOUR  = 22
)

// Info about one thing that can happen during an hour
type HourInfoPerson struct {
	who   string
	when  time.Time
	where string
}

// Info about everything that happens during an hour, per person
type HourInfo []*HourInfoPerson

// A plan is a collecion of plans for just a few months at a time
type Plans struct {
	all []*PeriodPlan
}

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
	cal, err := moskus.NewCalendar("nb_NO", true)
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
		for hour := FIRSTHOUR; hour <= LASTHOUR; hour++ {
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

func (pp *PeriodPlan) ViewHour(t time.Time) string {
	s := ""

	//PeriodPlan
	//year        int
	//fromMonth   int
	//uptoMonth   int
	//personPlans []*PersonPlan

	// if not the right year
	if t.Year() != pp.year {
		return ""
	}

	// if not within the month range
	if !((t.Month() >= time.Month(pp.fromMonth)) && (t.Month() < time.Month(pp.uptoMonth))) {
		return ""
	}

	for _, persplan := range pp.personPlans {
		// PersonPlan
		//who      string
		//workdays []*WorkDayAndLocation
		for _, wd := range persplan.workdays {
			// WorkDayAndLocation
			// dayoftheweek time.Weekday
			// fromHour     int
			// uptoHour     int
			// location     string

			//fmt.Printf("persplan %d, workday %d\n", i1, i2)

			// If not the right day of the week
			if wd.dayoftheweek != t.Weekday() {
				//fmt.Printf("Wrong day of the week! (%v and %v)\n", wd.dayoftheweek, t.Weekday())
				continue
			}

			// If not within the hour range
			if !((t.Hour() >= wd.fromHour) && (t.Hour() < wd.uptoHour)) {
				//fmt.Printf("Wrong hour range! (%v is not between %v and %v)\n", t.Hour(), wd.fromHour, wd.uptoHour)
				continue
			}

			// Found!
			s += fmt.Sprintf("%s %s at %s, %v at hour %v\n", t.String()[:10], persplan.who, wd.location, t.Weekday(), t.Hour())
		}
	}

	return s
}

func (pp *PeriodPlan) ViewDay(date time.Time) string {
	var t time.Time
	var hourString string
	s := ""
	for hour := FIRSTHOUR; hour <= LASTHOUR; hour++ {
		//fmt.Printf("hour: %d\n", hour)
		t = time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, time.UTC)
		hourString = pp.ViewHour(t)
		if hourString != "" {
			s += hourString + "\n"
		}
	}
	return s
}

// Make new plans, which is a collection of PeriodPlans
func NewPlans() *Plans {
	var plans Plans
	plans.all = make([]*PeriodPlan, 0)
	return &plans
}

// Add a PeriodPlan to the collection of plans
func (plans *Plans) AddPeriodPlan(pp *PeriodPlan) {
	plans.all = append(plans.all, pp)
}

// TODO: Create a function just like this that returns an HourInfo type instead of strings
func (plans *Plans) HourInfo(date time.Time) {
	fmt.Printf("What's up at %s?\n", date.String())
	s := ""
	for _, pp := range plans.all {
		s += pp.ViewHour(date)
	}
	if s == "" {
		fmt.Println("Nothing!")
	} else {
		fmt.Println(s)
	}
}

func main() {
	ppAlexander := NewPersonPlan("Alexander")
	ppAlexander.AddWorkday(time.Monday, 8, 15, "KNH")     // monday, from 8, up to 15
	ppAlexander.AddWorkday(time.Wednesday, 12, 17, "KOH") // wednesday, from 12, up to 17

	fmt.Println(ppAlexander.String())

	ppBob := NewPersonPlan("Bob")
	ppBob.AddWorkday(time.Monday, 9, 11, "KOH")   // monday, from 9, up to 11
	ppBob.AddWorkday(time.Thursday, 8, 10, "KNH") // wednesday, from 8, up to 10

	fmt.Println(ppBob.String())

	periodplan := NewPeriodPlan(2013, 1, 8)
	periodplan.AddPersonPlan(ppAlexander)
	periodplan.AddPersonPlan(ppBob)

	fmt.Println(periodplan.String())

	allPlans := NewPlans()
	allPlans.AddPeriodPlan(periodplan)

	fmt.Println("Info for all plans:")
	date := time.Date(2013, 3, 4, 10, 32, 0, 0, time.UTC)

	for i, pp := range allPlans.all {
		fmt.Printf("Plan %d\n", i)
		fmt.Println(pp.ViewDay(date))
	}

	allPlans.HourInfo(date)

	allPlans.HourInfo(time.Date(2013, 3, 7, 9, 14, 0, 0, time.UTC))
}
