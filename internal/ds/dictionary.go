package ds

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
)

type Pair struct {
	Key     interface{}
	Value   interface{}
	element *list.Element
}

type Dictionary []Pair

func NewKeyValue(input ...interface{}) *Pair {
	if len(input) == 2 {
		return &Pair{Key: input[0], Value: input[2]}
	}
	return &Pair{}
}

func NewDictionary() *Dictionary {
	return &Dictionary{}
}

func (dict *Dictionary) FindKey(key interface{}) string {
	for _, item := range *dict {
		if fmt.Sprint(item.Key) == fmt.Sprint(key) {
			return fmt.Sprint(item.Value)
		}
	}
	return ""
}

func (dict *Dictionary) FindValue(value interface{}) string {
	for _, item := range *dict {
		if fmt.Sprint(item.Value) == fmt.Sprint(value) {
			return fmt.Sprint(item.Key)
		}
	}
	return ""
}

func (dict *Dictionary) Push(key, value interface{}) {
	var new = append(*dict, Pair{Key: key, Value: value})
	*dict = new
}

func (dict *Dictionary) DeleteKey(key interface{}) {
	var new Dictionary
	for index, item := range *dict {
		if fmt.Sprint(item.Value) == fmt.Sprint(key) {
			new = remove(*dict, index)
			*dict = new
			break
		}
	}
}

func (dict *Dictionary) DeleteValue(value interface{}) {
	var new Dictionary
	for index, item := range *dict {
		if fmt.Sprint(item.Value) == fmt.Sprint(value) {
			new = remove(*dict, index)
			*dict = new
			break
		}
	}
}

func (dict *Dictionary) MapFromObject(ptr interface{}, key, value string) error {
	var new Dictionary
	var b, err = json.Marshal(ptr)
	if err != nil {
		return err
	}
	obj := gjson.Parse(string(b))
	for _, item := range obj.Array() {
		new = append(new, Pair{
			Key: item.Get(key).String(), Value: item.Get(value).String(),
		})

	}
	*dict = new
	return err
}

func remove(s Dictionary, i int) Dictionary {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func newElement(e *list.Element) *Pair {
	if e == nil {
		return nil
	}

	element := e.Value.(*orderedMapElement)

	return &Pair{
		element: e,
		Key:     element.key,
		Value:   element.value,
	}
}

// Next returns the message element, or nil if it finished.
func (e *Pair) Next() *Pair {
	return newElement(e.element.Next())
}

// Prev returns the previous element, or nil if it finished.
func (e *Pair) Prev() *Pair {
	return newElement(e.element.Prev())
}
