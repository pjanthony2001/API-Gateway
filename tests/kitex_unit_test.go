package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	client "github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"

	"github.com/cloudwego/kitex/pkg/loadbalance"

	"github.com/kitex-contrib/registry-nacos/resolver"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var kitexURLs = map[int32]string{1: "127.0.0.1:8888", 2: "127.0.0.1:8885"}

func JSONstring(msg string, additional string) (result string, err error) {
	data, err := json.Marshal(map[string]interface{}{"Msg": msg, "Additional": additional})
	result = string(data)
	return
}

func GenericServiceClient(kitex int32) (genericclient.Client, error) {

	g, err := GenericServiceConfig(kitex)

	if err != nil {
		return nil, err
	}

	return genericclient.NewClient(
		fmt.Sprintf("ExampleService%d", kitex),
		g,
		client.WithHostPorts(kitexURLs[kitex]),
	)
}

func GenericServiceClientWithResolver(kitex int32) (genericclient.Client, error) {

	g, err := GenericServiceConfig(kitex)

	if err != nil {
		return nil, err
	}

	resolverOption, err := NacosConfig()

	if err != nil {
		return nil, err
	}

	return genericclient.NewClient(
		fmt.Sprintf("ExampleService%d", kitex),
		g,
		resolverOption,
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()),
	)
}

func GenericServiceConfig(kitex int32) (generic.Generic, error) {

	p, err := generic.NewThriftFileProvider(fmt.Sprintf("../idl/example_service%d.thrift", kitex))
	if err != nil {
		return nil, err
	}

	return generic.JSONThriftGeneric(p)

}

func NacosConfig() (client.Option, error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig("127.0.0.1", 8848),
	}

	// the nacos client config
	cc := constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "info",
		Username:            "your-name",
		Password:            "your-password",
	}

	naco_client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	return client.WithResolver(resolver.NewNacosResolver(naco_client)), err
}
func TestKitex(t *testing.T) {

	var tests = []struct {
		name                   string
		sendMessage            string
		additionalInputs       string
		withResolver           bool
		withKitexServiceNumber int32
		wantMessage            string
	}{
		{"empty message should be empty", "", "", false, 1, "Parsed Message from Method 1: "},
		{"message should be message", "message", "", false, 1, "Parsed Message from Method 1: message"},
		{"empty with additional should be empty", "", "additional", false, 1, "Parsed Message from Method 1: "},
		{"msg with additional should be msg", "messages", "additional", false, 1, "Parsed Message from Method 1: messages"},
		{"empty message with resolver should be empty", "", "", true, 1, "Parsed Message from Method 1: "},
		{"message  with resolver should be message", "message", "", true, 1, "Parsed Message from Method 1: message"},
		{"empty with additional and resolver should be empty", "", "additional", true, 1, "Parsed Message from Method 1: "},
		{"msg with additional and resolver should be msg", "messages", "additional", true, 1, "Parsed Message from Method 1: messages"},
		{"empty message should be empty", "", "", false, 2, "Parsed Message from Method 1: "},
		{"message should be message", "message", "", false, 2, "Parsed Message from Method 1: message"},
		{"empty with additional should be empty", "", "additional", false, 2, "Parsed Message from Method 1: "},
		{"msg with additional should be msg", "messages", "additional", false, 2, "Parsed Message from Method 1: messages"},
		{"empty message with resolver should be empty", "", "", true, 2, "Parsed Message from Method 1: "},
		{"message  with resolver should be message", "message", "", true, 2, "Parsed Message from Method 1: message"},
		{"empty with additional and resolver should be empty", "", "additional", true, 2, "Parsed Message from Method 1: "},
		{"msg with additional and resolver should be msg", "messages", "additional", true, 2, "Parsed Message from Method 1: messages"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data, err := JSONstring(tt.sendMessage, tt.additionalInputs)

			if err != nil {
				t.Errorf("Error in parsing JSON: %s", err.Error())
			}

			var cli genericclient.Client

			if tt.withResolver {
				cli, err = GenericServiceClientWithResolver(tt.withKitexServiceNumber)
			} else {
				cli, err = GenericServiceClient(tt.withKitexServiceNumber)
			}

			if err != nil {
				t.Errorf("Error in creating client: %s", err.Error())
			}

			respRpc, err := cli.GenericCall(context.Background(), "ExampleMethod1", data)

			if err != nil {
				t.Errorf("Error in generic call: %s", err.Error())
			}

			t.Log(respRpc.(string))
			var sb map[string]interface{}
			json.Unmarshal([]byte(respRpc.(string)), &sb)

			if sb["Msg"] != tt.wantMessage {
				t.Errorf("got %s, want %s", sb, tt.wantMessage)
			}

			if len(sb) != 2 {
				t.Errorf("got %d, want %d", len(sb), 2)
			}
		})
	}
}
