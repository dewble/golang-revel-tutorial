package controllers

import (
	"goblog/app/models"
	"goblog/app/routes"
	"time"

	db "github.com/revel/modules/db/app"
	"github.com/revel/revel"
)

type Post struct {
	*revel.Controller
	db.Transactional
}

func (c Post) Index() revel.Result {
	var posts []models.Post
	rows, err := c.Txn.Query("select id, title, body, created_at, updated_at from posts order by created_at desc")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		post := models.Post{}
		if err := rows.Scan(&post.Id, &post.Title, &post.Body, &post.CreatedAt, &post.UpdatedAt); err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}

	return c.Render(posts)
}

func (c Post) New() revel.Result {
	post := models.Post{}
	return c.Render(post)
}

func (c Post) Create(title, body string) revel.Result {
	// 데이터베이스에 포스트 내용 저장
	_, err := c.Txn.Exec("insert into posts(title, body, created_at, updated_at)values(?,?,?,?)", title, body, time.Now(), time.Now())

	if err != nil {
		panic(err)
	}

	// 뷰에 Flash 메시지 전달
	c.Flash.Success("포스트 작성 완료")

	// 포스트 목록 화면으로 이동
	return c.Redirect(routes.Post.Index())
}
