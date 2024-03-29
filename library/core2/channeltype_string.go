// Code generated by "stringer -type=ChannelType"; DO NOT EDIT.

package core2

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ChannelUnknown-0]
	_ = x[ChannelPublic-1]
	_ = x[ChannelPrivileged-2]
	_ = x[ChannelDirectMessage-3]
	_ = x[ChannelGroupMessage-4]
	_ = x[ChannelReadOnly-5]
	_ = x[ChannelVoice-6]
	_ = x[ChannelForum-7]
}

const _ChannelType_name = "ChannelUnknownChannelPublicChannelPrivilegedChannelDirectMessageChannelGroupMessageChannelReadOnlyChannelVoiceChannelForum"

var _ChannelType_index = [...]uint8{0, 14, 27, 44, 64, 83, 98, 110, 122}

func (i ChannelType) String() string {
	if i >= ChannelType(len(_ChannelType_index)-1) {
		return "ChannelType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ChannelType_name[_ChannelType_index[i]:_ChannelType_index[i+1]]
}
