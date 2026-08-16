# BUG-002 复现

## 现象

更新分数接口接受负分。传入负数（如 -50）会成功写入排行榜，导致出现负分用户，污染排名。

## 触发步骤

1. 调用 UpdateScore("alice", -50)
2. 返回 nil（无错误）
3. alice 出现在排行榜中，分数为 -50

## 根因

internal/leaderboard/leaderboard.go 的 UpdateScore 方法不校验 score 是否为负，直接写入 scores map。

## 错误信息

```
expected error for negative score, got nil
```
