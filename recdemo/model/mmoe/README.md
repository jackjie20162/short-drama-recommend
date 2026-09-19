# MMoE V1 training / serving contract

The V1 model uses five sparse IDs plus eight dense features.

Sparse:
- user_id
- drama_id
- region_id
- language_id
- genre_id

Dense:
- user_watch_seconds
- user_completion
- user_pay_rate
- user_sessions
- drama_popularity
- drama_completion
- drama_pay_rate
- drama_age_days

The model has four objectives: CTR, watch-time bucket, completion and payment.

## Train

Use Python 3.10-3.12 with a PyTorch version that supports your platform.

`data/mmoe_train.csv` must contain the feature columns above plus:
`timestamp,ctr,watch_bucket,completion,pay`.

The loader sorts samples by timestamp and uses an 80/10/10 time split. No random split is used.

`python -m recdemo.model.mmoe.train --data data/mmoe_train.csv --output artifacts/mmoe.pt`

Training produces a checkpoint and a JSON metadata file next to it.

## Export

`python -m recdemo.model.mmoe.export_onnx --checkpoint artifacts/mmoe.pt --output deploy/triton/model_repository/mmoe/1/model.onnx`

The ONNX input names and order are part of the `mmoe.v1` contract. Do not reorder them without bumping the model version.

## Important

A real model file is not committed. `MMoE.Enabled` stays false until a trained ONNX model has been generated and deployed.
