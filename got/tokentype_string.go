// Code generated by "stringer -type=tokenType"; DO NOT EDIT.

package got

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[itemIgnore-0]
	_ = x[itemError-1]
	_ = x[itemStrictBlock-2]
	_ = x[itemNamedBlock-3]
	_ = x[itemEndBlock-4]
	_ = x[itemSubstitute-5]
	_ = x[itemInclude-6]
	_ = x[itemEnd-7]
	_ = x[itemConvert-8]
	_ = x[itemGo-9]
	_ = x[itemText-10]
	_ = x[itemRun-11]
	_ = x[itemString-12]
	_ = x[itemBool-13]
	_ = x[itemInt-14]
	_ = x[itemUInt-15]
	_ = x[itemFloat-16]
	_ = x[itemInterface-17]
	_ = x[itemBytes-18]
	_ = x[itemComment-19]
	_ = x[itemEOF-20]
	_ = x[itemBackup-21]
	_ = x[itemIf-22]
	_ = x[itemElse-23]
	_ = x[itemFor-24]
}

const _tokenType_name = "itemIgnoreitemErroritemStrictBlockitemNamedBlockitemEndBlockitemSubstituteitemIncludeitemEnditemConvertitemGoitemTextitemRunitemStringitemBoolitemIntitemUIntitemFloatitemInterfaceitemBytesitemCommentitemEOFitemBackupitemIfitemElseitemFor"

var _tokenType_index = [...]uint8{0, 10, 19, 34, 48, 60, 74, 85, 92, 103, 109, 117, 124, 134, 142, 149, 157, 166, 179, 188, 199, 206, 216, 222, 230, 237}

func (i tokenType) String() string {
	if i < 0 || i >= tokenType(len(_tokenType_index)-1) {
		return "tokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _tokenType_name[_tokenType_index[i]:_tokenType_index[i+1]]
}
