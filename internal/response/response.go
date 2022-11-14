package response

import "encoding/json"

type Response struct {
	Error string
	Data  interface{}
}

func (r *Response) ToByte() []byte {
	data, err := json.Marshal(r)
	if err != nil {
		return nil
	}
	return data
}
