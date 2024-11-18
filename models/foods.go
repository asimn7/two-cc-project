package models


import "gorm.io/gorm"

type Foods struct {
	ID				uint			`gorm:"primaryKey; autoIncrement" json:"id"`
	Buyer			*string			`json:"buyer"`
	Seller			*string			`json:"seller"`
	Consumer		*string			`json:"consumer"`
}

func MigrateFoods(db *gorm.DB)error{
	err := db.AutoMigrate(&Foods{})
	
	return err
}