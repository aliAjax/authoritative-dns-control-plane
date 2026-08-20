# Bug Reproduction

请求取消后，解析和写操作仍沿用不带取消语义的下游 context，客户端已经断开时后台工作不会及时停止。

复现：运行 collection.json 中的 service 定向测试。埋错基线会报告取消信号没有传到下游。
