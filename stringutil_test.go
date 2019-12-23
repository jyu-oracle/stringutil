package stringutil

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParseKeyValuePairs(t *testing.T) {
	cases := []struct {
		line   string
		prefix string
		m      map[string]string
	}{
		{"prefix key1=val1 key2=val2", "prefix", map[string]string{"key1": "val1", "key2": "val2"}},
		{"prefix key1=val1, key2=val2", "prefix", map[string]string{"key1": "val1", "key2": "val2"}},
		{"prefix key1=(val1) key2=val2", "prefix", map[string]string{"key1": "(val1)", "key2": "val2"}},
		{`prefix key1="val1" key2=val2`, "prefix", map[string]string{"key1": "val1", "key2": "val2"}},
		{`prefix key1="(val1)" key2=val2`, "prefix", map[string]string{"key1": "(val1)", "key2": "val2"}},
		{`prefix key1="(val1)", key2=val2`, "prefix", map[string]string{"key1": "(val1)", "key2": "val2"}},
		{`prefix key1="(val1)" , key2=val2`, "prefix", map[string]string{"key1": "(val1)", "key2": "val2"}},
		{"prefix key1=val1  , key2=val2  ,", "prefix", map[string]string{"key1": "val1", "key2": "val2"}},
		{"prefix key1 = val1 key2 = val2", "prefix", map[string]string{"key1": "val1", "key2": "val2"}},
		{"prefix key1 : val1 key2 : val2", "prefix", map[string]string{"key1": "val1", "key2": "val2"}},
		{"prefix key1:val1 key2:val2", "prefix", map[string]string{"key1": "val1", "key2": "val2"}},
		{"prefix key1:val1, key2:val2", "prefix", map[string]string{"key1": "val1", "key2": "val2"}},
		{"prefix key1:val1    key2:val2   ", "prefix", map[string]string{"key1": "val1", "key2": "val2"}},
		{"prefix key1:val1, key2:val2, key3:val3, key4:val4", "prefix", map[string]string{"key1": "val1", "key2": "val2", "key3": "val3", "key4": "val4"}},
		{"key1=val1 key2=val2", "", map[string]string{"key1": "val1", "key2": "val2"}},
		{"prefix key1=val1 key2=val2", "NOMATCH", nil},
		{"prefixkey1=val1 key2=val2", "prefix", map[string]string{"key1": "val1", "key2": "val2"}},
		{"key0=val0 prefix key1=val1 key2=val2", "prefix", map[string]string{"key1": "val1", "key2": "val2"}},
		{"prefix", "prefix", map[string]string{}},
		{`prefix key1="val1, from here to there" key2=val2`, "prefix", map[string]string{"key1": "val1, from here to there", "key2": "val2"}},
		{`prefix key1="val1, from here to there", key2=val2`, "prefix", map[string]string{"key1": "val1, from here to there", "key2": "val2"}},
		{`prefix key1="val1, key2=val2`, "prefix", map[string]string{"key1": `"val1`, "key2": "val2"}},
		{`waveform=0x101 (auto) mode=0x0 temp:4096 update region top=576, left=19, width=562, height=69 flags=0x0`, "", map[string]string{"waveform": "0x101", "mode": "0x0", "temp": "4096", "top": "576", "left": "19", "width": "562", "height": "69", "flags": "0x0"}},
		{`update end marker=122, end time=1558318772259, time taken=299 ms`, "", map[string]string{"marker": "122", "time": "1558318772259", "taken": "299"}},
		{`X11EM.onPointerEvent() RELEASED btn=1, (300,147) time=1558318772034`, "", map[string]string{"btn": "1", "time": "1558318772034"}},
		{"prefix key1=val1 key2=val2 key1=val3", "prefix", map[string]string{"key1": "val3", "key2": "val2"}},
		{`X11EM.onPointerEvent() RELEASED btn=1, pos=(300,147), time=1558318772034`, "", map[string]string{"btn": "1", "pos": "(300,147)", "time": "1558318772034"}},
	}

	for _, c := range cases {
		r := ParseKeyValuePairs(c.line, c.prefix)
		if debug {
			fmt.Println(r)
		}
		if !reflect.DeepEqual(c.m, r) {
			t.Errorf("ParseKeyValuePairs(%v,%v)=%v, want %v", c.line, c.prefix, r, c.m)
		}
	}
}

