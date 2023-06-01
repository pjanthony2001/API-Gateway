namespace go api

struct Response {
    1 : string Message
    2 : i32 Flag
}

struct Request {
    1 : string Message
    2 : i32 Flag
}

service Echo {
    Response echo (1 : Request req) (api.get="echo/query");
}