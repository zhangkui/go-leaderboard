# BUG-001 复现

## 现象

查询用户排名时，排名编号返回 0-based 而非 1-based。分数最高的用户（榜首）排名返回 0，而非 1；所有用户排名整体少 1。

## 触发步骤

1. 更新若干用户分数（如 alice=100, bob=200）
2. 查询榜首用户 bob 的排名：`GetRank("bob")`
3. 返回 0（期望 1）

## 根因

`internal/leaderboard/leaderboard.go` 的 `GetRank` 方法返回切片下标 `i`（0-based），而文档注释和契约要求 1-based 排名。

## 错误信息

```
expected rank 1 for top user, got 0
expected rank 2 for second user, got 1
```
