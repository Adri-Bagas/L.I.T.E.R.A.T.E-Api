package models

import (
	A "perpus_api/db"
	"sync"
)

type Media struct {
	Id        int    `json:"id" db:"id"`
	ModelName string `json:"model_name" db:"model_name"`
	ModelId   int    `json:"model_id" db:"model_id"`
	MediaType string `json:"media_type" db:"media_type"`
	Location  string `json:"location" db:"location"`
}

var MediaLock = sync.Mutex{}

func FindMedia(mName string, mId int) (*Media, error) {
	MediaLock.Lock()
	defer MediaLock.Unlock()

	var obj Media

	con := A.GetDB()

	sql := `
		SELECT * FROM public.medias where model_name = $1 and model_id = $2 ORDER BY id ASC ;
	`

	rows, err := con.Query(sql, mName, mId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&obj.Id,
			&obj.ModelName,
			&obj.ModelId,
			&obj.MediaType,
			&obj.Location,
		)
		if err != nil {
			return nil, err
		}
	}

	return &obj, nil
}

func DeleteMedia(id int64) (error) {
	MediaLock.Lock()
	defer MediaLock.Unlock()

	con := A.GetDB()

	sql := `
		DELETE FROM public.medias WHERE id = $1;
	`

	_, err := con.Exec(sql, id)

	if err != nil {
		return err
	}

	return nil
}