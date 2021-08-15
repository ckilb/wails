package runtime

import (
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/logger"
	"sync"
)

// eventListener holds a callback function which is invoked when
// the event listened for is emitted. It has a counter which indicates
// how the total number of events it is interested in. A value of zero
// means it does not expire (default).
type eventListener struct {
	callback func(...interface{}) // Function to call with emitted event data
	counter  int                  // The number of times this callback may be called. -1 = infinite
	delete   bool                 // Flag to indicate that this listener should be deleted
}

// Events handles eventing
type Events struct {
	log      *logger.Logger
	frontend frontend.Frontend

	// Go event listeners
	listeners  map[string][]*eventListener
	notifyLock sync.RWMutex
}

func (e *Events) Notify(name string, data ...interface{}) {
	e.notify(name, data...)
}

func (e *Events) On(eventName string, callback func(...interface{})) {
	e.registerListener(eventName, callback, -1)
}

func (e *Events) OnMultiple(eventName string, callback func(...interface{}), counter int) {
	e.registerListener(eventName, callback, counter)
}

func (e *Events) Once(eventName string, callback func(...interface{})) {
	e.registerListener(eventName, callback, 1)
}

func (e *Events) Emit(eventName string, data ...interface{}) {
	e.notify(eventName, data...)
	e.frontend.Notify(eventName, data...)
}

func (e *Events) Off(eventName string) {
	e.unRegisterListener(eventName)
}

// NewEvents creates a new log subsystem
func NewEvents(log *logger.Logger) *Events {
	result := &Events{
		log:       log,
		listeners: make(map[string][]*eventListener),
	}
	return result
}

// registerListener provides a means of subscribing to events of type "eventName"
func (e *Events) registerListener(eventName string, callback func(...interface{}), counter int) {
	// Create new eventListener
	thisListener := &eventListener{
		callback: callback,
		counter:  counter,
		delete:   false,
	}
	e.notifyLock.Lock()
	// Append the new listener to the listeners slice
	e.listeners[eventName] = append(e.listeners[eventName], thisListener)
	e.notifyLock.Unlock()
}

// unRegisterListener provides a means of unsubscribing to events of type "eventName"
func (e *Events) unRegisterListener(eventName string) {
	e.notifyLock.Lock()
	// Clear the listeners
	delete(e.listeners, eventName)
	e.notifyLock.Unlock()
}

// Notify for the given event name
func (e *Events) notify(eventName string, data ...interface{}) {

	// Get list of event listeners
	listeners := e.listeners[eventName]
	if listeners == nil {
		e.log.Trace("No listeners for event '%s'", eventName)
		return
	}

	// Lock the listeners
	e.notifyLock.Lock()

	// We have a dirty flag to indicate that there are items to delete
	itemsToDelete := false

	// Callback in goroutine
	for _, listener := range listeners {
		if listener.counter > 0 {
			listener.counter--
		}
		go listener.callback(data...)

		if listener.counter == 0 {
			listener.delete = true
			itemsToDelete = true
		}
	}

	// Do we have items to delete?
	if itemsToDelete == true {

		// Create a new Listeners slice
		var newListeners []*eventListener

		// Iterate over current listeners
		for _, listener := range listeners {
			// If we aren't deleting the listener, add it to the new list
			if !listener.delete {
				newListeners = append(newListeners, listener)
			}
		}

		// Save new listeners or remove entry
		if len(newListeners) > 0 {
			e.listeners[eventName] = newListeners
		} else {
			delete(e.listeners, eventName)
		}
	}

	// Unlock
	e.notifyLock.Unlock()
}

func (e *Events) SetFrontend(appFrontend frontend.Frontend) {
	e.frontend = appFrontend
}