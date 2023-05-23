package decoder

import jsoniter "github.com/json-iterator/go"

type JsonIter struct{}

func New() *JsonIter {
	return &JsonIter{}
}

func (j *JsonIter) Marshal(v interface{}) ([]byte, error) {
	return jsoniter.Marshal(v)
}

func (j *JsonIter) Unmarshal(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}
