# BUG_REPRO

## Bug 是什么

优惠券新增了每人限领数量 `user_limit` 字段；当一张券未设置该字段（值为 nil）时，领取流程会直接解引用这个 nil 指针。

## 如何触发

1. 创建一个不含 `user_limit` 的优惠券模板。
2. 基于该券创建一个发放批次。
3. 调用 `POST /api/user-coupons/issue` 领取该批次，或直接调用 `Service.Issue`。

## 错误信息

```text
panic: runtime error: invalid memory address or nil pointer dereference
```
