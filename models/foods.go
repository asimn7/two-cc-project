package models


import "gorm.io/gorm"

type Books struct {
	ID				uint			`gorm:"primaryKey; autoIncrement" json:"id"`
	Buyer			*string			`json:"buyer"`
	Seller			*string			`json:"seller"`
	Consumer		*string			`json:"consumer"`
}

func MigrateBooks(db *gorm.DB)error{
	err := db.AutoMigrate(&Books{})
	
	return err
}