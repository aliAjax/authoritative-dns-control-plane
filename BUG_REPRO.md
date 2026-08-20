# Bug Reproduction

内存仓库直接保存和返回带嵌套切片的领域对象，调用方修改 Values、NS 或快照 Records 后会污染仓库状态，并发访问时可触发数据竞争。

复现：运行 collection.json 中的三条 repository 定向测试。埋错基线中测试均失败，典型信息包括 `stored record changed through caller buffer`、`stored zone changed through caller slice` 和 `snapshot stored caller-owned values`。
