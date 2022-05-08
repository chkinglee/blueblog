package article_repo

import "time"

// Article
//go:generate gormgen -structs Article -input .
type Article struct {
	Id         string    //
	Uid        string    //
	Title      string    //
	Content    string    //
	CreateTime time.Time `gorm:"time"` //
}
