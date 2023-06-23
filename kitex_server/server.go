package main

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/server/genericserver"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/registry-nacos/example/hello/kitex_gen/api"
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
		TimeoutMs:           100000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "info",
		// more ...
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
	p, err := generic.NewThriftFileProvider("../idl/example_service.thrift")
	if err != nil {
		panic(err)
	}
	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		panic(err)
	}
	svr := genericserver.NewServer(
		new(GenericServiceImpl),
		g,
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "Hello"}),
		server.WithRegistry(registry.NewNacosRegistry(cli)))

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

type HelloImpl struct{}

func (h *HelloImpl) Echo(_ context.Context, req *api.Request) (resp *api.Response, err error) {
	resp = &api.Response{
		Message: req.Message,
	}
	return
}
