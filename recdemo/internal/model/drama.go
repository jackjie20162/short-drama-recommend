package model

type Drama struct {
	ID int64 `json:"id"`
	Title string `json:"title"`
	Cover string `json:"cover"`
	Region string `json:"region"`
	Lang string `json:"lang"`
	Status string `json:"status"`
	Popularity float64 `json:"popularity"`
	Completion float64 `json:"completion"`
	PayRate float64 `json:"pay_rate"`
	PublishedAt int64 `json:"published_at"`
	SeenBy map[string]bool `json:"seen_by,omitempty"`
}

type Query struct {
	UserID string
	Region string
	Lang string
	SeenIDs []int64
}

func Seeds() []Drama {
	return []Drama{
		{1,"Contract Wife","/covers/1.jpg","0","en","published",0.91,0.82,0.35,1720000000,map[string]bool{"u1":true}},
		{2,"Hidden Heiress","/covers/2.jpg","0","en","published",0.88,0.90,0.42,1721000000,map[string]bool{"u1":false}},
		{3,"CEO Next Door","/covers/3.jpg","0","en","published",0.78,0.71,0.22,1722000000,map[string]bool{}},
		{4,"Revenge Contract","/covers/4.jpg","0","zh","published",0.84,0.77,0.31,1723000000,map[string]bool{}},
		{5,"Love After Divorce","/covers/5.jpg","0","en","published",0.76,0.69,0.18,1724000000,map[string]bool{}},
		{6,"My Secret Billionaire","/covers/6.jpg","0","en","published",0.86,0.88,0.39,1725000000,map[string]bool{}},
		{7,"The Silent CEO","/covers/7.jpg","0","en","published",0.73,0.66,0.15,1726000000,map[string]bool{}},
		{8,"Second Chance","/covers/8.jpg","0","es","published",0.69,0.72,0.19,1727000000,map[string]bool{}},
		{9,"Billionaire's Secret","/covers/9.jpg","0","en","published",0.82,0.79,0.36,1728000000,map[string]bool{}},
		{10,"Offline Test","/covers/10.jpg","0","en","offline",0.99,0.99,0.99,1729000000,map[string]bool{}},
		{11,"Love in Tokyo","/covers/11.jpg","1","ja","published",0.92,0.86,0.28,1730000000,map[string]bool{}},
		{12,"Seoul Contract","/covers/12.jpg","2","ko","published",0.90,0.84,0.33,1731000000,map[string]bool{}},
		{13,"Desert Promise","/covers/13.jpg","3","ar","published",0.87,0.80,0.29,1732000000,map[string]bool{}},
		{14,"Paris After Midnight","/covers/14.jpg","4","fr","published",0.81,0.75,0.25,1733000000,map[string]bool{}},
		{15,"Spanish Hearts","/covers/15.jpg","5","es","published",0.79,0.78,0.27,1734000000,map[string]bool{}},
		{16,"Global Heiress","/covers/16.jpg","0","en","published",0.94,0.91,0.45,1735000000,map[string]bool{}},
		{17,"Rainy Marriage","/covers/17.jpg","0","en","published",0.74,0.70,0.20,1736000000,map[string]bool{}},
		{18,"CEO's Secret","/covers/18.jpg","0","en","published",0.89,0.83,0.41,1737000000,map[string]bool{}},
		{19,"After the Vow","/covers/19.jpg","0","en","published",0.71,0.68,0.17,1738000000,map[string]bool{}},
		{20,"Last Summer","/covers/20.jpg","0","en","published",0.77,0.73,0.21,1739000000,map[string]bool{}},
	}
}
