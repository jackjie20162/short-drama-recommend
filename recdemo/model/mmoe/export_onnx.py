import argparse
import torch
from .model import MMoE

def export(args):
    ckpt = torch.load(args.checkpoint, map_location="cpu", weights_only=False)
    model = MMoE(ckpt["cardinalities"])
    model.load_state_dict(ckpt["state_dict"])
    model.eval()
    sparse = [torch.zeros(1, dtype=torch.long) for _ in ckpt["sparse_features"]]
    dense = torch.zeros(1, len(ckpt["dense_features"]), dtype=torch.float32)
    inputs = tuple(sparse + [dense])
    names = list(ckpt["sparse_features"]) + ["dense_features"]
    outputs = ["p_ctr", "watch_logits", "p_completion", "p_pay"]
    torch.onnx.export(
        model, inputs, args.output,
        input_names=names, output_names=outputs,
        opset_version=17,
        dynamic_axes={name: {0: "batch"} for name in names + outputs},
    )
    print(f"exported {args.output}; version={ckpt['model_version']}")

if __name__ == "__main__":
    p = argparse.ArgumentParser()
    p.add_argument("--checkpoint", required=True)
    p.add_argument("--output", default="deploy/triton/model_repository/mmoe/1/model.onnx")
    export(p.parse_args())
