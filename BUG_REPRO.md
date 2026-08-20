# Bug Reproduction

HTTP 客户端重试汇总错误时丢失 transport/read error 的 sentinel identity，取消请求仍可能进入下一轮重试，状态错误也缺少最终 attempt 上下文。

复现：运行 collection.json 中的四条 protocol 定向测试。埋错基线无法通过错误链、取消传播和状态重试上下文断言。
