package nearclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Request struct {
	Version string          `json:"jsonrpc,omitempty"`
	Id      string          `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type FunctionCallParams struct {
	Request_type string `json:"request_type"`
	Finality     string `json:"finality"`
	Account_id   string `json:"account_id"`
	Method_name  string `json:"method_name"`
	Args_base64  string `json:"args_base64"`
}

type Response struct {
	Version string          `json:"jsonrpc,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Id      string          `json:"id,omitempty"`
}

type FunctionCallResponse struct {
	Result      json.RawMessage `json:"result,omitempty"`
	Error       string          `json:"error,omitempty"`
	Logs        []string        `json:"logs,omitempty"`
	BlockHeight int64           `json:"block_height,omitempty"`
	BlockHash   string          `json:"block_hash,omitempty"`
}

type Client struct {
	URL string
}

func RequestJSON(method string, params interface{}) ([]byte, error) {
	enc_params, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	request := Request{
		Version: "2.0",
		Id:      "dontcare",
		Method:  method,
		Params:  enc_params,
	}

	enc, err_r := json.Marshal(request)
	if err != nil {
		return nil, err_r
	}

	return enc, nil
}

func (c *Client) MakeRequest(data []byte) ([]byte, error) {
	body := bytes.NewBuffer(data)
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post(c.URL, "application/json", body)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	//Read the response body
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	var resp_dec Response

	if err := json.Unmarshal(resp_body, &resp_dec); err != nil {
		return nil, err
	}

	return resp_dec.Result, nil
}

func (c *Client) FunctionCall(account_id string, method_name string, args string) (*FunctionCallResponse, error) {
	params := FunctionCallParams{
		Request_type: "call_function",
		Finality:     "final",
		Account_id:   account_id,
		Method_name:  method_name,
		Args_base64:  args,
	}

	data, err := RequestJSON("query", params)
	if err != nil {
		return nil, err
	}

	resp, err := c.MakeRequest(data)
	if err != nil {
		return nil, err
	}

	var resp_dec FunctionCallResponse

	if err := json.Unmarshal(resp, &resp_dec); err != nil {
		return nil, err
	}

	return &resp_dec, nil
}

//Returns general status of current validator nodes.
func (c *Client) Status() ([]byte, error) {
	params := make([]string, 0)

	data, err := RequestJSON("status", params)
	if err != nil {
		return nil, err
	}

	resp, err := c.MakeRequest(data)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
