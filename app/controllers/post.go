package controllers

import (
	"goblog/app/models"
	"goblog/app/routes"

	"github.com/revel/revel"
)

// GormController를 임베디드 필드로 지정
type Post struct {
	GormController
}

// Order 메서드와 Find 메서드로 전체 포스트 조회
func (c Post) Index() revel.Result {
	var posts []models.Post
	c.Txn.Order("created_at desc").Find(&posts)
	return c.Render(posts)
}

// First 메서드로 id에 해당하는 포스트 조회
func (c Post) Show(id int) revel.Result {
	var post models.Post
	c.Txn.First(&post, id)
	c.Txn.Where(&models.Comment{PostId: id}).Find(&post.Comments)
	return c.Render(post)
}

// Save 메서드로 포스트 수정
func (c Post) Update(id int, title, body string) revel.Result {
	var post models.Post
	c.Txn.First(&post, id)
	post.Title = title
	post.Body = body

	c.Txn.Save(&post)

	// view 에 flash 메시지 전달
	c.Flash.Success("포스트 수정 완료")

	// 포스트 상세 보기 화면으로 이동
	return c.Redirect(routes.Post.Show(id))
}

// Create 메서드로 포스트 생성
func (c Post) Create(title, body string) revel.Result {
	post := models.Post{Title: title, Body: body}
	c.Txn.Create(&post)

	// 뷰에 Flash 메시지 전달
	c.Flash.Success("포스트 작성 완료")

	// 포스트 목록 화면으로 이동
	return c.Redirect(routes.Post.Index())
}

// Delete 메서드로 포스트 삭제
func (c Post) Destroy(id int) revel.Result {
	c.Txn.Where("post_id = ?", id).Delete(&models.Comment{})
	c.Txn.Where("id = ?", id).Delete(&models.Post{})

	//View에 Flash 메시지 전달
	c.Flash.Success("포스트 삭제 완료")

	// 포스트 목록 화면으로 이동
	return c.Redirect(routes.Post.Index())
}

func (c Post) New() revel.Result {
	post := models.Post{}
	return c.Render(post)
}

//
func (c Post) Edit(id int) revel.Result {
	var post models.Post
	c.Txn.First(&post, id)

	return c.Render(post)
}
