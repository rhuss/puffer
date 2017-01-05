package calendar

import "time"

type NextEvents struct {
	TodayEvents          *[]TimedEvent
	TomorrowEvents       *[]TimedEvent
	TomorrowAllDayEvents *[]Event
}

// A calendar event happening today or tomorrow
type TimedEvent struct {
	Start *time.Time
	End   *time.Time

	Event
}

type Event struct {
	Calendar string
	Summary  string
}
