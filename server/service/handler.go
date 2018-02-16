package service

import (
	"encoding/json"
	"fmt"

	"github.com/Ullaakut/cameradar"
	"github.com/Ullaakut/cameradar/server/adaptor/jsonrpc2"
	v "gopkg.in/go-playground/validator.v9"
)

func (c *Cameradar) handleRequest(message string) {
	var ret []cmrdr.Stream
	var JSONRPCErr jsonrpc2.Error

	var request jsonrpc2.Request
	err := json.Unmarshal([]byte(message), &request)
	if err != nil {
		JSONRPCErr = jsonrpc2.Error{
			Code:    jsonrpc2.ParseError,
			Message: jsonrpc2.ParseErrorMessage,
			Data:    err.Error(),
		}
		c.respondToClient(ret, "", JSONRPCErr)
		return
	}

	validate := v.New()
	err = validate.Struct(request)
	if err != nil {
		JSONRPCErr = jsonrpc2.Error{
			Code:    jsonrpc2.InvalidRequest,
			Message: jsonrpc2.InvalidRequestMessage,
			Data:    err.Error(),
		}
		c.respondToClient(ret, request.ID, JSONRPCErr)
		return
	}

	c.SetOptions(request.Params)

	switch request.Method {
	case "discover":
		ret, err = c.Discover()
	case "attack_credentials":
		ret, err = c.AttackCredentials()
	case "attack_route":
		ret, err = c.AttackRoute()
	case "discover_and_attack":
		c.discoverAndAttack(request.ID)
		return
	default:
		JSONRPCErr = jsonrpc2.Error{
			Code:    jsonrpc2.MethodNotFound,
			Message: jsonrpc2.MethodNotFoundMessage,
			Data:    "method" + request.Method + "not found",
		}
	}
	if err != nil {
		JSONRPCErr = jsonrpc2.Error{
			Code:    jsonrpc2.InternalError,
			Message: jsonrpc2.InternalErrorMessage,
			Data:    err.Error(),
		}
	}

	c.respondToClient(ret, request.ID, JSONRPCErr)
}

func (c *Cameradar) discoverAndAttack(ID string) {
	var JSONRPCErr jsonrpc2.Error

	streams, err := c.Discover()
	if err != nil {
		c.respondToClient(streams, ID, jsonrpc2.Error{
			Code:    jsonrpc2.InternalError,
			Message: jsonrpc2.InternalErrorMessage,
			Data:    err.Error(),
		})
		return
	}
	c.respondToClient(streams, ID, JSONRPCErr)

	streams, err = c.AttackRoute()
	if err != nil {
		c.respondToClient(streams, ID, jsonrpc2.Error{
			Code:    jsonrpc2.InternalError,
			Message: jsonrpc2.InternalErrorMessage,
			Data:    err.Error(),
		})
		return
	}
	c.respondToClient(streams, ID, JSONRPCErr)

	streams, err = c.AttackCredentials()
	if err != nil {
		c.respondToClient(streams, ID, jsonrpc2.Error{
			Code:    jsonrpc2.InternalError,
			Message: jsonrpc2.InternalErrorMessage,
			Data:    err.Error(),
		})
		return
	}
	c.respondToClient(streams, ID, JSONRPCErr)

	for _, stream := range streams {
		if stream.RouteFound == false {
			streams, err = c.AttackCredentials()
			if err != nil {
				c.respondToClient(streams, ID, jsonrpc2.Error{
					Code:    jsonrpc2.InternalError,
					Message: jsonrpc2.InternalErrorMessage,
					Data:    err.Error(),
				})
				return
			}
			c.respondToClient(streams, ID, JSONRPCErr)
			return
		}
	}
}

func (c *Cameradar) respondToClient(result []cmrdr.Stream, ID string, JSONRPCErr jsonrpc2.Error) {
	r := jsonrpc2.Response{
		JSONRPC: "2.0",
		Result:  result,
		Error:   JSONRPCErr,
		ID:      ID,
	}

	response, err := json.Marshal(r)
	if err != nil {
		c.toClient <- "{\"jsonrpc\": \"2.0\",\"result\":null,\"error\":{\"code\":" + fmt.Sprint(jsonrpc2.InternalError) + ",\"" + jsonrpc2.InternalErrorMessage + "\",\"data\":\"could not marshal response\"},\"id\":\"" + ID + "\"}"
	}

	c.toClient <- string(response)
}
