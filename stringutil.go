package stringutil

import (
	"fmt"
	"regexp"
	"strings"
)

const debug = false

var rexKeyValue, rexData, rexComplexData *regexp.Regexp

func stripDoubleQuote(str string) string {
	if len(str) < 2 {
		return str
	}
	if str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	return str
}

func removePrefix(line *string, prefix string) bool {
	if prefix != "" {
		idx := strings.Index(*line, prefix)
		if idx == -1 {
			return false
		}
		*line = (*line)[idx+len(prefix):]
	}
	return true
}

func initRexKeyValue() {
	/*
		( key [=|:] value )
			value format : value, "value", (value)

		([^\s]+) \s* [=|:] \s* ("[^"]+" | \([^\)]+\) | [^\s,]+ )

		\s             whitespace (== [\t\n\f\r ])
		"[^"]+"        double quoted string
		[^\s,]+        string not containing whitespace and ',' character
	*/
	if rexKeyValue == nil {
		rexKeyValue = regexp.MustCompile(`([^\s]+)\s*[=|:]\s*("[^"]+"|\([^\)]+\)|[^\s,]+)`)
	}
}

func initRexData() {
	/*
		( value )
			value format : value, "value", (value)

		("[^"]+" | \([^\)]+\) | [^\s,]+)
	*/
	if rexData == nil {
		rexData = regexp.MustCompile(`("[^"]+"|\([^\)]+\)|[^\s,]+)`)
	}
}

func initRexComplexData() {
	/*
		( key [=|:] value ) | ( value )
			value format : value, "value", (value)

		([^\s]+) \s* [=|:] \s* ("[^"]+" | \([^\)]+\) | [^\s,]+ )  |  ("[^"]+" | \([^\)]+\) | [^\s,]+)
	*/
	if rexComplexData == nil {
		rexComplexData = regexp.MustCompile(`([^\s]+)\s*[=|:]\s*("[^"]+"|\([^\)]+\)|[^\s,]+)|("[^"]+"|\([^\)]+\)|[^\s,]+)`)
	}
}

// GetParsedKeyValueMap returns key-value pairs directly parsed from the input string as a map.
// Expected input pattern : ( key  [=|:] value )
// Returns nil if prefix matching fails.
func GetParsedKeyValueMap(line string, prefix string) map[string]string {
	if !removePrefix(&line, prefix) {
		return nil
	}
	if debug {
		fmt.Println(line)
	}
	retMap := make(map[string]string)

	initRexKeyValue()
	data := rexKeyValue.FindAllStringSubmatch(line, -1)
	for _, v := range data {
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
		retMap[key] = value
		if debug {
			fmt.Printf("[%v]->[%v]\n", key, value)
		}
	}
	return retMap
}

// GetSlicedDataMap returns sliced data as a map.
// 'values' are directly parsed from the input string, and the corresponding 'keys' are provided as field param.
// If "_" string is specified for fields, the default key value ("#" + index) will be used. Otherwise the specified field value will be used as key.
// Expected input pattern : ( value )
// Returns nil if prefix matching fails.
func GetSlicedDataMap(line string, prefix string, fields []string) map[string]string {
	if !removePrefix(&line, prefix) {
		return nil
	}
	if debug {
		fmt.Println(line)
	}
	retMap := make(map[string]string)

	initRexData()
	data := rexData.FindAllStringSubmatch(line, -1)
	for i, v := range data {
		if len(v) != 2 {
			return nil
		}
		var key string
		if i < len(fields) && fields[i] != "_" {
			key = fields[i]
		} else {
			key = fmt.Sprintf("#%d", i+1)
		}
		value := stripDoubleQuote(v[1])
		retMap[key] = value
		if debug {
			fmt.Printf("[%v]->[%v]\n", key, value)
		}
	}
	return retMap
}

// GetSlicedDataArray returns sliced data as an array.
// Expected input pattern : ( value )
// Returns nil if prefix matching fails.
func GetSlicedDataArray(line string, prefix string) []string {
	if !removePrefix(&line, prefix) {
		return nil
	}
	if debug {
		fmt.Println(line)
	}
	retArray := make([]string, 0)

	initRexData()
	data := rexData.FindAllStringSubmatch(line, -1)
	for i, v := range data {
		if len(v) != 2 {
			return nil
		}
		value := stripDoubleQuote(v[1])
		if debug {
			fmt.Printf("(%v) [%v]\n", i, value)
		}
		retArray = append(retArray, value)
	}
	return retArray
}

// GetParsedComplexDataMap returns complex data (key-value pair or valueonly data) parsed from the input string as a map.
// 'values' are directly parsed from the input string, and the 'keys' are either from input string or fields param.
// For the case of key-value pair, the default 'key' is directly from the input string.
// For the case of valueonly data, the default 'key' is "# + index" string.
// If "_" string is specified for fields, the default 'key' will be used. Otherwise the specified field value will be used as key (this is overriding)
// Expected input pattern : ( key  [=|:] value ) | ( value )
// Returns nil if prefix matching fails.
func GetParsedComplexDataMap(line string, prefix string, fields []string) map[string]string {
	if !removePrefix(&line, prefix) {
		return nil
	}
	if debug {
		fmt.Println(line)
	}
	retMap := make(map[string]string)

	initRexComplexData()
	data := rexComplexData.FindAllStringSubmatch(line, -1)
	for i, v := range data {
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
		retMap[key] = value
		if debug {
			fmt.Printf("[%v]->[%v]\n", key, value)
		}
	}
	return retMap
}

// GetSlicedComplexDataArray returns complex data (key-value pair or valueonly data) parsed from the input string as an array.
// key-value pair is normalized to use "=".
// Expected input pattern : ( key  [=|:] value ) | ( value )
// Returns nil if prefix matching fails.
func GetSlicedComplexDataArray(line string, prefix string) []string {
	if !removePrefix(&line, prefix) {
		return nil
	}
	if debug {
		fmt.Println(line)
	}
	retArray := make([]string, 0)

	initRexComplexData()
	data := rexComplexData.FindAllStringSubmatch(line, -1)
	for i, v := range data {
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
		retArray = append(retArray, str)
	}
	return retArray
}
