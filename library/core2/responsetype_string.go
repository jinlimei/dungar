// Code generated by "stringer -type=ResponseType"; DO NOT EDIT.

package core2

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ResponseTypeUnknown-0]
	_ = x[ResponseTypeBasic-1]
	_ = x[ResponseTypeReaction-2]
	_ = x[ResponseTypeDirectMessage-3]
}

const _ResponseType_name = "ResponseTypeUnknownResponseTypeBasicResponseTypeReactionResponseTypeDirectMessage"

var _ResponseType_index = [...]uint8{0, 19, 36, 56, 81}

func (i ResponseType) String() string {
	if i >= ResponseType(len(_ResponseType_index)-1) {
		return "ResponseType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ResponseType_name[_ResponseType_index[i]:_ResponseType_index[i+1]]
}
