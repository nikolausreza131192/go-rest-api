package controllers

import (
	"github.com/nikolausreza131192/pos/entity"
	"github.com/nikolausreza131192/pos/repository"
)

// Item control all process related with item
type Item interface {
	GetAll() []entity.Item
	GetByID(id int) entity.Item
}

type itemController struct {
	itemRP repository.ItemRepo
}

// ItemControllerParam will be used as repository parameter
type ItemControllerParam struct {
	ItemRP repository.ItemRepo
}

// NewItem initialize item controller
func NewItem(param ItemControllerParam) Item {
	return &itemController{
		itemRP: param.ItemRP,
	}
}

// GetAll return all items
func (c *itemController) GetAll() []entity.Item {
	return c.itemRP.GetAll()
}

// GetByID return item by specified ID
func (c *itemController) GetByID(id int) entity.Item {
	return c.itemRP.GetByID(id)
}
