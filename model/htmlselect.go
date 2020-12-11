package model

import (
	"gorm.io/gorm"
)

//HTMLSelect used to fill the html select
type HTMLSelect struct {
	gorm.Model
	Name     string `gorm:"type:varchar(40)"`
	ParentID uint
}

//Node used to build trees
type Node struct {
	ID       uint    `json:"id"`
	ParentID uint    `json:"pid"`
	Name     string  `json:"name"`
	Child    []*Node `json:"child"`
}

//Response used to json response
type Response struct {
	Data []*Node `json:"data"`
}

//GetTrees get the htmlselect tress
func GetTrees(pid uint) []*Node {
	var htmls []HTMLSelect
	DB.Where("parent_id = ?", pid).Find(&htmls)
	var tree []*Node
	for _, v := range htmls {
		child := GetTrees(v.ID)
		node := &Node{
			ID:       v.ID,
			Name:     v.Name,
			ParentID: v.ParentID,
		}
		node.Child = child
		tree = append(tree, node)
	}
	return tree
}

//Insert Insert values into htmlselect database
func Insert(name string, pid uint) {
	var hs HTMLSelect
	hs.Name = name
	hs.ParentID = pid
	DB.Create(&hs)
}
