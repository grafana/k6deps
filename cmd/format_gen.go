// Code generated by "enumer -type=format -transform=lower -trimprefix format -output format_gen.go"; DO NOT EDIT.

package cmd

import (
	"fmt"
	"strings"
)

const _formatName = "jsontextjs"

var _formatIndex = [...]uint8{0, 4, 8, 10}

const _formatLowerName = "jsontextjs"

func (i format) String() string {
	if i < 0 || i >= format(len(_formatIndex)-1) {
		return fmt.Sprintf("format(%d)", i)
	}
	return _formatName[_formatIndex[i]:_formatIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _formatNoOp() {
	var x [1]struct{}
	_ = x[formatJSON-(0)]
	_ = x[formatText-(1)]
	_ = x[formatJS-(2)]
}

var _formatValues = []format{formatJSON, formatText, formatJS}

var _formatNameToValueMap = map[string]format{
	_formatName[0:4]:       formatJSON,
	_formatLowerName[0:4]:  formatJSON,
	_formatName[4:8]:       formatText,
	_formatLowerName[4:8]:  formatText,
	_formatName[8:10]:      formatJS,
	_formatLowerName[8:10]: formatJS,
}

var _formatNames = []string{
	_formatName[0:4],
	_formatName[4:8],
	_formatName[8:10],
}

// formatString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func formatString(s string) (format, error) {
	if val, ok := _formatNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _formatNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to format values", s)
}

// formatValues returns all values of the enum
func formatValues() []format {
	return _formatValues
}

// formatStrings returns a slice of all String values of the enum
func formatStrings() []string {
	strs := make([]string, len(_formatNames))
	copy(strs, _formatNames)
	return strs
}

// IsAformat returns "true" if the value is listed in the enum definition. "false" otherwise
func (i format) IsAformat() bool {
	for _, v := range _formatValues {
		if i == v {
			return true
		}
	}
	return false
}