# Kong as a Shield Proxy

This is the opinionated POC implementation of using Kong as a shield proxy. The setup is using kubernetes and this POC is using [k3d](https://k3d.io/v5.7.3/). This is using Kong db-less mode so the routing configuration needs to be passed in the helm values.

## Steps

1. Init kong and run shield service

Initialize kong

```
./init.sh
```

Run Shield Server

```bash
go run ./../.. server start
```

Kong will be running with the routing config passed in the helm values.

2. Prepare Kong

Now, you could port forward kong pod. Kong admin service could be exposed from port `8001`

```bash
kubectl port-forward kong-kong-84b7b977c-fdsl5 8001:8001
```

Verify that routing is configured with this curl

```
➜ curl -i -X GET http://localhost:8001/services | jq .
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   784  100   784    0     0  40805      0 --:--:-- --:--:-- --:--:-- 41263
jq: parse error: Invalid numeric literal at line 1, column 9

~
➜ curl -X GET http://localhost:8001/services | jq .
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   784  100   784    0     0  49968      0 --:--:-- --:--:-- --:--:-- 52266
{
  "next": null,
  "data": [
    {
      "tls_verify": null,
      "enabled": true,
      "protocol": "http",
      "created_at": 1724337791,
      "updated_at": 1724337791,
      "client_certificate": null,
      "write_timeout": 60000,
      "name": "shield-http",
      "id": "31dd9b4f-dca8-5d81-9cb8-07f7fb13a678",
      "tags": null,
      "ca_certificates": null,
      "tls_verify_depth": null,
      "read_timeout": 60000,
      "port": 8000,
      "connect_timeout": 60000,
      "host": "host.k3d.internal",
      "retries": 5,
      "path": null
    },
    {
      "tls_verify": null,
      "enabled": true,
      "protocol": "http",
      "created_at": 1724337791,
      "updated_at": 1724337791,
      "client_certificate": null,
      "write_timeout": 60000,
      "name": "shield-grpc",
      "id": "9d526ec4-b1b7-554a-9f43-4508340e61be",
      "tags": null,
      "ca_certificates": null,
      "tls_verify_depth": null,
      "read_timeout": 60000,
      "port": 8081,
      "connect_timeout": 60000,
      "host": "host.k3d.internal",
      "retries": 5,
      "path": null
    }
  ]
}
```

3. Test HTTP and GRPC call

Test calling shield HTTP API through Kong.

```
➜ curl http://localhost:8888/admin/admin/ping
{"status":"SERVING"}
```

Test calling shield GRPC API through Kong

```
➜ grpcurl -vv -plaintext localhost:8888 gotocompany.shield.v1beta1.ShieldService/ListProjects


Resolved method descriptor:
rpc ListProjects ( .gotocompany.shield.v1beta1.ListProjectsRequest ) returns ( .gotocompany.shield.v1beta1.ListProjectsResponse ) {
  option (.google.api.http) = { get: "/v1beta1/projects" };
  option (.grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = { tags: [ "Project" ], summary: "Get all Project" };
}

Request metadata to send:
(empty)

Response headers received:
content-type: application/grpc
date: Mon, 26 Aug 2024 09:01:45 GMT
server: openresty
via: kong/3.6.1
x-kong-proxy-latency: 0
x-kong-request-id: 7a6156070b05d13438cea24d20585315
x-kong-upstream-latency: 32

Estimated response size: 2308 bytes

Response contents:
<< Body >>
```
