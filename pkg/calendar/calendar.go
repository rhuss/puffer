package calendar

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gcalendar "google.golang.org/api/calendar/v3"
	"sort"
	"log"
)

// helper for extracting calendar data
type calendarItem struct {
	id   string
	name string
}

type ByStart []TimedEvent

func (a ByStart) Len() int           { return len(a) }
func (a ByStart) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStart) Less(i, j int) bool { return a[i].Start.Unix() < a[j].Start.Unix() }

// FetchToken reaches out to Google OAuth with the clients credentials to get a valid token
// Interactive user action is required here
func FetchToken(jsonKey []byte) (*oauth2.Token, error) {
	config, err := google.ConfigFromJSON(jsonKey, gcalendar.CalendarReadonlyScope)
	if err != nil {
		return nil, err
	}
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, err
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}
	return tok, nil
}

// CloseEvent fetches the events from today (or tomorrow if none are there for today)
func GetNextEvents(oauthToken *oauth2.Token, jsonKey []byte, todayCalendarNames []string, allDayCalendarNames [] string) (*NextEvents, error) {
	ctx := context.Background()

	config, err := google.ConfigFromJSON(jsonKey, gcalendar.CalendarReadonlyScope)
	if err != nil {
		return nil, err
	}
	client := config.Client(ctx, oauthToken)

	srv, err := gcalendar.New(client)
	if err != nil {
		return nil, err
	}

	calendars, err := getCalendarIds(srv, todayCalendarNames)
	if err != nil {
		return nil, err
	}
	log.Printf("Calendars: %v", calendars)

	start, midnight, end := getTimeWindow()

	todayEvents, err := extractEventsForCalendars(srv, calendars, start, midnight)
	if err != nil {
		return nil, err
	}
	var tomorrowEvents *[]TimedEvent
	if len(*todayEvents) == 0 {
		tomorrowEvents, err = extractEventsForCalendars(srv, calendars, start, end)
		if err != nil {
			return nil, err
		}
	}

	allDayCalendars, err := getCalendarIds(srv, allDayCalendarNames)
	if err != nil {
		return nil, err
	}

	tomorrowAllDayEvents, err := extractAllDayEvents(srv, allDayCalendars, midnight, end)
	if err != nil {
		return nil, err
	}

	return &NextEvents{
		TodayEvents:    todayEvents,
		TomorrowEvents: tomorrowEvents,
		TomorrowAllDayEvents: tomorrowAllDayEvents,
	}, nil
}
func extractAllDayEvents(srv *gcalendar.Service, calendars []calendarItem, start time.Time, end time.Time) (*[]Event, error) {
	ret := []Event{}
	for _, cal := range calendars {
		events, err := getEvents(srv, cal, start, end)
		if err != nil {
			return nil, err
		}

		if len(events.Items) > 0 {
			for _, i := range events.Items {
				// If the DateTime is an empty string the Event is an all-day Event.
				// So only Date is available.
				if i.Start.DateTime == "" {
					event := Event{
						Calendar: cal.name,
						Summary:  i.Summary,
					}
					ret = append(ret, event)
				}
			}
		}
	}
	if len(ret) == 0 {
		return nil, nil
	} else {
		return &ret, nil
	}
}

func extractEventsForCalendars(srv *gcalendar.Service, calendars []calendarItem, start time.Time, end time.Time) (*[]TimedEvent, error) {
	ret := []TimedEvent{}
	for _, cal := range calendars {
		events, err := getEvents(srv, cal, start, end)
		if err != nil {
			return nil, err
		}
		e, err := extractTimedEvents(cal, events)
		if err != nil {
			return nil, err
		}
		ret = append(ret, *e...)
	}
	sort.Sort(ByStart(ret))
	return &ret, nil
}

func extractTimedEvents(cal calendarItem, events *gcalendar.Events) (*[]TimedEvent, error) {
	ret := []TimedEvent{}
	if len(events.Items) > 0 {
		for _, i := range events.Items {
			// If the DateTime is an empty string the Event is an all-day Event.
			// So only Date is available.
			if i.Start.DateTime != "" {
				startTime, err := parseRfc339(i.Start.DateTime)
				if err != nil {
					return nil, err
				}
				endTime, err := parseRfc339(i.End.DateTime)
				if err != nil {
					return nil, err
				}
				event := TimedEvent{
					Start:     startTime,
					End:       endTime,
					Event: Event{
						Calendar: cal.name,
						Summary:  i.Summary,
					},
				}
				ret = append(ret, event)
			}
		}
	}
	return &ret, nil
}

func getEvents(srv *gcalendar.Service, cal calendarItem, start time.Time, end time.Time) (*gcalendar.Events, error) {
	events, err := srv.Events.List(cal.id).ShowDeleted(false).
		SingleEvents(true).TimeMin(start.Format(time.RFC3339)).TimeMax(end.Format(time.RFC3339)).OrderBy("startTime").Do()
	return events, err
}

func getCalendarIds(srv *gcalendar.Service, summaries []string) ([]calendarItem, error) {
	cl, err := srv.CalendarList.List().Do()
	if err != nil {
		return nil, err
	}
	ret := []calendarItem{}

	for _, i := range cl.Items {
		for _, j := range summaries {
			if j == i.Summary {
				ret = append(ret, calendarItem{
					id:   i.Id,
					name: i.Summary,
				})
			}
		}
	}
	return ret, nil
}

func parseRfc339(rfc3339 string) (*time.Time, error) {
	t, err := time.Parse(time.RFC3339, rfc3339)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// Get the time window for which to fetch events
func getTimeWindow() (time.Time, time.Time, time.Time) {
	start := time.Now().Add(-2 * time.Hour)
	year, month, day := start.Date()
	lastMidnight := time.Date(year, month, day, 0, 0, 0, 0, start.Location())
	midnight := lastMidnight.Add(time.Hour * 24)
	end := lastMidnight.Add(2 * time.Hour * 24)
	return start, midnight, end
}
