package main

//Baobei 百宝箱单元
type Baobei struct {
	CommonModle
	ProjectID uint   `json:"-"`
	Desc      string `json:"desc"`
	Name      string `json:"name"`
	Tags      string `json:"itags"`
	URL       string `json:"url"`
	Operator  string `json:"operator"`
	InFavor   bool   `gorm:"-" json:"in_favor"`
	Users     []User `gorm:"many2many:user_baobeis;" json:"-"`
}
