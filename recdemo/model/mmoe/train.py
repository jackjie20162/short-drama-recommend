import argparse
import json
from pathlib import Path
import torch
from torch import nn
from torch.utils.data import DataLoader
from .dataset import BehaviorDataset
from .model import MMoE
from .schema import SPARSE_FEATURES, DENSE_FEATURES, DEFAULT_CARDINALITIES, MODEL_VERSION

def _cardinalities(ds):
    values = []
    for idx, name in enumerate(SPARSE_FEATURES):
        maximum = max((r[0][idx] for r in ds.rows), default=0)
        values.append(max(DEFAULT_CARDINALITIES[name], maximum + 1))
    return values

def _loader(ds, batch):
    return DataLoader(ds, batch_size=batch, shuffle=False)

def _run_epoch(model, loader, opt, pay_weight):
    train = opt is not None
    model.train(train)
    bce = nn.BCEWithLogitsLoss()
    ce = nn.CrossEntropyLoss()
    total = 0.0
    count = 0
    for sparse, dense, ctr, watch, completion, pay in loader:
        args = list(sparse) + [dense]
        out = model(*args)
        loss = (
            bce(out[0], ctr)
            + ce(out[1], watch)
            + 0.8 * bce(out[2], completion)
            + pay_weight * bce(out[3], pay)
        )
        if train:
            opt.zero_grad()
            loss.backward()
            opt.step()
        total += float(loss.detach())
        count += 1
    return total / max(count, 1)

def train(args):
    ds = BehaviorDataset(args.data)
    if len(ds) < 10:
        raise ValueError("training dataset is too small; provide real behavior samples")
    train_ds, val_ds, test_ds = ds.split_time()
    cards = _cardinalities(ds)
    model = MMoE(cards)
    opt = torch.optim.AdamW(model.parameters(), lr=args.lr, weight_decay=1e-4)
    history = []
    for epoch in range(args.epochs):
        loss = _run_epoch(model, _loader(train_ds, args.batch_size), opt, args.pay_weight)
        with torch.no_grad():
            val_loss = _run_epoch(model, _loader(val_ds, args.batch_size), None, args.pay_weight)
        history.append({"epoch": epoch + 1, "train_loss": loss, "val_loss": val_loss})
        print(f"epoch={epoch+1} train_loss={loss:.6f} val_loss={val_loss:.6f}")
    Path(args.output).parent.mkdir(parents=True, exist_ok=True)
    torch.save({
        "state_dict": model.state_dict(),
        "model_version": MODEL_VERSION,
        "sparse_features": SPARSE_FEATURES,
        "dense_features": DENSE_FEATURES,
        "cardinalities": cards,
        "history": history,
        "test_samples": len(test_ds),
    }, args.output)
    meta = Path(args.output).with_suffix(".json")
    meta.write_text(json.dumps({
        "model_version": MODEL_VERSION,
        "sparse_features": SPARSE_FEATURES,
        "dense_features": DENSE_FEATURES,
        "cardinalities": cards,
        "history": history,
        "test_samples": len(test_ds),
    }, ensure_ascii=False, indent=2), encoding="utf-8")

if __name__ == "__main__":
    p = argparse.ArgumentParser()
    p.add_argument("--data", required=True)
    p.add_argument("--output", default="artifacts/mmoe.pt")
    p.add_argument("--epochs", type=int, default=5)
    p.add_argument("--batch-size", type=int, default=512)
    p.add_argument("--lr", type=float, default=1e-3)
    p.add_argument("--pay-weight", type=float, default=1.2)
    train(p.parse_args())
