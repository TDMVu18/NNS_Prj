package routes

import (
	"GoAPI/controller"
	"GoAPI/middleware"
	"github.com/gin-gonic/gin"
)

func CreateRouter(app *gin.Engine) {

	auth := app.Group("/auth")
	{
		auth.POST("/api/register", controller.Register)
		auth.POST("/api/login", controller.Login)
		auth.GET("/api/logout", controller.Logout)
		auth.GET("/web", controller.Auth)
		auth.GET("/web/signup", controller.SignUp)
	}
	person := app.Group("/person")
	{
		person.Use(middleware.AuthMiddleware())
		info := person.Group("/info")
		{
			info.POST("/api", controller.AddPerson)
			info.GET("/web/profile", controller.ShowProfile)
			info.POST("/api/upload", controller.Upload)
			info.GET("/web", controller.ListPerson)
			info.POST("/api/appearance", controller.ToggleAppearance)
			info.POST("/api/update", controller.UpdatePersonById)
			info.POST("/api/delete", controller.DeletePersonById)

		}
		salary := person.Group("/salary")
		{
			salary.POST("", controller.SalaryAdd)
			salary.GET("", controller.ListSalary)
			salary.POST("/delete", controller.DeleteSalaryById)
			salary.POST("/update", controller.UpdateSalaryById)
			salary.GET("/level", controller.GetSalaryLevels)
		}
		office := person.Group("/office")
		{
			office.POST("", controller.OfficeAdd)
			office.GET("", controller.ListOffice)
			office.POST("/delete", controller.DeleteOfficeById)
			office.POST("/update", controller.UpdateOfficeById)
			office.GET("/name", controller.GetOfficeName)
		}
	}
}
