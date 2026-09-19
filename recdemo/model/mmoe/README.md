# MMoE training

This module trains CTR, watch-time bucket, completion and pay jointly with shared experts and task-specific gates/towers.

CSV columns are the feature names in dataset.py followed by ctr, watch_bucket, completion and pay. Feature order is a model contract and must not change without a model version bump.

Train:
python -m recdemo.model.mmoe.train --data data/mmoe_train.csv --output artifacts/mmoe.pt

Export:
python -m recdemo.model.mmoe.export_onnx --checkpoint artifacts/mmoe.pt --output deploy/triton/model_repository/mmoe/1/model.onnx

The ONNX binary is generated separately and is intentionally not committed to source control.
