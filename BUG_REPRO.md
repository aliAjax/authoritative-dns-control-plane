# Bug Reproduction

策略评估把缺失 ID 或非法 CIDR 当成普通不匹配，校验错误没有沿调用链返回给调用方。

复现：运行 collection.json 中的两条 policy_engine 定向测试。埋错基线不会明确报告非法规则。
