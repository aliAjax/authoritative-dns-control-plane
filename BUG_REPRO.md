# Bug Reproduction

rollback 复用旧快照中的嵌套 records，并把旧状态直接带回发布流程，随后 publish 读取到被污染的状态或共享值。

复现：运行 collection.json 中的两条 snapshot_release 定向测试。埋错基线无法保持 rollback 后的记录值和发布状态一致。
