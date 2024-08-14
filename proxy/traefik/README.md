# Traefik as a Shield Proxy

This is the opinionated POC implementation of using Traefik as a shield proxy. The setup is using kubernetes and this POC is using [k3d](https://k3d.io/v5.7.3/).

## Steps

1. Init traefik and run shield service

Initialize traefik

```
./init.sh
```

Run Shield Server

```bash
go run ./../.. server start
```

Traefik will pick up config `static.yml` and spin up a pod. Siren server is running and expose path `/admin/traefikrule` and traefik pod will periodically fetches dynamic config from the endpoint.

2. Prepare Traefik

Now, you could port forward traefik pod to the host port in port `8888`

```bash
kubectl port-forward traefik-76644c5668-fsstc -n default 8888:8888
```

You also could expose traefik admin UI in port `9001`

```bash
kubectl port-forward traefik-76644c5668-fsstc -n default 9001:9001
```

3. Test HTTP and GRPC call

Test calling shield HTTP API through traefik

```
➜ curl -i -X GET http://localhost:8888/admin/traefikrule
HTTP/1.1 200 OK
Content-Length: 1233
Content-Type: text/plain; charset=utf-8

Date: Thu, 22 Aug 2024 14:18:42 GMT

<< Body >>
```

Test calling shield GRPC API through traefik

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
date: Thu, 22 Aug 2024 14:24:03 GMT

Estimated response size: 2308 bytes

Response contents:
<< Body >>

```

Check Shield logs, Shield API `/admin/traefikcheck` is being called since we are using `forward-auth` middleware plugin. In this case, authorization is being offloaded to the external component.
