package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

func ReadJson(r io.ReadCloser, obj interface{}) error {

	body, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	if err := r.Close(); err != nil {
		panic(err)
	}

	return ReadJsonByByte(body, obj)
}

func ReadJsonByByte(body []byte, obj interface{}) error {
	mdz := json.NewDecoder(bytes.NewBuffer(body))

	mdz.UseNumber()
	err := mdz.Decode(obj)

	if err != nil {
		return err
	}
	return nil
}

func WriteJson(w io.Writer, obj interface{}) {

	if obj == nil {
		io.WriteString(w, "{}")
		return
	}

	if objStr, ok := obj.(string); ok {
		if objStr == "" {
			io.WriteString(w, "{}")
			return
		}
	}
	jsonData, _ := json.Marshal(obj)
	io.WriteString(w, string(jsonData))
}

//将对象转换为JSON
func ToJson(obj interface{}) string {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	return string(jsonData)
}
