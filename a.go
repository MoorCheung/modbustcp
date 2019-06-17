package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	var m = make(map[string]interface{})
	m["a"] = 1
	m["v"] = 2
	m["b"] = 3
	bytes, e := json.Marshal(m)
	fmt.Println(e,string(bytes))
}