package database

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	Id       int
	Author   string
	Content  string
	Created  time.Time
	Modified time.Time
}

func Create(p Post) error {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY AUTOINCREMENT, author TEXT, content TEXT, created INT, modified INT)
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert into posts(author, content, created, modified) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	vals := []interface{}{}
	var t = time.Now().UnixMilli()
	vals = append(vals, p.Author, p.Content, t, t)
	_, err = stmt.Exec(vals...)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func Read() ([]Post, error) {
	var posts []Post
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("select id, author, content, created, modified from posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p Post
		var tc int64
		var tm int64
		err = rows.Scan(&p.Id, &p.Author, &p.Content, &tc, &tm)
		if err != nil {
			return nil, err
		}
		p.Created = time.UnixMilli(tc)
		p.Modified = time.UnixMilli(tm)
		posts = append(posts, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func ReadOne(id int) (*Post, error) {
	var post Post
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT author, content, created FROM posts WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var tc int64
	err = stmt.QueryRow(id).Scan(&post.Author, &post.Content, &tc)
	post.Created = time.UnixMilli(tc)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func Update(p Post) error {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("UPDATE posts SET author = ?, content = ?, modified = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	vals := []interface{}{}
	vals = append(vals, p.Author, p.Content, time.Now().UnixMilli(), p.Id)
	_, err = stmt.Exec(vals...)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func Delete(id int) error {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("DELETE FROM posts WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
