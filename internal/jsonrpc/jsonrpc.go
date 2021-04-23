package jsonrpc

import (
	"encoding/json"
	"strconv"
)

/* error example
{
        "jsonrpc":      "2.0",
        "error":        {
                "code": -32000,
                "message":      "Execution Error"
        },
        "id":   "1"
}
*/

// Request w/o params
type Request struct {
	JSONRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
}

// RequestWithParams w/ params
type RequestWithParams struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// Response is exported
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// Error is exported
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Error returns error Code + Message
func (err *Error) Error() string {
	return strconv.Itoa(err.Code) + " " + err.Message
}

// NewRequest returns new *Request without params
func NewRequest(id string, method string) *Request {

	req := &Request{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
	}
	return req
}

// NewRequestWithParams *RequestWithParams struct
func NewRequestWithParams(id string, method string, params interface{}) *RequestWithParams {

	req := &RequestWithParams{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}
	return req
}

// Buffer method returns JSON encoded data
func (request *Request) Buffer() ([]byte, error) {
	return json.Marshal(request)
}

// Parse methos prases JSON and stores result in a struct
func (reply *Response) Parse(bytes []byte) error {
	return json.Unmarshal(bytes, reply)
}

// IsError checks whether response contains error
func (reply *Response) IsError() bool {
	return reply.Error != nil
}
