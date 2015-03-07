package goics

import (
	"strings"
)

const (
	vParamSep = ";"
)

// IcsNode is a basic token.., with, key, val, and extra params
// to define the type of val.
type IcsNode struct {
	Key    string
	Val    string
	Params map[string]string
}

// how many params has a token
func (n *IcsNode) ParamsLen() int {
	if n.Params == nil {
		return 0
	}
	return len(n.Params)
}

// If only has one extra param, returns it..
// as key, val
func (n *IcsNode) GetOneParam() (string, string) {
	if n.ParamsLen() == 0 {
		return "", ""
	}
	var key, val string
	for k, v := range n.Params {
		key, val = k, v
		break
	}
	return key, val
}

// Decodes a line extracting key, val and extra params
// linked to key..
func DecodeLine(line string) *IcsNode {
	if strings.Contains(line, keySep) == false {
		return &IcsNode{}
	}
	key, val := getKeyVal(line)
	//@todo test if val containes , multipleparams
	if strings.Contains(key, vParamSep) == false {
		return &IcsNode{
			Key: key,
			Val: val,
		}
	} else {
		// Extract key
		first_param := strings.Index(key, vParamSep)
		real_key := key[0:first_param]
		n := &IcsNode{
			Key: real_key,
			Val: val,
		}
		// Extract params
		params := key[first_param+1:]
		n.Params = decode_params(params)
		return n
	}
	return nil
}

// decode extra params linked in key val in the form
// key;param1=val1:val
func decode_params(arr string) map[string]string {

	p := make(map[string]string)
	var is_quoted = false
	var is_param bool = true
	var cur_param string
	var cur_val string
	for _, c := range arr {
		switch {
		// if string is quoted, wait till next quote
		// and capture content
		case c == '"':
			if is_quoted == false {
				is_quoted = true
			} else {
				p[cur_param] = cur_val
				is_quoted = false
			}
		case c == '=' && is_quoted == false:
			is_param = false
		case c == ';' && is_quoted == false:
			is_param = true
			p[cur_param] = cur_val
			cur_param = ""
			cur_val = ""
		default:
			if is_param {
				cur_param = cur_param + string(c)
			} else {
				cur_val = cur_val + string(c)
			}
		}
	}
	p[cur_param] = cur_val
	return p

}

// Returns a key, val... for a line..
func getKeyVal(s string) (key, value string) {
	p := strings.SplitN(s, keySep, 2)
	return p[0], p[1]
}
