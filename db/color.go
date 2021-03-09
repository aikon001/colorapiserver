package db

import (
	"encoding/hex"
	"errors"
	"unsafe"

	"github.com/aikon001/colorapiserver/models"
)

func (db Database) GetAllColors() (*models.ColorList, error) {
	list := &models.ColorList{}
	rows, err := db.Conn.Query("SELECT * FROM colors ORDER BY id DESC")
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var color models.Color
		err := rows.Scan(&color.ID, &color.Name, &color.hexadecimal, &color.R, &color.G, &color.B, &color.CreatedAt)
		if err != nil {
			return list, err
		}
		list.Colors = append(list.Colors, color)

	}
	return list, nil
}

func (db Database) AddColor(color *models.Color) error {
	var id int
	var createdAt string

	query := `INSERT INTO colors (name, hexadecimal,R,G,B) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at`

	if len(color.hexadecimal) != 0 {
		byt, _ := hex.DecodeString(color.hexadecimal)
		err := db.Conn.QueryRow(query, color.Name, color.hexadecimal, byt[0], byt[1], byt[2]).Scan(&id, &createdAt)
		if err != nil {
			return err
		}

	} else if unsafe.Sizeof(color.R)+unsafe.Sizeof(color.G)+unsafe.Sizeof(color.B) != 0 {
		rgb := (*[3]byte)(unsafe.Pointer(&color.R))[:]
		err := db.Conn.QueryRow(query, color.Name, color.hexadecimal, rgb[2], rgb[1], rgb[0]).Scan(&id, &createdAt)
		if err != nil {
			return err
		}

	} else {
		return errors.New("No hexadecimal provided [Neither RGB provided!]")
	}
	color.ID = id
	color.CreatedAt = createdAt
	return nil

}

func (db Database) GetItemById(colorId int) (models.Color, error) {
	color := models.Color{}
	query := `SELECT * FROM colors WHERE id = $1;`
	row := db.Conn.QueryRow(query, colorId)
	err := row.Scan(&color.ID, &color.Name, &color.hexadecimal, &color.R, &color.G, &color.B, &color.CreatedAt)
	return color, err
}

func (db Database) DeleteItem(colorId int) error {
	query := `DELETE FROM items WHERE id = $1;`
	_, err := db.Conn.Exec(query, colorId)
	return err
}