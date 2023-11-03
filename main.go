package main

import (
	"GoAPI/routes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"time"
)

type SessionFuncMap struct {
	Session sessions.Session
}

func (s SessionFuncMap) Get(key interface{}) interface{} {
	return s.Session.Get(key)
}

func add(num, page int) int {
	return num + 6*(page-1) + 1
}

func shortenData(input string, maxLength int) string {
	if len(input) <= maxLength {
		return input
	}
	return input[:maxLength] + "..."
}

func formatTimestamp(timestamp primitive.DateTime) string {
	// Chuyển đổi primitive.DateTime thành time.Time
	t := time.Unix(int64(timestamp)/1000, 0) // Chia cho 1000 để chuyển đổi thành giây

	// Định dạng thời gian sang múi giờ của bạn
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	formattedTime := t.In(loc).Format("2006-01-02 15:04:05")

	return formattedTime
}
func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("your-secret-key"))
	r.Use(sessions.Sessions("my-session", store))

	r.SetFuncMap(template.FuncMap{
		"session": func(ctx *gin.Context) SessionFuncMap {
			return SessionFuncMap{Session: sessions.Default(ctx)}
		},
		"add":             add,
		"formatTimestamp": formatTimestamp,
		"shortenData":     shortenData,
	})
	routes.CreateRouter(r)
	r.Static("/js", "./js")
	r.Static("/uploads", "./uploads")
	r.LoadHTMLGlob("templates/*")
	r.Run(":3000")
}
