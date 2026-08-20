# Bug Reproduction

缓存保存和返回 DNS packet 时复用了调用方的字节切片。修改 Put 输入或 Get 返回值都会污染后续读取。

复现：运行 collection.json 中的 resolver_runtime 定向测试。埋错基线会出现缓存响应被调用方 buffer 改写。
