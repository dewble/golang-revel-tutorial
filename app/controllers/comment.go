package controllers

import (
	"goblog/app/models"
	"goblog/app/routes"

	"github.com/revel/revel"
)

type Comment struct {
	GormController
}

func (c Comment) Create(postId int, body, commenter string) revel.Result {
	// 댓글 저장
	comment := models.Comment{PostId: postId, Body: body, Commenter: commenter}
	c.Txn.Create(&comment)

	// 뷰에 Flash 메시지 전달
	c.Flash.Success("댓글 작성 완료")

	// 포스트 상세 보기 화면으로 이동
	return c.Redirect(routes.Post.Show(postId))
}

func (c Comment) Destroy(postId, id int) revel.Result {
	// 댓글 삭제
	c.Txn.Where("id=?", id).Delete(&models.Comment{})

	// 뷰에 Flash 메시지 전달
	c.Flash.Success("댓글 삭제 완료")

	// 포스트 상세 보기 화면으로 이동
	return c.Redirect(routes.Post.Show(postId))
}
