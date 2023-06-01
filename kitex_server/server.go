package main

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/server/genericserver"
)

func main() {
	// Parse IDL with Local Files
	// YOUR_IDL_PATH thrift file path,eg: ./idl/example.thrift
	p, err := generic.NewThriftFileProvider("../idl/example_service.thrift")
	if err != nil {
		panic(err)
	}
	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		panic(err)
	}
	svr := genericserver.NewServer(new(GenericServiceImpl), g)
	if err != nil {
		panic(err)
	}
	err = svr.Run()
	if err != nil {
		fmt.Println("this is the error runtime")
		panic(err)
	}
	// resp is a JSON string
}

type GenericServiceImpl struct {
}

func (g *GenericServiceImpl) GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
	var requestData map[string]interface{}
	err = jsoniter.Unmarshal([]byte(request.(string)), &requestData)
	for key, element := range requestData {
		fmt.Println("Key:", key, "=>", "Element:", element)
	}
	Message := requestData["Msg"].(string)
	return "{\"Msg\" : \"" + Message + "\"}", err
}
