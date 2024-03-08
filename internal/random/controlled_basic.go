package random

import "log"

// BasicControl is our basic variant of AdvancedControl which just provides
// a Chance and ID to get off the ground, plus some internal analytics.
type BasicControl struct {
	// ID is the unique ID for this control
	ID string

	// Chance is the % chance to execute as true.
	// Range is [0,1] where 0.0 is always false, 1.0 always true
	Chance float64

	// callCount is the count of calls to Execute for this struct
	callCount int

	// trueCount is how many times Execute returned true
	trueCount int

	// falseCount is how many times Execute returned false
	falseCount int
}

// BasicControlGroup is for managing the basic controls (which are
// just an ID and a % chance)
type BasicControlGroup struct {
	controls map[string]*BasicControl
}

// ToAnalytics converts the control to an Analytics struct
func (bc *BasicControl) ToAnalytics() Analytics {
	return Analytics{
		ID:         bc.ID,
		TrueCount:  bc.trueCount,
		FalseCount: bc.falseCount,
		CallCount:  bc.callCount,
		Chance:     bc.Chance,
	}
}

// Execute just does a basic execution of the control
func (bc *BasicControl) Execute() bool {
	bc.callCount++

	if RandomlyTrue(bc.Chance) {
		bc.trueCount++
		return true
	}

	bc.falseCount++
	return false
}

// Define defines a control in the controls map
func (bc *BasicControlGroup) Define(ctrl *BasicControl) {
	bc.controls[ctrl.ID] = ctrl
}

// ExecuteID takes in only an ID and executes it. Will always return false
// if there is no ID in the controls map
func (bc *BasicControlGroup) ExecuteID(ID string) bool {
	ctrl, ok := bc.controls[ID]

	if !ok {
		log.Printf("WARN: ExecuteID called on ID '%s' which is not set\n", ID)
		return false
	}

	return ctrl.Execute()
}

// Execute behaves like RandomlyTrue but will store analytics in the
// control group. If it is not defined, it will be defined.
func (bc *BasicControlGroup) Execute(id string, chance float64) bool {
	ctrl, ok := bc.controls[id]
	if !ok {
		ctrl = &BasicControl{
			ID:     id,
			Chance: chance,
		}

		bc.controls[id] = ctrl
	}

	// Allow for the opportunity for the chance to change dynamically
	// even though we have the same ID
	ctrl.Chance = chance

	return ctrl.Execute()
}

// Analytics returns a list of analytics for used/activated calls
func (bc *BasicControlGroup) Analytics() map[string]Analytics {
	out := make(map[string]Analytics)

	for id, ctrl := range bc.controls {
		out[id] = ctrl.ToAnalytics()
	}

	return out
}

// NewBasicControlGroup builds a new NewBasicControlGroup for you bby
func NewBasicControlGroup() *BasicControlGroup {
	return &BasicControlGroup{
		controls: make(map[string]*BasicControl),
	}
}
