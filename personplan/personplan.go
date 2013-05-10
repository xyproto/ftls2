package main

import (
	"time"
	"fmt"
	"strconv"

	"github.com/xyproto/norwegiantime"
)

type Day struct {
	dayoftheweek time.Weekday
	fromHour     int
	uptoHour     int
}

type PersonPlan struct {
	who      string
	workdays []*Day
}

type PeriodPlan struct {
	year int
	fromMonth int
	uptoMonth int
	personPlans []*PersonPlan
}

func NewPersonPlan(who string) *PersonPlan {
	var pp PersonPlan
	pp.who = who
	return &pp
}

func (pp *PersonPlan) AddWorkday(dayoftheweek time.Weekday, fromHour, uptoHour int) {
	newday := &Day{dayoftheweek, fromHour, uptoHour}
	pp.workdays = append(pp.workdays, newday)
}

func (pp *PersonPlan) String() string {
	s := "User: " + pp.who + "\n"
	s += "-----------------------------------------------\n"
	for _, day := range pp.workdays {
		s += "\n"
		s += "\t" + norwegiantime.N2d(day.dayoftheweek, "en") + " (" + norwegiantime.N2d(day.dayoftheweek, "no") + ")\n"
		s += "\tFrom this hour: \t" + strconv.Itoa(day.fromHour) + "\n"
		s += "\tUp to this hour:\t" + strconv.Itoa(day.uptoHour) + "\n"
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

func (pp *PeriodPlan) String() string {
	s := fmt.Sprintf("From %d, month %d\n", pp.year, pp.fromMonth)
	s += fmt.Sprintf("Up to %d, month %d\n", pp.year, pp.uptoMonth)
	for day := 0; day < 7; day++ {
		for _, persplan := range pp.personPlans {
			for _, personday := range persplan.workdays {
				if personday.dayoftheweek == time.Weekday(day) {
					s += persplan.who + " on " + time.Weekday(day).String() + "\n"
					// TODO: Add hours too?
				}
			}
		}
	}
	return s
}

func main() {
	ppAlexander := NewPersonPlan("Alexander")
	ppAlexander.AddWorkday(time.Monday, 8, 15) // monday, from 8, up to 15
	ppAlexander.AddWorkday(time.Wednesday, 12, 17) // wednesday, from 12, up to 17

	fmt.Println(ppAlexander.String())

	ppBob := NewPersonPlan("Bob")
	ppBob.AddWorkday(time.Tuesday, 9, 11) // monday, from 9, up to 11
	ppBob.AddWorkday(time.Thursday, 8, 10) // wednesday, from 8, up to 10

	fmt.Println(ppBob.String())

	periodplan := NewPeriodPlan(2013, 1, 8)
	periodplan.AddPersonPlan(ppAlexander)
	periodplan.AddPersonPlan(ppBob)

	fmt.Println(periodplan.String())
}

