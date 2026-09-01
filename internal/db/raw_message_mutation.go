package db

import "sync"

var (
	rawMessageMutationLock     sync.RWMutex
	rawMessageMutationListener func(reason string)
)

// SetRawMessageMutationListener registers the callback invoked after stored
// raw-message rows are edited or deleted outside the normal learning path.
func SetRawMessageMutationListener(fn func(reason string)) {
	rawMessageMutationLock.Lock()
	defer rawMessageMutationLock.Unlock()

	rawMessageMutationListener = fn
}

// NotifyRawMessageMutation invokes the registered raw-message mutation
// listener, if any.
func NotifyRawMessageMutation(reason string) {
	rawMessageMutationLock.RLock()
	fn := rawMessageMutationListener
	rawMessageMutationLock.RUnlock()

	if fn != nil {
		fn(reason)
	}
}
