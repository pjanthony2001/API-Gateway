package main

import (
	"context"
	"fmt"
	"net"

	jsoniter "github.com/json-iterator/go"

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/server/genericserver"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/registry-nacos/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func main() {

	// the nacos server config
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}

	// the nacos client config
	cc := constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log/kitex2",
		CacheDir:            "/tmp/nacos/cache/kitex2",
		LogLevel:            "info",
		// Username:            "your-name",
		// Password:            "your-password",
	}

	cli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	// Parse IDL with Local Files
	// YOUR_IDL_PATH thrift file path,eg: ./idl/example.thrift
	p, err := generic.NewThriftFileProvider("../idl/example_service1.thrift")
	if err != nil {
		panic(err)
	}
	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		panic(err)
	}

	add, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8885")
	if err != nil {
		panic(err)
	}

	svr := genericserver.NewServer(
		new(GenericServiceImpl),
		g,
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "ExampleService2"}),
		server.WithServiceAddr(add),
		server.WithRegistry(registry.NewNacosRegistry(cli)),
	)

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
	if method == "ExampleMethod1" {
		return ExampleMethod1(ctx, request)
	}

	return ctx, fmt.Errorf("Did not find available method: was given %s", method)
}

func ExampleMethod1(ctx context.Context, request interface{}) (response interface{}, err error) {
	var requestData map[string]interface{}
	err = jsoniter.Unmarshal([]byte(request.(string)), &requestData)
	for key, element := range requestData {
		fmt.Println("Key:", key, "=>", "Element:", element)
	}
	Message := "Parsed Message from Service 2 and Method 1: " + requestData["Msg"].(string)

	jsonstring, err := jsoniter.MarshalToString(
		struct {
			Msg string
		}{
			Msg: Message,
		})

	return jsonstring, err
}
