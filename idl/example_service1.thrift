include "base.thrift"
namespace go kitex.test.server

struct ExampleReq {
    1: required string Msg,
    255: base.Base Base,
}
struct ExampleResp {
    1: required string Msg,
    255: base.BaseResp BaseResp,
}

service ExampleService1 {
    ExampleResp ExampleMethod1(1: ExampleReq req),
    ExampleResp ExampleMethod2(1: ExampleReq req),
}