func TestParseValues(t *testing.T) {
	cases := []struct {
		line   string
		prefix string
		fields string
		m      map[string]string
	}{
		{"prefix val1 val2", "prefix", "key1 key2", map[string]string{"key1": "val1", "key2": "val2"}},
		{"prefix val1, val2", "prefix", "key1 key2", map[string]string{"key1": "val1", "key2": "val2"}},
		{"prefix val1   ,   val2", "prefix", "key1 key2", map[string]string{"key1": "val1", "key2": "val2"}},
		{"prefix val1 val2 val3 val4", "prefix", "", map[string]string{"#1": "val1", "#2": "val2", "#3": "val3", "#4": "val4"}},
		{"prefix val1 val2 val3 val4", "prefix", "key1 key2", map[string]string{"key1": "val1", "key2": "val2", "#3": "val3", "#4": "val4"}},
		{"prefix val1 val2", "NOMATCH", "key1 key2", nil},
		{"prefixval1 val2", "prefix", "key1 key2", map[string]string{"key1": "val1", "key2": "val2"}},
		{"uselessvalues prefix val1 val2", "prefix", "key1 key2", map[string]string{"key1": "val1", "key2": "val2"}},
		{"prefix", "prefix", "", map[string]string{}},
		{`prefix "val1, from here to there" val2`, "prefix", "key1 key2", map[string]string{"key1": "val1, from here to there", "key2": "val2"}},
		{`prefix "val1, from here to there", val2`, "prefix", "key1 key2", map[string]string{"key1": "val1, from here to there", "key2": "val2"}},
		{`prefix (val1, from here to there) val2`, "prefix", "key1 key2", map[string]string{"key1": "(val1, from here to there)", "key2": "val2"}},
		{`prefix "val1 val2`, "prefix", "key1 key2", map[string]string{"key1": `"val1`, "key2": "val2"}},
		{"prefix val0 val1 val2 val3", "prefix", "_ key1 key2 key3", map[string]string{"#1": "val0", "key1": "val1", "key2": "val2", "key3": "val3"}},
		{"prefix val1() val2", "prefix", "key1 key2", map[string]string{"key1": "val1()", "key2": "val2"}},
		{"prefix val1 ()val2", "prefix", "key1 key2", map[string]string{"key1": "val1", "key2": "()val2"}},
		{"prefix val1 () val2", "prefix", "key1 _ key2", map[string]string{"key1": "val1", "#2": "()", "key2": "val2"}},
		{"prefix s1=val1 s2=val2", "prefix", "key1 key2", map[string]string{"key1": "s1=val1", "key2": "s2=val2"}},
	}

	for _, c := range cases {
		r := ParseValues(c.line, c.prefix, strings.Fields(c.fields))
		if debug {
			fmt.Println(r)
		}
		if !reflect.DeepEqual(c.m, r) {
			t.Errorf("ParseValues(%v,%v,%v)=%v, want %v", c.line, c.prefix, c.fields, r, c.m)
		}
	}
}

func TestParseComplexData(t *testing.T) {
	cases := []struct {
		line   string
		prefix string
		fields string
		m      map[string]string
	}{
		{"prefix key1=val1 val2 key3=val3", "prefix", "_ key2 _", map[string]string{"key1": "val1", "key2": "val2", "key3": "val3"}},
		{`prefix key1="val1" val2 key3=val3`, "prefix", "_ key2 _", map[string]string{"key1": "val1", "key2": "val2", "key3": "val3"}},
		{"prefix key1=(val1) val2 key3=val3", "prefix", "_ key2 _", map[string]string{"key1": "(val1)", "key2": "val2", "key3": "val3"}},
		{"prefix key1:val1 val2 key3:val3", "prefix", "_ key2 _", map[string]string{"key1": "val1", "key2": "val2", "key3": "val3"}},
		{"prefix key1:val1 val2 key3:val3, key4:val4", "prefix", "_ key2 _", map[string]string{"key1": "val1", "key2": "val2", "key3": "val3", "key4": "val4"}},
		{"prefix key1=val1 val2 key3=val3", "prefix", "MyKey1 key2 MyKey3", map[string]string{"MyKey1": "val1", "key2": "val2", "MyKey3": "val3"}},
		{"prefix key1=val1 val2 key3=val3", "prefix", "_ _ _", map[string]string{"key1": "val1", "#2": "val2", "key3": "val3"}},
		{"prefix key1=val1 val2 key3=val3", "prefix", "", map[string]string{"key1": "val1", "#2": "val2", "key3": "val3"}},
		{`X11EM.onPointerEvent() RELEASED btn=1, (300,147) time=1558318772034`, "", "evt subtype _ pos _", map[string]string{"evt": "X11EM.onPointerEvent()", "subtype": "RELEASED", "btn": "1", "pos": "(300,147)", "time": "1558318772034"}},
		{`perfScenario CVM received X11 ButtonRelease button=1 time=1560928559.821221`, "CVM received", "_ evt", map[string]string{"#1": "X11", "evt": "ButtonRelease", "button": "1", "time": "1560928559.821221"}},
		{`prefix key1="val1", val2`, "prefix", "_ key2", map[string]string{"key1": "val1", "key2": "val2"}},
		{`prefix key1="val1,val11, val12", val2`, "prefix", "_ key2", map[string]string{"key1": "val1,val11, val12", "key2": "val2"}},
		{`prefix key1=val1, "val2"`, "prefix", "_ key2", map[string]string{"key1": "val1", "key2": "val2"}},
		{"prefix val1 val2", "NOMATCH", "key1 key2", nil},
		{"MMUPageSize:           4 kB", "", "_ unit", map[string]string{"MMUPageSize": "4", "unit": "kB"}},
	}

	for _, c := range cases {
		r := ParseComplexData(c.line, c.prefix, strings.Fields(c.fields))
		if debug {
			fmt.Println(r)
		}
		if !reflect.DeepEqual(c.m, r) {
			t.Errorf("ParseComplexData(%v, %v, %v)=%v, want %v", c.line, c.prefix, c.fields, r, c.m)
		}
	}
}
