# BENZHI_README

## 项目说明

- 项目：zhangkui/go-leaderboard
- 项目用途：in-memory real-time leaderboard service
- Go 工具链：`golang:1.22`
- 前端工具链：无

## 标准构建、运行和测试命令

进入容器后执行：

```bash
cd '/app' && GOTOOLCHAIN=local go build ./...
cd '/app' && GOTOOLCHAIN=local go run ./cmd
cd '/app' && GOTOOLCHAIN=local go test ./...
```

## Docker 构建和进入容器

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh go-leaderboard-bug2-amd64 linux/amd64
./build_benzhi_docker.sh go-leaderboard-bug2-arm64 linux/arm64
docker run -it go-leaderboard-bug2-amd64:latest
docker run -it --platform linux/arm64 go-leaderboard-bug2-arm64:latest
```

## 题目验证命令

1. 预期退出码 0：`go test ./internal/leaderboard -run "TestPrivateUpdateScoreRejectsNegative" -count=1 -v`
2. 预期退出码 0：`go test -buildvcs=false -count=1 ./...`

## Bug 复现

Bug 现象、触发步骤和完整错误信息见 `BUG_REPRO.md`。
