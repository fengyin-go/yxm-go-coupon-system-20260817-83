# 优惠券系统（coupon-system）

一个纯 Go 标准库（`net/http`，零第三方依赖）实现的优惠券管理后端服务。采用 `cmd/internal/pkg` 标准工程分层，可编译、可测试、可运行。

## 功能特性

- 优惠券模板管理：满减（threshold）/ 立减（fixed）两类券，支持有效期与发行总量控制。
- 发放批次：一张券可分多批次发放，实时跟踪已发放数量与剩余可发量。
- 用户券：从批次领取，内置状态机 `unused → used / expired`，券码全局唯一。
- 核销：核销用户券时自动生成使用记录，并同步更新券的已核销计数。
- 使用记录：按用户/券维度查询核销历史。

> 金额字段（`value` / `min_spend` / `amount`）统一使用 **int64「分」** 表示，避免浮点精度问题。

## 目录结构

```
origin/
├── cmd/server/main.go
├── internal/
│   ├── app/          # 依赖装配
│   ├── config/       # 环境变量配置
│   ├── model/        # 领域模型 + 校验 + 状态机
│   ├── store/        # 数据访问接口 + 内存实现
│   ├── service/      # 业务逻辑（领取/核销闭环）
│   └── handler/      # HTTP 路由 + 处理器
└── pkg/
    ├── httpx/        # 统一响应/分页/JSON
    ├── idgen/        # ID 与短码生成
    └── logger/       # 分级日志
```

## 运行

```bash
cd origin
go run ./cmd/server
# 或
go build -o coupon-server ./cmd/server && ./coupon-server
```

环境变量：

| 变量 | 默认 | 说明 |
|------|------|------|
| `PORT` | 8080 | 监听端口 |
| `ADDR` | `:8080` | 监听地址（优先于 PORT） |
| `MAX_PAGE_SIZE` | 100 | 分页最大每页条数 |
| `LOG_LEVEL` | info | debug / info / warn / error |

## API 一览

统一响应结构：`{"code":0,"message":"ok","data":...}`；错误码 400/404/409/500。

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/coupons` | 创建优惠券模板 |
| GET | `/api/coupons` | 列表（`type`/`status`/`keyword`/`page`/`size`） |
| GET | `/api/coupons/{id}` | 券详情 |
| PUT | `/api/coupons/{id}` | 更新券 |
| DELETE | `/api/coupons/{id}` | 删除券 |
| POST | `/api/batches` | 创建发放批次 |
| GET | `/api/batches` | 批次列表（`coupon_id`） |
| GET | `/api/batches/{id}` | 批次详情 |
| DELETE | `/api/batches/{id}` | 删除批次 |
| POST | `/api/user-coupons/issue` | 领取券（`user_id`+`batch_id`） |
| POST | `/api/user-coupons/{id}/use` | 核销券（`order_id`） |
| GET | `/api/user-coupons` | 用户券列表（`user_id`/`status`） |
| GET | `/api/user-coupons/{id}` | 用户券详情 |
| GET | `/api/usages` | 使用记录列表（`user_id`/`coupon_id`） |
| GET | `/api/usages/{id}` | 使用记录详情 |

## 业务闭环示例

```bash
# 1. 创建券（满 100 减 10，即 1000 分门槛、1000 分面值）
curl -s -X POST localhost:8080/api/coupons -d '{
  "name":"满100减10","type":"threshold","value":1000,"min_spend":10000,
  "total_count":100,"start_at":"2026-01-01T00:00:00Z","end_at":"2026-12-31T00:00:00Z"
}'

# 2. 创建发放批次（拿上面返回的券 id）
curl -s -X POST localhost:8080/api/batches -d '{"name":"双十一批次","coupon_id":"<coupon_id>","total_count":10}'

# 3. 领取券（拿批次 id）
curl -s -X POST localhost:8080/api/user-coupons/issue -d '{"user_id":"u1001","batch_id":"<batch_id>"}'

# 4. 核销券（拿用户券 id）
curl -s -X POST localhost:8080/api/user-coupons/<user_coupon_id>/use -d '{"order_id":"o20260816"}'
```

## 测试

```bash
go test ./...
```
