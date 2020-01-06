package stringutil

import (
	"fmt"
	"regexp"
	"strings"
)

const debug = false

var keyValuePattern, valuesOnlyPattern, pairsAndValuesPattern *regexp.Regexp

func stripDoubleQuote(str string) string {
	if len(str) < 2 {
		return str
	}
	if str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	return str
}

func removePrefix(str *string, prefix string) bool {
	if prefix != "" {
		idx := strings.Index(*str, prefix)
		if idx == -1 {
			return false
		}
		*str = (*str)[idx+len(prefix):]
	}
	return true
}

func initKeyValuePattern() {
	/*
		( key [=|:] value )
			value format : value, "value", (value)

		([^\s]+) \s* [=|:] \s* ("[^"]+" | \([^\)]+\) | [^\s,]+ )

		\s             whitespace (== [\t\n\f\r ])
		"[^"]+"        double quoted string
		[^\s,]+        string not containing whitespace and ',' character
	*/
	if keyValuePattern == nil {
		keyValuePattern = regexp.MustCompile(`([^\s]+)\s*[=|:]\s*("[^"]+"|\([^\)]+\)|[^\s,]+)`)
	}
}

func initValuesOnlyPattern() {
	/*
		( value )
			value format : value, "value", (value)

		("[^"]+" | \([^\)]+\) | [^\s,]+)
	*/
	if valuesOnlyPattern == nil {
		valuesOnlyPattern = regexp.MustCompile(`("[^"]+"|\([^\)]+\)|[^\s,]+)`)
	}
}

func initPairsAndValuesPattern() {
	/*
		( key [=|:] value ) | ( value )
			value format : value, "value", (value)

		([^\s]+) \s* [=|:] \s* ("[^"]+" | \([^\)]+\) | [^\s,]+ )  |  ("[^"]+" | \([^\)]+\) | [^\s,]+)
	*/
	if pairsAndValuesPattern == nil {
		pairsAndValuesPattern = regexp.MustCompile(`([^\s]+)\s*[=|:]\s*("[^"]+"|\([^\)]+\)|[^\s,]+)|("[^"]+"|\([^\)]+\)|[^\s,]+)`)
	}
}

// ExtractKeyValuePairs returns key-value pairs directly parsed from the input string as a map.
// Expected input pattern : ( key  [=|:] value )
// Returns nil if prefix matching fails.
func ExtractKeyValuePairs(str string, prefix string) map[string]string {
	if !removePrefix(&str, prefix) {
		return nil
	}
	if debug {
		fmt.Println(str)
	}
	pairs := make(map[string]string)

	initKeyValuePattern()
	matches := keyValuePattern.FindAllStringSubmatch(str, -1)
	for _, v := range matches {
		/*
			v[0] : match of the entire expression (key=value)
			v[1] : match of the 1st parenthesized subexpr (key)
			v[2] : match of the 2nd parenthesized subexpr (value)
		*/
		if len(v) != 3 {
			return nil
		}
		key := v[1]
		value := stripDoubleQuote(v[2])
		pairs[key] = value
		if debug {
			fmt.Printf("[%v]->[%v]\n", key, value)
		}
	}
	return pairs
}

// ExtractValuesWithFields returns sliced input as a map.
// 'values' are directly parsed from the input string, and the corresponding 'keys' are provided as fields param.
// If "_" string is specified for fields, the default key value ("#" + index) will be used. Otherwise the specified field value will be used as key.
// Expected input pattern : ( value )
// Returns nil if prefix matching fails.
func ExtractValuesWithFields(str string, prefix string, fields []string) map[string]string {
	if !removePrefix(&str, prefix) {
		return nil
	}
	if debug {
		fmt.Println(str)
	}
	pairs := make(map[string]string)

	initValuesOnlyPattern()
	matches := valuesOnlyPattern.FindAllStringSubmatch(str, -1)
	for i, v := range matches {
		if len(v) != 2 {
			return nil
		}
		var key, value string
		if i < len(fields) && fields[i] != "_" {
			key = fields[i]
		} else {
			key = fmt.Sprintf("#%d", i+1)
		}
		value = stripDoubleQuote(v[1])
		pairs[key] = value
		if debug {
			fmt.Printf("[%v]->[%v]\n", key, value)
		}
	}
	return pairs
}

