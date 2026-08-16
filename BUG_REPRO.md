# BUG-003 复现

## 现象

查询不存在用户的排名时，服务崩溃（panic）而不是返回错误。

## 触发步骤

1. GET /rank/nobody（nobody 不存在）
2. 服务 panic（nil pointer dereference）

## 根因

internal/api/handler.go 的 rank handler 调用 GetUserEntry 后直接解引用 entry.UserID，不存在的用户返回 nil，导致 panic。

## 错误信息

```
runtime error: invalid memory address or nil pointer dereference
```
