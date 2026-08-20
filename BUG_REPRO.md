# Bug Reproduction

Anycast 规划器只过滤严格更旧的 fencing token，相同 token 的 intent 会重复提交，状态更新交错时旧 token 也可能重新生成意图。

复现：运行 collection.json 中的两条 anycast_control 定向测试并启用 race 检查。埋错基线会接受 stale 或重复 fencing intent。