// Split returns sliced data as an array.
// Expected input pattern : ( value )
// Returns nil if prefix matching fails.
func Split(str string, prefix string) []string {
	if !removePrefix(&str, prefix) {
		return nil
	}
	if debug {
		fmt.Println(str)
	}
	ret := make([]string, 0)

	initValuesOnlyPattern()
	matches := valuesOnlyPattern.FindAllStringSubmatch(str, -1)
	for i, v := range matches {
		if len(v) != 2 {
			return nil
		}
		value := stripDoubleQuote(v[1])
		if debug {
			fmt.Printf("(%v) [%v]\n", i, value)
		}
		ret = append(ret, value)
	}
	return ret
}

// ExtractKeyValuePairsWithFields returns key-value pair and/or value-only data parsed from the input string as a map.
// 'values' are directly parsed from the input string, and the 'keys' are either from input string or fields param.
// For the case of key-value pair, the default 'key' is directly from the input string.
// For the case of valueonly data, the default 'key' is "# + index" string.
// If "_" string is specified for fields, the default 'key' will be used. Otherwise the specified field value will be used as key (this is overriding)
// Expected input pattern : ( key  [=|:] value ) | ( value )
// Returns nil if prefix matching fails.
func ExtractKeyValuePairsWithFields(str string, prefix string, fields []string) map[string]string {
	if !removePrefix(&str, prefix) {
		return nil
	}
	if debug {
		fmt.Println(str)
	}
	pairs := make(map[string]string)

	initPairsAndValuesPattern()
	matches := pairsAndValuesPattern.FindAllStringSubmatch(str, -1)
	for i, v := range matches {
		/*
			v[0] : match of the entire expression (key=value | valueonly)
			v[1] : match of the 1st parenthesized subexpr (key)
			v[2] : match of the 2nd parenthesized subexpr (value)
			v[3] : match of the 3rd parenthesized subexpr (valueonly)
		*/
		if len(v) != 4 {
			return nil
		}
		var key, value string
		if v[1] == "" && v[2] == "" {
			if i < len(fields) && fields[i] != "_" {
				key = fields[i]
			} else {
				key = fmt.Sprintf("#%d", i+1)
			}
			value = stripDoubleQuote(v[3])
		} else {
			if i < len(fields) && fields[i] != "_" {
				key = fields[i]
			} else {
				key = v[1]
			}
			value = stripDoubleQuote(v[2])
		}
		pairs[key] = value
		if debug {
			fmt.Printf("[%v]->[%v]\n", key, value)
		}
	}
	return pairs
}

// SplitPairsAndValues returns key-value pair and/or value-only data parsed from the input string as an array.
// key-value pair is normalized to use "=".
// Expected input pattern : ( key  [=|:] value ) | ( value )
// Returns nil if prefix matching fails.
func SplitPairsAndValues(str string, prefix string) []string {
	if !removePrefix(&str, prefix) {
		return nil
	}
	if debug {
		fmt.Println(str)
	}
	ret := make([]string, 0)

	initPairsAndValuesPattern()
	matches := pairsAndValuesPattern.FindAllStringSubmatch(str, -1)
	for i, v := range matches {
		/*
			v[0] : match of the entire expression (key=value | valueonly)
			v[1] : match of the 1st parenthesized subexpr (key)
			v[2] : match of the 2nd parenthesized subexpr (value)
			v[3] : match of the 3rd parenthesized subexpr (valueonly)
		*/
		if len(v) != 4 {
			return nil
		}
		var str string
		if v[1] == "" && v[2] == "" {
			value := stripDoubleQuote(v[3])
			str = value
		} else {
			key := v[1]
			value := stripDoubleQuote(v[2])
			str = fmt.Sprintf("%s=%s", key, value)
		}
		if debug {
			fmt.Printf("(%v) [%v]\n", i, str)
		}
		ret = append(ret, str)
	}
	return ret
}
