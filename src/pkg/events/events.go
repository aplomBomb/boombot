package events

// Event is a string representing a
//listener for a particular Discord event
type Event = string

//EventList is a slice containing all
//currently listened to Discord events
//along with other helper utils
type EventList struct {
	events []Event
	length int
}

// New instantiates a new EventList.
func (el EventList) New(events []string) EventList {
	return EventList{
		events: events,
		length: len(events),
	}
}

// Add method is for appending new
//listeners to the active EventList
func (el EventList) Add(event string) {

}

// Remove method is for removing existing
//listeners from the active EventList
func (el EventList) Remove(event string) {

}

//Length simply returns the number of
//event listeners in the current EventList
func (el EventList) Length(eventList EventList) int {
	return len(eventList.events)
}

//List simply lists all the event
//listeners in the current EventList
func (el EventList) List(eventList EventList) []string {
	list := []string{""}
	for _, event := range eventList.events {
		return append(list, event)
	}
	return list
}
