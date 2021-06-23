package controllers

import (
	"database/sql"
	"fmt"
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

func getPost(txn *sql.Tx, id int) (models.Post, error) {
	post := models.Post{}
	err := txn.QueryRow("select id, title, body, created_at, updated_at from posts where id=?", id).
		Scan(&post.Id, &post.Title, &post.Body, &post.CreatedAt, &post.UpdatedAt)

	switch {
	case err == sql.ErrNoRows:
		return post, fmt.Errorf("No post with that ID - %d.", id)
	case err != nil:
		return post, err
	}

	// 포스트의 댓글 조회
	post.Comments = getComments(txn, id)
	return post, nil
}

func getComments(txn *sql.Tx, postId int) (comments []models.Comment) {
	rows, err := txn.Query("select id, body, commenter, post_id, created_at, updated_at from comments where post_id=? order by created_at desc", postId)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		comment := models.Comment{}
		if err := rows.Scan(&comment.Id, &comment.Body, &comment.Commenter, &comment.PostId, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
			panic(err)
		}
		comments = append(comments, comment)
	}
	return
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

// 각 포스트 보여주기
// show action method 작성
func (c Post) Show(id int) revel.Result {
	post, err := getPost(c.Txn, id)
	if err != nil {
		panic(err)
	}

	return c.Render(post)
}

// getPost() 함수를 사용하여 매개변수로 전달된 id에 해당하는 포스트를 조회해서 뷰로 렌더링
func (c Post) Edit(id int) revel.Result {
	post, err := getPost(c.Txn, id)
	if err != nil {
		panic(err)
	}

	return c.Render(post)
}

// 사용자가 form 으로 전송한 포스트 정보를 데이터베이스에 저장하는 Update 액션
func (c Post) Update(id int, title, body string) revel.Result {

	// 포스트 내용 수정
	if _, err := c.Txn.Exec("update posts set title=?, body=?, updated_at=? where id=?", title, body, time.Now(), id); err != nil {
		panic(err)
	}

	// view 에 flash 메시지 전달
	c.Flash.Success("포스트 수정 완료")

	// 포스트 상세 보기 화면으로 이동
	return c.Redirect(routes.Post.Show(id))
}

func (c Post) Destroy(id int) revel.Result {
	// 포스트 삭제
	if _, err := c.Txn.Exec("delete from posts where id=?", id); err != nil {
		panic(err)
	}

	//View에 Flash 메시지 전달
	c.Flash.Success("포스트 삭제 완료")

	// 포스트 목록 화면으로 이동
	return c.Redirect(routes.Post.Index())
}
