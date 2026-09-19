import csv
from pathlib import Path
import torch
from torch.utils.data import Dataset
FEATURES=["user_watch_seconds","user_completion","user_pay_rate","user_sessions","drama_popularity","drama_completion","drama_pay_rate","drama_age_days","same_region","same_language","genre_affinity","is_new","is_continue","hour","weekday","country_score","language_score","price_score"]
class BehaviorDataset(Dataset):
    def __init__(self,path):
        rows=[]
        with Path(path).open("r",encoding="utf-8",newline="") as f:
            for row in csv.DictReader(f):
                x=[float(row.get(k,0) or 0) for k in FEATURES]; rows.append((x,int(row.get("ctr",0) or 0),int(row.get("watch_bucket",0) or 0),int(row.get("completion",0) or 0),int(row.get("pay",0) or 0)))
        self.rows=rows
    def __len__(self): return len(self.rows)
    def __getitem__(self,i):
        x,ctr,watch,completion,pay=self.rows[i]; return torch.tensor(x,dtype=torch.float32),torch.tensor(ctr,dtype=torch.float32),torch.tensor(watch,dtype=torch.long),torch.tensor(completion,dtype=torch.float32),torch.tensor(pay,dtype=torch.float32)
