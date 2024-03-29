// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package core2

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TokenSpace-0]
	_ = x[TokenWord-1]
	_ = x[TokenUserID-2]
	_ = x[TokenChanID-3]
	_ = x[TokenGroupID-4]
	_ = x[TokenRoleID-5]
	_ = x[TokenEmoticon-6]
	_ = x[TokenURL-7]
}

const _TokenType_name = "TokenSpaceTokenWordTokenUserIDTokenChanIDTokenGroupIDTokenRoleIDTokenEmoticonTokenURL"

var _TokenType_index = [...]uint8{0, 10, 19, 30, 41, 53, 64, 77, 85}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
