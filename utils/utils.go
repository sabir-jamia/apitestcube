package utils

import (
	"bytes"
	"encoding/json"
)

func GetPostData(data string, userInput interface{}) {
	var emptyBytes []byte
	escapedJson := bytes.NewBuffer(emptyBytes)
	json.HTMLEscape(escapedJson, []byte(data))
	json.Unmarshal(escapedJson.Bytes(), &userInput)
}