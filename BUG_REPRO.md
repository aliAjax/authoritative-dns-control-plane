# Bug Reproduction

AXFR 编码在错误分支直接返回，跳过缓冲 writer 的最终 Flush；空记录时头部可能留在缓冲区，底层 writer 的 flush 错误也会丢失。

复现：运行 collection.json 中的四条 protocol 定向测试。埋错基线无法观察完整头部或正确收到 writer 错误。
