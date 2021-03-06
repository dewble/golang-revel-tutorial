package controllers

import (
	"goblog/app/models"

	"golang.org/x/crypto/bcrypt"

	"github.com/revel/revel"
)

type App struct {
	GormController
	CurrentUser *models.User
}

func (c App) Login() revel.Result {
	return c.Render()
}

// username, password로 인증 확인 후 세션 정보를 생성
func (c App) CreateSession(username, password string) revel.Result {
	var user models.User

	// username으로 사용자 조회
	// First find first record that match given conditions, order by primary key
	c.Txn.Where(&models.User{Username: username}).First(&user)

	// bcrpyt 패키지의 CompareHashAndPassword 함수로 패스워드 비교
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	// 패스워드가 일치하면 세션 생성후 포스트 목록 화면으로 이동
	if err == nil {
		authKey := revel.Sign(user.Username)
		c.Session["authKey"] = authKey
		c.Session["username"] = user.Username
		c.Flash.Success("Welcome, " + user.Name)
		return c.Redirect(Post.Index)
	}

	// 세션 정보를 모두 제거하고 홈으로 이동
	for k := range c.Session {
		delete(c.Session, k)
	}
	c.Flash.Out["username"] = username
	c.Flash.Error("Login failed")
	return c.Redirect(Home.Index)
}

// 로그아웃 처리
func (c App) DestroySession() revel.Result {
	// clear session
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Render(Home.Index)
}

func (c *App) setCurrentUser() revel.Result {
	// 뷰에서 currentUser를 사용할 수 있도록 RenderArgs에 CurrentUser를 추가
	defer func() {
		if c.CurrentUser != nil {
			c.ViewArgs["currentUser"] = c.CurrentUser // use ViewArgs instead of RenderArgs
		} else {
			delete(c.ViewArgs, "currentUser")
		}
	}()

	// 세션에서 username과 authKey를 가져옴
	username, ok := c.Session["username"]
	if !ok || username == "" {
		return nil
	}

	authKey, ok := c.Session["authKey"]
	if !ok || authKey == "" {
		return nil
	}

	// revel의 Verify 함수를 사용해 authKey가 유효한지 확인
	// authKey가 유효하면 username으로 사용자를 조회하고 컨트롤러의 CurrentUser에 저장
	if match := revel.Verify(username, authKey); match {
		var user models.User
		c.Txn.Where(&models.User{Username: username}).First(&user)
		if &user != nil {
			c.CurrentUser = &user
		}
	}
	return nil
}

/*
type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	greeting := "Good"
	return c.Render(greeting)
}

func (c App) Hello(email string) revel.Result {
	// Parameter 로 email을 받아서 view에 email을 전달한다.
	//c.Validation.Required(email).Message("email is need")
	//c.Validation.MinSize(email, 3).Message("type over 3 word")

	c.Validation.Email(email).Message("Email 형식을 확인하세요 user@mail.com ")

	if c.Validation.HasErrors() {
		c.Validation.Keep() // Flash cookie에 VaildationErros 저장
		//c.FlashParams()
		return c.Redirect(App.Index)
	}
	return c.Render(email)
}
*/
