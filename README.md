# Authoritative DNS and Anycast Control Plane

纯 Go 1.23 实现的权威 DNS 控制平面示例。服务提供 zone/record draft 发布、不可变快照、策略解析、健康多数派、Anycast 路由意图、DNS wire 查询和 DoH/UDP/TCP 接入。

## Run

```sh
go run ./cmd/control-plane
curl -s localhost:8080/healthz
curl -s -X POST localhost:8080/api/v1/zones -d '{"name":"example.com."}'
```

默认仅使用内存存储，PostgreSQL/Redis/对象存储通过 `internal/repository` 接口替换。生产部署清单位于 `deployments/`，OpenAPI 在 `api/openapi.yaml`，gRPC 原型在 `api/proto/control_plane.proto`。

## Components

`zone_domain` 聚合 zone 和记录版本；`snapshot_release` 负责校验、serial、发布和回滚；`policy_engine` 解释地域/ASN/ECS/权重规则；`probe_domain` 执行滑动窗口多数派；`resolver_runtime` 提供 wire 编解码、缓存和传输；`dnssec` 与 `anycast_control` 管理密钥与宣告意图；`workers` 负责后台收敛。

