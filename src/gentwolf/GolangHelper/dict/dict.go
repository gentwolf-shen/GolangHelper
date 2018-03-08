package dict

import (
	"encoding/json"
	"io/ioutil"
)

var items map[string]string

func Load(filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &items)
	return err
}

func LoadFromStr(str string) error {
	return json.Unmarshal([]byte(str), &items)
}

func Get(key string) string {
	return items[key]
}

func Set(key, value string) {
	items[key] = value
}
