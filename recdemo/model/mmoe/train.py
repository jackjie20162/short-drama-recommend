import argparse
from pathlib import Path
import torch
from torch import nn
from torch.utils.data import DataLoader
from .dataset import BehaviorDataset,FEATURES
from .model import MMoE
def train(args):
    ds=BehaviorDataset(args.data)
    if not len(ds): raise ValueError("training dataset is empty")
    loader=DataLoader(ds,batch_size=args.batch_size,shuffle=True); model=MMoE(len(FEATURES)); opt=torch.optim.AdamW(model.parameters(),lr=args.lr,weight_decay=1e-4)
    bce=nn.BCEWithLogitsLoss(); ce=nn.CrossEntropyLoss()
    for epoch in range(args.epochs):
        model.train(); total=0.0
        for x,ctr,watch,completion,pay in loader:
            out=model(x); loss=bce(out["p_ctr"],ctr)+ce(out["watch_logits"],watch)+bce(out["p_completion"],completion)+args.pay_weight*bce(out["p_pay"],pay)
            opt.zero_grad(); loss.backward(); opt.step(); total+=loss.item()
        print(f"epoch={epoch+1} loss={total/max(1,len(loader)):.6f}")
    Path(args.output).parent.mkdir(parents=True,exist_ok=True); torch.save({"state_dict":model.state_dict(),"feature_names":FEATURES},args.output)
if __name__=="__main__":
    p=argparse.ArgumentParser(); p.add_argument("--data",required=True); p.add_argument("--output",default="artifacts/mmoe.pt"); p.add_argument("--epochs",type=int,default=5); p.add_argument("--batch-size",type=int,default=256); p.add_argument("--lr",type=float,default=1e-3); p.add_argument("--pay-weight",type=float,default=2.0); train(p.parse_args())
