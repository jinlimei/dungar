package core2

// IncomingBlocks is our (theoretically) standardized means of working with the more
// complex data setups from the Slack API (and potentially other APIs)
type IncomingBlocks struct {
	BlockSet []IncomingBlock
}

// IncomingBlock is a given Block with a special Type and potentially a sub-set of elements
// via Elements
type IncomingBlock struct {
	BlockID  string
	Type     string
	Elements []IncomingBlockElement
}

// IncomingBlockElement is a specific element of an IncomingBlock with a specific Type
// and potentially further children via Elements
type IncomingBlockElement struct {
	Type     string
	URL      string
	Text     string
	Elements []IncomingBlockElement
}
