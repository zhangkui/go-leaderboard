# go-leaderboard

## 项目说明

一个基于内存的实时排行榜服务，支持分数更新、排名查询、TopN 获取和分页查询，提供 HTTP 接口供游戏或社交应用使用。

## 标准命令

go build ./...   # 编译
go test ./...    # 测试
go run ./cmd     # 启动（监听 :8080）

## 使用方式

启动后通过 HTTP 接口操作：

- `POST /score` — 更新分数，body: `{"user_id":"alice","score":100}`
- `GET /rank/:user` — 查询用户排名，如 `/rank/alice`
- `GET /topn?n=10` — 获取前 N 名
- `GET /leaderboard?page=1&size=10` — 分页查询

数据全部存内存，无外部依赖。
