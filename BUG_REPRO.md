# Bug Reproduction

DNSSEC 签名入口没有统一拒绝不存在、退役、算法为零或资料不完整的 key，异常输入可能进入签名流程并表现为不稳定成功。

复现：运行 collection.json 中的四条 dnssec 定向测试。埋错基线会错误接受缺失 key 或非法签名参数。
