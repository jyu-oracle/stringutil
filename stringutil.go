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
		key  [=|:]  ( "value" | value ) {,}
		\s* ([^\s,]+) \s*  [=|:]  \s* ( ("[^"]+") | ([^\s,]+) )  [,]?

		\s             whitespace (== [\t\n\f\r ])
		"[^"]+"        double quoted string
		[^\s,]+        string not containing whitespace and ',' character
	*/
	rex := regexp.MustCompile(`\s*([^\s,]+)\s*[=|:]\s*(("[^"]+")|([^\s,]+))[,]?`)
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
		( "value" | value ) {,}
		\s* ( ("[^"]+") | ([^\s,]+) )  [,]?

		\s             whitespace (== [\t\n\f\r ])
		"[^"]+"        double quoted string
		[^\s,]+        string not containing whitespace and ',' character
	*/
	rex := regexp.MustCompile(`\s*(("[^"]+")|([^\s,]+))[,]?`)
	data := rex.FindAllStringSubmatch(line, -1)

	for i, v := range data {
		var key string
		if i < len(fields) {
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
