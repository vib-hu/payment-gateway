package formatters

import "encoding/json"

type JsonFormatter struct{}

func NewJsonFormatter() *JsonFormatter {
	return &JsonFormatter{}
}

func (f *JsonFormatter) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (f *JsonFormatter) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
