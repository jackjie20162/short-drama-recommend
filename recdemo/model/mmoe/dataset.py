import csv
from pathlib import Path
import torch
from torch.utils.data import Dataset
from .schema import SPARSE_FEATURES, DENSE_FEATURES, WATCH_BUCKETS

SPARSE = SPARSE_FEATURES
DENSE = DENSE_FEATURES

class BehaviorDataset(Dataset):
    def __init__(self, path):
        rows = []
        with Path(path).open("r", encoding="utf-8", newline="") as f:
            for row in csv.DictReader(f):
                sparse = [int(float(row.get(k, 0) or 0)) for k in SPARSE]
                dense = [float(row.get(k, 0) or 0) for k in DENSE]
                watch = int(row.get("watch_bucket", 0) or 0)
                watch = max(0, min(WATCH_BUCKETS - 1, watch))
                rows.append((
                    sparse,
                    dense,
                    float(row.get("ctr", 0) or 0),
                    watch,
                    float(row.get("completion", 0) or 0),
                    float(row.get("pay", 0) or 0),
                    int(row.get("timestamp", 0) or 0),
                ))
        rows.sort(key=lambda r: r[-1])
        self.rows = rows

    def __len__(self):
        return len(self.rows)

    def __getitem__(self, i):
        sparse, dense, ctr, watch, completion, pay, _ = self.rows[i]
        return (
            tuple(torch.tensor(v, dtype=torch.long) for v in sparse),
            torch.tensor(dense, dtype=torch.float32),
            torch.tensor(ctr, dtype=torch.float32),
            torch.tensor(watch, dtype=torch.long),
            torch.tensor(completion, dtype=torch.float32),
            torch.tensor(pay, dtype=torch.float32),
        )

    def split_time(self, train_ratio=0.8, val_ratio=0.1):
        n = len(self.rows)
        a = int(n * train_ratio)
        b = int(n * (train_ratio + val_ratio))
        return self._subset(0, a), self._subset(a, b), self._subset(b, n)

    def _subset(self, start, end):
        ds = BehaviorDataset.__new__(BehaviorDataset)
        ds.rows = self.rows[start:end]
        return ds
