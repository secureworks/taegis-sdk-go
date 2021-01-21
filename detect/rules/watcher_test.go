package rules

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type Change struct {
	rules []*Rule
	err   error
	ts    time.Time
}

type MockWatcherClient struct {
	changes       []*Change
	currentChange int
	lastEventType *RuleEventType
	lastRuleType  *RuleType
	m             *sync.Mutex
}

func (t *MockWatcherClient) GetChangesSince(timestamp time.Time, eventType *RuleEventType, ruleType *RuleType) ([]*Rule, error) {
	t.m.Lock()
	defer t.m.Unlock()

	// This will go through all the changes once, then just always
	// return nil, nil, which means no changes and no errors.
	if t.currentChange >= len(t.changes) {
		return nil, nil
	}

	change := t.changes[t.currentChange]

	t.lastEventType = eventType
	t.lastRuleType = ruleType

	t.currentChange++

	return change.rules, change.err
}

func Test_RuleWatcher(t *testing.T) {
	require := require.New(t)

	r1 := &Rule{
		Name: "Rule #1",
	}
	r2 := &Rule{
		Name: "Rule #2",
	}
	r3 := &Rule{
		Name: "Rule #3",
	}
	r4 := &Rule{
		Name: "Rule #4",
	}

	tc := &MockWatcherClient{
		changes: []*Change{
			{
				rules: nil,
				err:   nil,
			},
			{
				rules: nil,
				err:   fmt.Errorf("failed to get changes"),
			},
			{
				rules: []*Rule{r1, r2},
				err:   nil,
			},
			{
				rules: nil,
				err:   nil,
			},
			{
				rules: []*Rule{r3},
				err:   nil,
			},
			{
				rules: []*Rule{r4},
				err:   nil,
			},
		},
		m: new(sync.Mutex),
	}

	changes := []*Change{}
	cb := func(rules []*Rule) {
		changes = append(changes, &Change{
			rules: rules,
			ts:    time.Now(),
		})
	}

	waitTime := 10 * time.Millisecond
	args := RuleWatcherArgs{
		Client:   tc,
		HowOften: waitTime,
		Callback: cb,
	}
	// First make sure the minimum watch time is checked
	w, err := NewRuleWatcher(args)
	require.Nil(w, "watcher should be nil")
	require.EqualError(err, "provided duration 10ms is too small, 1 minute minimum watching time")

	// Now allow a shorter time
	args.allowShortTime = true
	w, err = NewRuleWatcher(args)
	require.NotNil(w, "watcher should not be nil")
	require.NoError(err)

	w.Watch()

	// Give it a little time
	time.Sleep(100 * time.Millisecond)

	// We need to lock the test client before doing assertions to avoid data races
	tc.m.Lock()

	require.Lenf(changes, 3, "there should be three changes")

	// Check the rule from the changes
	require.Len(changes[0].rules, 2)
	require.Equal("Rule #1", changes[0].rules[0].Name)
	require.Equal("Rule #2", changes[0].rules[1].Name)
	require.Len(changes[1].rules, 1)
	require.Equal("Rule #3", changes[1].rules[0].Name)
	require.Len(changes[2].rules, 1)
	require.Equal("Rule #4", changes[2].rules[0].Name)

	// Make sure the timing seems right
	// The difference between change 2 and change 1 should be at least double the wait time
	require.True(changes[1].ts.Sub(changes[0].ts) > waitTime*2,
		"the time difference between changes 1 and 2 is too small")
	// The difference between change 3 and change 2 should be about the wait time
	require.True(changes[2].ts.Sub(changes[1].ts) > waitTime,
		"the time difference between changes 2 and 3 is too small")

	require.Nil(tc.lastEventType)
	require.Nil(tc.lastRuleType)

	tc.m.Unlock()

	w.Shutdown(context.TODO())

	// Now make sure the event type and rule type is passed through

	// Reset test client
	tc.changes = []*Change{
		{
			rules: nil,
			err:   nil,
		},
	}
	tc.currentChange = 0

	eventType := RuleEventTypeAuth
	args.EventType = &eventType
	ruleType := RuleTypeRegex
	args.RuleType = &ruleType

	w, err = NewRuleWatcher(args)
	require.NotNil(w, "watcher should not be nil")
	require.NoError(err)

	w.Watch()

	time.Sleep(50 * time.Millisecond)

	tc.m.Lock()
	require.Equal(&eventType, tc.lastEventType)
	require.Equal(&ruleType, tc.lastRuleType)
	tc.m.Unlock()

	w.Shutdown(context.TODO())
}
