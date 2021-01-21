package rules

import (
	"context"
	"fmt"
	"log"
	"time"
)

const minimumWatchTime = 1 * time.Minute

// RuleWatchCallback is a callback function which will get called when some
// modified rules have been found by the RuleWatcher.
type RuleWatchCallback func(rules []*Rule)

// GetChangesSinceClient represents things that implement the GetChangesSince
// method, which would normally be the Client.
type GetChangesSinceClient interface {
	GetChangesSince(timestamp time.Time, eventType *RuleEventType, ruleType *RuleType) ([]*Rule, error)
}

var _ GetChangesSinceClient = &Client{}

// RuleWatcherArgs represents the arguments needed when running RuleWatcher.
type RuleWatcherArgs struct {
	Client    GetChangesSinceClient
	HowOften  time.Duration
	Callback  RuleWatchCallback
	RuleType  *RuleType
	EventType *RuleEventType

	// Used for the test
	allowShortTime bool
}

// RuleWatcher will watch the API for rule changes and notify a callback with
// modified rules.
type RuleWatcher struct {
	args   RuleWatcherArgs
	ctx    context.Context
	cancel context.CancelFunc
	fin    chan struct{}
}

// NewRuleWatcher builds a RuleWatcher which will use the Client provided in
// the args to poll the API every HowOften, and will call the provided Callback
// with any new or updated rules. The RuleType and EventType are optional, and
// if set will be sent to the GetChangesSince method of the provided Client.
//
// The HowOften value must be 1 minute or larger, otherwise an error is
// returned.
func NewRuleWatcher(args RuleWatcherArgs) (*RuleWatcher, error) {
	if !args.allowShortTime && args.HowOften < minimumWatchTime {
		return nil, fmt.Errorf("provided duration %s is too small, 1 minute minimum watching time", args.HowOften)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &RuleWatcher{
		args:   args,
		ctx:    ctx,
		cancel: cancel,
		fin:    make(chan struct{}),
	}, nil
}

// Watch will start the rule watcher watching for rule changes.
func (w *RuleWatcher) Watch() {
	go w.watcherLoop()
}

// Shutdown will shutdown the watcher.
func (w *RuleWatcher) Shutdown(ctx context.Context) (err error) {
	w.cancel()

	select {
	case <-ctx.Done():
		err = ctx.Err()
	case <-w.fin:
	}

	return
}

func (w *RuleWatcher) watcherLoop() {
	var start time.Time
	for {
		start = time.Now()
		select {
		case <-time.After(w.args.HowOften):
			rules, err := w.args.Client.GetChangesSince(start, w.args.EventType, w.args.RuleType)
			if err != nil {
				log.Printf("error getting changed rules: %s\n", err)
			} else {
				if rules != nil {
					w.args.Callback(rules)
				}
			}
		case <-w.ctx.Done():
			close(w.fin)
			return
		}
	}
}
