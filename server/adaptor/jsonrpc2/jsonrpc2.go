// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jsonrpc2

import "github.com/EtixLabs/cameradar"

// http://www.jsonrpc.org/specification
const (
	ParseError     = -32700 // Invalid JSON was received by the server. An error occurred on the server while parsing the JSON text.
	InvalidRequest = -32600 // The JSON sent is not a valid Request object.
	MethodNotFound = -32601 // The method does not exist / is not available.
	InvalidParams  = -32602 // Invalid method parameter(s).
	InternalError  = -32603 // Internal JSON-RPC error.
)

const (
	// ParseErrorMessage is the message associated with the ParseError error
	ParseErrorMessage = "Parse error"
	// InvalidRequestMessage is the message associated with the InvalidRequest error
	InvalidRequestMessage = "Invalid Request"
	// MethodNotFoundMessage is the message associated with the MethodNotFound error
	MethodNotFoundMessage = "Method not found"
	// InvalidParamsMessage is the message associated with the InvalidParams error
	InvalidParamsMessage = "Invalid params"
	// InternalErrorMessage is the message associated with the InternalError error
	InternalErrorMessage = "Internal error"
)

// Request represents a JSONRPC request
type Request struct {
	JSONRPC string        `json:"jsonrpc" validate:"eq=2.0"`
	Method  string        `json:"method" validate:"required"`
	Params  cmrdr.Options `json:"params" validate:"required"`
	ID      string        `json:"id"`
}

// Response represents a JSONRPC response
type Response struct {
	JSONRPC string         `json:"jsonrpc" validate:"eq=2.0"`
	Result  []cmrdr.Stream `json:"result"`
	Error   Error          `json:"error"`
	ID      string         `json:"id"`
}

// Error represents a JSONRPC response's error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}
