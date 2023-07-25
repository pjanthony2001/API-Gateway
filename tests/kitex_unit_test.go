package tests

import (
	"context"
	"encoding/json"
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

const kitexURL = "127.0.0.1:8888"

func JSONstring(msg string, additional string) (result string, err error) {
	data, err := json.Marshal(map[string]interface{}{"Msg": msg, "Additional": additional})
	result = string(data)
	return
}

func GenericServiceClient() (genericclient.Client, error) {

	g, err := GenericServiceConfig()

	if err != nil {
		return nil, err
	}

	return genericclient.NewClient(
		"ExampleService",
		g,
		client.WithHostPorts(kitexURL),
	)
}

func GenericServiceClientWithResolver() (genericclient.Client, error) {

	g, err := GenericServiceConfig()

	if err != nil {
		return nil, err
	}

	resolverOption, err := NacosConfig()

	if err != nil {
		return nil, err
	}

	return genericclient.NewClient(
		"ExampleService",
		g,
		resolverOption,
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()),
	)
}

func GenericServiceConfig() (generic.Generic, error) {

	p, err := generic.NewThriftFileProvider("../idl/example_service.thrift")
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
		name             string
		wantMsg          string
		additionalInputs string
		withResolver     bool
	}{
		{"empty message should be empty", "", "", false},
		{"message should be message", "message", "", false},
		{"empty with additional should be empty", "", "additional", false},
		{"msg with additional should be msg", "messages", "additional", false},
		{"empty message with resolver should be empty", "", "", true},
		{"message  with resolver should be message", "message", "", true},
		{"empty with additional and resolver should be empty", "", "additional", true},
		{"msg with additional and resolver should be msg", "messages", "additional", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data, err := JSONstring(tt.wantMsg, tt.additionalInputs)

			if err != nil {
				t.Errorf("Error in parsing JSON: %s", err.Error())
			}

			var cli genericclient.Client

			if tt.withResolver {
				cli, err = GenericServiceClientWithResolver()
			} else {
				cli, err = GenericServiceClient()
			}

			if err != nil {
				t.Errorf("Error in creating client: %s", err.Error())
			}

			respRpc, err := cli.GenericCall(context.Background(), "ExampleMethod", data)

			if err != nil {
				t.Errorf("Error in generic call: %s", err.Error())
			}

			t.Log(respRpc.(string))
			var sb map[string]interface{}
			json.Unmarshal([]byte(respRpc.(string)), &sb)

			if sb["Msg"] != tt.wantMsg {
				t.Errorf("got %s, want %s", sb, tt.wantMsg)
			}

			if len(sb) != 2 {
				t.Errorf("got %d, want %d", len(sb), 2)
			}
		})
	}
}
