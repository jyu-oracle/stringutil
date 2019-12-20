package stringutil

import (
	"fmt"
	"regexp"
	"strings"
)

const debug = false

func stripDoubleQuote(str string) string {
	if len(str) < 2 {
		return str
	}
	if str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	return str
}

// ParseKeyValuePairs returns key-value pairs directly parsed from the input string.
// Expected input pattern : key  [=|:]  ( "value" | value )
func ParseKeyValuePairs(line string, prefix string) map[string]string {
	if prefix != "" {
		idx := strings.Index(line, prefix)
		if idx == -1 {
			return nil
		}
		line = line[idx+len(prefix):]
	}
	if debug {
		fmt.Println(line)
	}
	retMap := make(map[string]string)

	/*
		key  [=|:]  ( "value" | value )
		([^\s]+) \s* [=|:] \s* ( ("[^"]+") | ([^\s,]+) )

		\s             whitespace (== [\t\n\f\r ])
		"[^"]+"        double quoted string
		[^\s,]+        string not containing whitespace and ',' character
	*/
	rex := regexp.MustCompile(`([^\s]+)\s*[=|:]\s*(("[^"]+")|([^\s,]+))`)
	data := rex.FindAllStringSubmatch(line, -1)

	for _, v := range data {
		key := v[1]
		value := stripDoubleQuote(v[2])
		retMap[key] = value
		if debug {
			fmt.Printf("[%v]->[%v]\n", key, value)
		}
	}
	return retMap
}

// ParseValues returns key-value pairs.
// Keys are explicitly provided as fields param, and values are parsed from the input string.
// If specify "_" string for fields, the default key value ("#" + index) will be used. Otherwise the specified field value will be used as key.
// For the case of value which doesn't have matching key, the default key("#" + index) is used.
// Expected input pattern : ( "value" | value )
func ParseValues(line string, prefix string, fields []string) map[string]string {
	if prefix != "" {
		idx := strings.Index(line, prefix)
		if idx == -1 {
			return nil
		}
		line = line[idx+len(prefix):]
	}
	if debug {
		fmt.Println(line)
	}
	retMap := make(map[string]string)

	/*
		( "value" | value )
		( ("[^"]+") | ([^\s,]+) )

		\s             whitespace (== [\t\n\f\r ])
		"[^"]+"        double quoted string
		[^\s,]+        string not containing whitespace and ',' character
	*/
	rex := regexp.MustCompile(`(("[^"]+")|([^\s,]+))`)
	data := rex.FindAllStringSubmatch(line, -1)

	for i, v := range data {
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

// ParseData returns key-value pairs and values parsed from the input string.
// fields param is used to supply additional information for keys. If specify "_" string for fields, the default key value (key value parsed from input or #idx) will be used. Otherwise the specified field value will override the default key value.
// Expected input pattern : ( { key  [=|:] }  ( "value" | value ) )
func ParseData(line string, prefix string, fields []string) map[string]string {
	if prefix != "" {
		idx := strings.Index(line, prefix)
		if idx == -1 {
			return nil
		}
		line = line[idx+len(prefix):]
	}
	if debug {
		fmt.Println(line)
	}
	retMap := make(map[string]string)

	rex1 := regexp.MustCompile(`([^\s]+)`)
	rex2 := regexp.MustCompile(`([^\s]+)[=|:](("[^"]+")|([^\s,]+))`)
	data1 := rex1.FindAllStringSubmatch(line, -1)

	for i, v := range data1 {
		var key string
		value := stripDoubleQuote(v[1])
		data2 := rex2.FindAllStringSubmatch(value, -1)
		if len(data2) == 0 {
			if i < len(fields) && fields[i] != "_" {
				key = fields[i]
			} else {
				key = fmt.Sprintf("#%d", i+1)
			}
			retMap[key] = value
			if debug {
				fmt.Printf("[%v]->[%v]\n", key, value)
			}
		} else {
			for _, vv := range data2 {
				if i < len(fields) && fields[i] != "_" {
					key = fields[i]
				} else {
					key = vv[1]
				}
				value = stripDoubleQuote(vv[2])
				retMap[key] = value
				if debug {
					fmt.Printf("* [%v]->[%v]\n", key, value)
				}
			}
		}
	}
	return retMap
}
