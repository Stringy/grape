package reddit

type Redditor struct {
	Name     string
	LKarma   int  `json:"link_karma"`
	CKarma   int  `json:"comment_karma"`
	IsFriend bool `json:"is_friend"`
	HasMail  bool `json:"has_mail"`
	IsOver18 bool `json:"over_18"`
	IsGold   bool `json:"is_gold"`
	IsMod    bool `json:"is_mod"`
	Cookie   string
	ModHash  string
}

func (r *Redditor) IsLoggedIn() bool {
	return r.Cookie != "" && r.ModHash != ""
}
