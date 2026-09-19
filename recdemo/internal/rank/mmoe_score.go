package rank

func MMoEScore(p Prediction)float64{mids:=[]float64{5,20,45,120,240,450,900};watch:=0.0;for i,v:=range p.WatchProbs{if i<len(mids){watch+=v*mids[i]}};watch/=900;if watch>1{watch=1};return 0.20*p.CTR+0.35*watch+0.25*p.Completion+0.20*p.Pay}
