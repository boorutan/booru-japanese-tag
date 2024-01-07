package translate

import (
	"encoding/csv"
	"github.com/boorutan/booru-japanese-tag/db"
	"math/rand"
	"os"
)

type Tag struct {
	Name           string
	PostCount      int
	Alias          string
	Id             int
	TranslatedName string
}

func GetTag() Tag {
	rows, _ := db.DB.Query("SELECT id, name, post_count, alias FROM tag WHERE translated = false AND category = 0 OR category = 4 ORDER BY post_count DESC LIMIT 5")
	var tags []Tag
	for rows.Next() {
		tag := Tag{}
		_ = rows.Scan(&tag.Id, &tag.Name, &tag.PostCount, &tag.Alias)
		tags = append(tags, tag)
	}
	return tags[rand.Intn(len(tags))]
}

func UpdateTag(en, ja string) error {
	_, err := db.DB.Exec("UPDATE tag SET translated = true, translated_name = ? WHERE name = ?", ja, en)
	return err
}

func ImportDanbooruTag() error {
	_, _ = db.Execute("DELETE FROM tag WHERE translated = false")
	file, err := os.Open("danbooru.csv")
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}
	for _, v := range rows {
		_, err = db.DB.Exec(
			"INSERT INTO tag (name, category, post_count, alias, translated) VALUES (?, ?, ?, ?, false)",
			v[0],
			v[1],
			v[2],
			v[3],
		)
		if err != nil {
			continue
		}
	}
	return nil
}
