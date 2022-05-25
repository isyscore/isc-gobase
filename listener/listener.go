package listener

type EventListener func(event BaseEvent)

var eventWatcherMaps map[string][]EventListener

type BaseEvent interface {
	Name() string
}

func PublishEvent(event BaseEvent) {
	if eventWatcherMaps == nil {
		return
	}
	if eventWatchers, exist := eventWatcherMaps[event.Name()]; exist {
		for _, eventWatcher := range eventWatchers {
			eventWatcher(event)
		}
	}
}

func AddListener(eventName string, eventListener EventListener) {
	if eventWatcherMaps == nil {
		eventWatcherMaps = map[string][]EventListener{}
	}
	if eventWatchers, exist := eventWatcherMaps[eventName]; exist {
		eventWatchers = append(eventWatchers, eventListener)
		eventWatcherMaps[eventName] = eventWatchers
	} else {
		eventWatchers = []EventListener{}
		eventWatchers = append(eventWatchers, eventListener)
		eventWatcherMaps[eventName] = eventWatchers
	}
}
