package repository

import (
	"authentication-service/database"
	"authentication-service/model"
)

func DownloadItem(id string, item *model.ItemWithAssets) error {
	db := database.GrabDB()

	if err := db.Get(&item.Item, "SELECT * FROM products WHERE identifier=$1 LIMIT 1", id); err != nil {
		return err
	}

	if err := db.Select(&item.Assets, "SELECT url FROM product_images WHERE product=$1", id); err != nil {
		return err
	}

	return nil
}

func DownloadAllItems(items *[]model.ItemWithAssets) error {
	db := database.GrabDB()

	var ids []string
	if err := db.Select(&ids, "SELECT identifier FROM products"); err != nil {
		return err
	}

	for _, id := range ids {

		//Retrieving raw database item
		var raw model.DatabaseItem
		if err := db.Get(&raw, "SELECT * FROM products WHERE identifier=$1 LIMIT 1", id); err != nil {
			return err
		}

		//Retrieving assets
		var assets []string
		if err := db.Select(&assets, "SELECT url FROM product_images WHERE product=$1", id); err != nil {
			return err
		}

		var user model.User
		if err := db.Get(&user, "SELECT * FROM users WHERE username=$1", raw.Seller); err != nil {
			return err
		}

		item := model.Item{
			Description: raw.Description,
			Title:       raw.Title,
			Identifier:  raw.Identifier,
			Seller:      user,
			Price:       raw.Price,
			Stock:       raw.Stock,
			Category:    raw.Category,
		}
		*items = append(*items, model.ItemWithAssets{item, assets})
	}

	return nil
}
