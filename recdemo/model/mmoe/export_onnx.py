import argparse
import torch
from .model import MMoE
from .dataset import FEATURES
def export(args):
    model=MMoE(len(FEATURES)); ckpt=torch.load(args.checkpoint,map_location="cpu",weights_only=False); model.load_state_dict(ckpt["state_dict"]); model.eval(); x=torch.zeros(1,len(FEATURES),dtype=torch.float32)
    torch.onnx.export(model,x,args.output,input_names=["features"],output_names=["p_ctr","watch_logits","p_completion","p_pay"],opset_version=17,dynamic_axes={"features":{0:"batch"},"p_ctr":{0:"batch"},"watch_logits":{0:"batch"},"p_completion":{0:"batch"},"p_pay":{0:"batch"}}); print(f"exported {args.output}; feature_count={len(FEATURES)}")
if __name__=="__main__":
    p=argparse.ArgumentParser(); p.add_argument("--checkpoint",required=True); p.add_argument("--output",default="deploy/triton/model_repository/mmoe/1/model.onnx"); export(p.parse_args())
