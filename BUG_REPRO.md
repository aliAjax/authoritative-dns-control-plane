# Bug Reproduction

Worker 的 Start 不是幂等的，多个并发启动会创建重复循环；Stop 在未启动或重复调用时没有统一生命周期保护，可能关闭同一通道并触发 panic。

复现：运行 collection.json 中的两条 operations 定向测试并启用 race 检查。埋错基线会出现 `close of closed channel` 或数据竞争。
