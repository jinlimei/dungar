package random

import (
	"fmt"
	"log"
	"time"
)

// AdvControlGroup is our grouping for AdvancedControl to handle them
// in an organized way
type AdvControlGroup struct {
	controls map[string]*AdvancedControl
}

// Analytics provides a shortened version of AdvancedControl which
// just pertains/contains the data we've captured for running it.
type Analytics struct {
	ID         string  `json:"id"`
	TrueCount  int     `json:"true_count"`
	FalseCount int     `json:"false_count"`
	CallCount  int     `json:"call_count"`
	Chance     float64 `json:"chance"`
}

// AdvancedControl is our advanced handling of chance-based events
type AdvancedControl struct {
	// ID is the actual unique ID of this control
	ID string

	// MaximumOccurrence determines that, when triggered as true,
	// the RNG can trigger true again.
	// For example: If set to 30m with a 100% chance, it will only ever be
	//              true every 30 minutes
	MaximumOccurrence time.Duration

	// MinimumOccurrence determines that, when triggered, the minimum amount of
	// time before the RNG should be triggering true again.
	// For example: If set at 15m with a 0% chance, the RNG will always be true
	//              at a 15m mark if it hasn't been true within 15m
	MinimumOccurrence time.Duration

	// IncrementIncrease when set to a positive number will slowly increase the
	// chance increment until the event occurs, in which case the chance is reset
	// to its original value.
	// Requires calls to be evaluated (will not evaluate behind the scenes)
	IncrementIncrease float64

	// IncreaseDuration defines how often the increment should increase. Requires
	// calls to the function to be evaluated (will not evaluate behind the scenes)
	IncreaseDuration time.Duration

	// DefaultChance is the base level chance in which the function calculates
	// returning true. Range of [0,1] -- 0.0 being always false and 1.0 being always
	// true. Decimal values to determine % (so 0.5 is 50%)
	DefaultChance float64

	// currentChance is the current value of chance
	// (DefaultChance + seconds * IncrementIncrease)
	currentChance float64

	// Last time the chance returned true. Used in calculating MinimumOccurrence
	// and MaximumOccurrence
	lastTimeTrue time.Time

	// Last time the chance was called & executed. Used in calculating IncrementIncrease,
	// based on IncreaseDuration
	lastTimeCalled time.Time

	// isReady returns whether or not the function has been initialized.
	isReady bool

	// callCount defines the # of calls this function has received in runtime.
	// useful for doing analytics on function calls.
	callCount int

	// trueCount defines the # of trues we run through
	trueCount int

	// falseCount defines the # of falses we run through
	falseCount int
}

// String defines an easy string value of the Analytics
func (ra *Analytics) String() string {
	return fmt.Sprintf(
		"ID=%s,True=%d,False=%d,Calls=%d,Chance=%0.6f",
		ra.ID,
		ra.TrueCount,
		ra.FalseCount,
		ra.CallCount,
		ra.Chance,
	)
}

func (rc *AdvancedControl) initIfNotReady() {
	if !rc.isReady {
		rc.Init()
	}
}

// Init will initialize the advanced control
func (rc *AdvancedControl) Init() {
	rc.currentChance = rc.DefaultChance
	rc.lastTimeCalled = time.Unix(0, 0)
	rc.lastTimeTrue = time.Unix(0, 0)
	rc.isReady = true
}

// ToAnalytics will conver the control into an Analytics struct
func (rc *AdvancedControl) ToAnalytics() Analytics {
	return Analytics{
		ID:         rc.ID,
		TrueCount:  rc.trueCount,
		FalseCount: rc.falseCount,
		CallCount:  rc.callCount,
		Chance:     rc.DefaultChance,
	}
}

// Execute runs all the necessary calculations to make sure our RNG is happy.
func (rc *AdvancedControl) Execute() bool {
	rc.initIfNotReady()

	now := time.Now()
	outVal := false

	if rc.MaximumOccurrence > 0.0 && now.Before(rc.lastTimeTrue.Add(rc.MaximumOccurrence)) {
		outVal = false
	} else if rc.MinimumOccurrence > 0.0 && now.After(rc.lastTimeTrue.Add(rc.MinimumOccurrence)) {
		outVal = true
	} else if RandomlyTrue(rc.currentChance) {
		outVal = true
	} else if rc.IncrementIncrease > 0.0 && !timeIsZero(rc.lastTimeCalled) {
		diff := now.Sub(rc.lastTimeCalled)
		multiplier := diff.Seconds() / rc.IncreaseDuration.Seconds()

		//log.Printf("Diff: %v, Multiplier: %v,", diff, multiplier)

		//log.Printf("Pre-Calc Chance: %0.6f\n", rc.currentChance)
		rc.currentChance += rc.IncrementIncrease * multiplier
		//log.Printf("Post-Calc Chance: %0.6f\n", rc.currentChance)
		// 5 decimal places
		//rc.currentChance = math.Round(rc.currentChance * 100_000) / 100_000
		//log.Printf("Rounded Chance: %0.6f\n", rc.currentChance)

	}

	if outVal {
		// reset increment since we've returned true once. Yay!
		rc.currentChance = rc.DefaultChance
		rc.lastTimeTrue = now
		rc.trueCount++
	} else {
		rc.falseCount++
	}

	rc.lastTimeCalled = now
	rc.callCount++

	return outVal
}

// Define defines a control in the group
func (cr *AdvControlGroup) Define(rng *AdvancedControl) {
	cr.controls[rng.ID] = rng
}

// Execute provides the execution of the control in the group
func (cr *AdvControlGroup) Execute(ID string) bool {
	ctrl, ok := cr.controls[ID]

	if !ok {
		log.Printf("WARN: Execute ID '%s' called when not defined\n", ID)
		return false
	}

	return ctrl.Execute()
}

// Analytics returns a set of Analytics records for providing
// value to things
func (cr *AdvControlGroup) Analytics() map[string]Analytics {
	out := make(map[string]Analytics, len(cr.controls))

	for id, ctrl := range cr.controls {
		out[id] = ctrl.ToAnalytics()
	}

	return out
}

// NewAdvControlGroup builds a new AdvControlGroup for you bby
func NewAdvControlGroup() *AdvControlGroup {
	return &AdvControlGroup{
		controls: make(map[string]*AdvancedControl),
	}
}
