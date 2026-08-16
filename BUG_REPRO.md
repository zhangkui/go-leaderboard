# BUG-004 复现

## 现象

分页查询排行榜时结果错位。查第 1 页（page=1, size=10）返回的不是前 10 名，而是第 11-20 名。

## 触发步骤

1. 添加 25 个用户（分数 1000 到 976）
2. GET /leaderboard?page=1&size=10
3. 返回第 11-20 名（score 990-981），而非第 1-10 名（score 1000-991）

## 根因

internal/api/handler.go:85 的 leaderboard handler 计算 offset 时用了 `page * size`（应为 `(page-1) * size`），page=1 时 offset=10，跳过第一页。

## 错误信息

```
expected top score 1000 first on page 1, got 990
```
