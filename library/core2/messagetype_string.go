// Code generated by "stringer -type=MessageType"; DO NOT EDIT.

package core2

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[MessageTypeUnknown-0]
	_ = x[MessageTypeBasic-1]
	_ = x[MessageTypeMe-2]
	_ = x[MessageTypeChanged-3]
	_ = x[MessageTypeDeleted-4]
	_ = x[MessageTypeBroadcast-5]
}

const _MessageType_name = "MessageTypeUnknownMessageTypeBasicMessageTypeMeMessageTypeChangedMessageTypeDeletedMessageTypeBroadcast"

var _MessageType_index = [...]uint8{0, 18, 34, 47, 65, 83, 103}

func (i MessageType) String() string {
	if i >= MessageType(len(_MessageType_index)-1) {
		return "MessageType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _MessageType_name[_MessageType_index[i]:_MessageType_index[i+1]]
}
