package controllers

import (
	"github.com/revel/revel"
)

// 웹 어플리케이션 초기화 로직 등록
func init() {
	revel.OnAppStart(InitDB)
	revel.InterceptMethod((*GormController).Begin, revel.BEFORE)
	revel.InterceptMethod((*GormController).Commit, revel.AFTER)
	revel.InterceptMethod((*GormController).Rollback, revel.FINALLY)

	// 모든 액션마다 setCurrentUser가 수행되도록 인터셉터로 등록
	revel.InterceptMethod((*App).setCurrentUser, revel.BEFORE)
}
