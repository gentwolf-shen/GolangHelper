package dict

import (
	"encoding/json"
	"io/ioutil"
)

var items map[string]string

func Init(filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &items)
	return err
}

func Get(key string) string {
	return items[key]
}

func Set(key, value string) {
	items[key] = value
}
