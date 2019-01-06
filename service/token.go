package service

// Token ...
type Token struct {
	OID           string `json:"oid"`
	Name          string `json:"name"`
	Pwd           string `json:"pwd"`
	EffectiveTime int64  `json:"effectiveTime"`
}
