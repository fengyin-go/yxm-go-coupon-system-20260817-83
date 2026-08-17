# BUG_REPRO

## Bug 是什么

领取用户券的“检查剩余数量 → 创建用户券 → 更新批次和券计数”不是原子操作，并发领取同一个批次时可能超发。

## 如何触发

并发调用 `POST /api/user-coupons/issue` 领取一个 `total_count` 很小的批次，或多个 goroutine 同时调用 `Service.Issue`。

## 错误信息

并发请求最终返回的成功数量可能超过批次 `total_count`；使用 `-race` 运行时可观察到 data race 报告。
