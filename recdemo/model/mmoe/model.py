import torch
from torch import nn
from .schema import SPARSE_FEATURES, DENSE_FEATURES, WATCH_BUCKETS

class Expert(nn.Module):
    def __init__(self, input_dim, hidden_dim):
        super().__init__()
        self.net = nn.Sequential(
            nn.Linear(input_dim, hidden_dim),
            nn.ReLU(),
            nn.Linear(hidden_dim, input_dim),
            nn.ReLU(),
        )

    def forward(self, x):
        return self.net(x)

class Tower(nn.Module):
    def __init__(self, input_dim, hidden_dim, out_dim=1):
        super().__init__()
        self.net = nn.Sequential(
            nn.Linear(input_dim, hidden_dim),
            nn.ReLU(),
            nn.Linear(hidden_dim, out_dim),
        )

    def forward(self, x):
        return self.net(x)

class MMoE(nn.Module):
    def __init__(
        self,
        cardinalities=None,
        embedding_dim=16,
        dense_dim=len(DENSE_FEATURES),
        expert_count=4,
        expert_hidden=128,
        tower_hidden=64,
    ):
        super().__init__()
        if cardinalities is None:
            cardinalities = [1_000_000, 1_000_000, 256, 256, 512]
        self.embeddings = nn.ModuleList([
            nn.Embedding(int(size), embedding_dim) for size in cardinalities
        ])
        input_dim = embedding_dim * len(SPARSE_FEATURES) + dense_dim
        self.experts = nn.ModuleList([
            Expert(input_dim, expert_hidden) for _ in range(expert_count)
        ])
        self.gates = nn.ModuleDict({
            k: nn.Linear(input_dim, expert_count)
            for k in ("ctr", "watch", "completion", "pay")
        })
        self.towers = nn.ModuleDict({
            "ctr": Tower(input_dim, tower_hidden),
            "watch": Tower(input_dim, tower_hidden, WATCH_BUCKETS),
            "completion": Tower(input_dim, tower_hidden),
            "pay": Tower(input_dim, tower_hidden),
        })

    def _mix(self, x, task):
        experts = torch.stack([e(x) for e in self.experts], dim=1)
        weights = torch.softmax(self.gates[task](x), dim=-1).unsqueeze(-1)
        return (experts * weights).sum(dim=1)

    def forward(self, *inputs):
        sparse = inputs[:len(SPARSE_FEATURES)]
        dense = inputs[len(SPARSE_FEATURES)]
        embedded = [emb(s.long()) for emb, s in zip(self.embeddings, sparse)]
        x = torch.cat(embedded + [dense.float()], dim=-1)
        return (
            self.towers["ctr"](self._mix(x, "ctr")).squeeze(-1),
            self.towers["watch"](self._mix(x, "watch")),
            self.towers["completion"](self._mix(x, "completion")).squeeze(-1),
            self.towers["pay"](self._mix(x, "pay")).squeeze(-1),
        )
