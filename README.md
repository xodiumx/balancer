# URL Balancer

## Description

- Example of balancing service
- when sending a request to grpc-server (balancer) every 10 requests sends redirect url of original server,
  otherwise CDN url is sent depending on zone `s1,s2,s3,n`.
- different subdomains s1,s2,s3,etc. are cached
  and redirect to the original url comes for each 10 requests each one separately
- 2 containers, a balancer and a gateway for direct request are configured in the service

## Install

- Clone repository

- Help command

```shell
make help
```

- create `.env` file in `./dev/docker`

- Run app in docker

```shell
make app
```

- Start logs:

```Text
2025-06-26 23:09:40 balancer-1  | 2025/06/26 20:09:40 [config] CDN_HOST=cdn.host.original.lol.322 DEBUG=true FREQUENCY=10
2025-06-26 23:09:40 balancer-1  | 2025-06-26T20:09:40.848Z      INFO    server/server.go:60     🚀 gRPC server started       {"addr": ":50051"}
2025-06-26 23:09:40 gateway     | 2025/06/26 20:09:40 🌐 Gateway started on :8080
```

## Use cases

### First

- For test balancer u can use [grpcurl](https://github.com/fullstorydev/grpcurl) for make requests on app

```shell
grpcurl -plaintext -d '{"video":"http://s2.origin-cluster/video/123/xcg2djHckad.m3u8"}' \
  localhost:50051 balancer.VideoBalancer/GetRedirect
```

- Output:

```text
{
  "redirectUrl": "http://cdn.host.original.lol.322/s2/video/123/xcg2djHckad.m3u8"
}
```

- Every 10 requests redirect to the original url:

```shell
grpcurl -plaintext -d '{"video":"http://s2.origin-cluster/video/123/xcg2djHckad.m3u8"}' \
  localhost:50051 balancer.VideoBalancer/GetRedirect
```

- Output:

```text
{
  "redirectUrl": "http://s2.origin-cluster/video/123/xcg2djHckad.m3u8"
}
```

### Second

- [ghz](https://github.com/bojand/ghz) is used for load testing.

- Run command:

```shell
make nt
```

- Metrics mac m1

```text
Summary:
  Count:        1000000
  Total:        56.37 s
  Slowest:      82.79 ms
  Fastest:      0.62 ms
  Average:      5.16 ms
  Requests/sec: 17738.45

Response time histogram:
  0.622  [1]      |
  8.839  [943533] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  17.056 [49309]  |∎∎
  25.273 [3572]   |
  33.490 [1600]   |
  41.708 [814]    |
  49.925 [412]    |
  58.142 [352]    |
  66.359 [302]    |
  74.576 [94]     |
  82.793 [11]     |

Latency distribution:
  10 % in 2.86 ms
  25 % in 3.54 ms
  50 % in 4.56 ms
  75 % in 6.05 ms
  90 % in 7.77 ms
  95 % in 9.08 ms
  99 % in 14.01 ms

Status code distribution:
  [OK]   1000000 responses
```

- Logs

```text
...
2025-06-26 18:36:12 balancer-1  | 2025-06-26T15:36:12.712Z      INFO    handler/handler.go:55   Request data{"original URL": "http://s1.origin-cluster/video/123/xcg2djHckad.m3u8", "redirect URL": "http://cdn.host.original.lol.322/s1/video/123/xcg2djHckad.m3u8", "request_number": 1000006}
2025-06-26 18:36:12 balancer-1  | 2025-06-26T15:36:12.712Z      INFO    zap/options.go:212      finished unary call with code OK    {"grpc.start_time": "2025-06-26T15:36:12Z", "grpc.request.deadline": "2025-06-26T15:36:32Z", "system": "grpc", "span.kind": "server", "grpc.service": "balancer.VideoBalancer", "grpc.method": "GetRedirect", "grpc.code": "OK", "grpc.time_ms": 0.005}
2025-06-26 18:36:12 balancer-1  | 2025-06-26T15:36:12.712Z      INFO    handler/handler.go:55   Request data{"original URL": "http://s1.origin-cluster/video/123/xcg2djHckad.m3u8", "redirect URL": "http://cdn.host.original.lol.322/s1/video/123/xcg2djHckad.m3u8", "request_number": 1000007}
2025-06-26 18:36:12 balancer-1  | 2025-06-26T15:36:12.712Z      INFO    zap/options.go:212      finished unary call with code OK    {"grpc.start_time": "2025-06-26T15:36:12Z", "grpc.request.deadline": "2025-06-26T15:36:32Z", "system": "grpc", "span.kind": "server", "grpc.service": "balancer.VideoBalancer", "grpc.method": "GetRedirect", "grpc.code": "OK", "grpc.time_ms": 0.006}
2025-06-26 18:36:12 balancer-1  | 2025-06-26T15:36:12.712Z      INFO    handler/handler.go:55   Request data{"original URL": "http://s1.origin-cluster/video/123/xcg2djHckad.m3u8", "redirect URL": "http://cdn.host.original.lol.322/s1/video/123/xcg2djHckad.m3u8", "request_number": 1000008}
2025-06-26 18:36:12 balancer-1  | 2025-06-26T15:36:12.712Z      INFO    zap/options.go:212      finished unary call with code OK    {"grpc.start_time": "2025-06-26T15:36:12Z", "grpc.request.deadline": "2025-06-26T15:36:32Z", "system": "grpc", "span.kind": "server", "grpc.service": "balancer.VideoBalancer", "grpc.method": "GetRedirect", "grpc.code": "OK", "grpc.time_ms": 0.005}
2025-06-26 18:36:12 balancer-1  | 2025-06-26T15:36:12.712Z      INFO    handler/handler.go:55   Request data{"original URL": "http://s1.origin-cluster/video/123/xcg2djHckad.m3u8", "redirect URL": "http://cdn.host.original.lol.322/s1/video/123/xcg2djHckad.m3u8", "request_number": 1000009}
2025-06-26 18:36:12 balancer-1  | 2025-06-26T15:36:12.712Z      INFO    zap/options.go:212      finished unary call with code OK    {"grpc.start_time": "2025-06-26T15:36:12Z", "grpc.request.deadline": "2025-06-26T15:36:32Z", "system": "grpc", "span.kind": "server", "grpc.service": "balancer.VideoBalancer", "grpc.method": "GetRedirect", "grpc.code": "OK", "grpc.time_ms": 0.018}
2025-06-26 18:36:12 balancer-1  | 2025-06-26T15:36:12.712Z      WARN    handler/handler.go:31   Request data{"Redirect url": "http://s1.origin-cluster/video/123/xcg2djHckad.m3u8", "request_number": 1000010}
2025-06-26 18:36:12 balancer-1  | 2025-06-26T15:36:12.712Z      INFO    zap/options.go:212      finished unary call with code OK    {"grpc.start_time": "2025-06-26T15:36:12Z", "grpc.request.deadline": "2025-06-26T15:36:32Z", "system": "grpc", "span.kind": "server", "grpc.service": "balancer.VideoBalancer", "grpc.method": "GetRedirect", "grpc.code": "OK", "grpc.time_ms": 0.027}
```

### Third

- Test redirect via gateway, run:

```shell
curl -X POST http://localhost:8080/watch -v \
  -H "Content-Type: application/json" \
  -d '{"video": "http://s1.origin-cluster/video/123/vid.m3u8"}'
```

- Output:

```text
Note: Unnecessary use of -X or --request, POST is already inferred.
* Host localhost:8080 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> POST /watch HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.7.1
> Accept: */*
> Content-Type: application/json
> Content-Length: 56
>
* upload completely sent off: 56 bytes
< HTTP/1.1 302 Found
< Location: http://cdn.host.original.lol.322/s1/video/123/vid.m3u8
< Date: Thu, 26 Jun 2025 17:09:37 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
```
