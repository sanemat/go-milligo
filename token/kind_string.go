// Code generated by "stringer -type=Kind"; DO NOT EDIT.

package token

import "strconv"

const _Kind_name = "RESERVEDNUMEOF"

var _Kind_index = [...]uint8{0, 8, 11, 14}

func (i Kind) String() string {
	if i < 0 || i >= Kind(len(_Kind_index)-1) {
		return "Kind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Kind_name[_Kind_index[i]:_Kind_index[i+1]]
}
