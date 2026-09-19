# MMoE 落地说明

链路：ES 召回/过滤 → Feature Builder → MMoE（CTR/WatchTime/Completion/Pay）→ 组合分数 → Rerank → TopN。

代码位置：model/mmoe 训练与导出；internal/feature/builder.go 特征；internal/rank/mmoe_client.go Triton；internal/rank/mmoe_score.go 组合分数；deploy/triton/model_repository/mmoe/config.pbtxt Triton 配置。

ONNX 二进制不提交 Git。先训练生成 model.onnx，再放到 Triton 的 mmoe/1/。模型不可用时继续使用现有 formula.go fallback。

下一阶段把用户特征从 seed/内存替换成 user-rpc + behavior-rpc/feature-rpc，并增加模型版本、校准和 A/B 实验。
