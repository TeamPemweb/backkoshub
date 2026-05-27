package main

import (
	"net/http"

	"github.com/TeamPemweb/backkoshub/controllers"
	"github.com/TeamPemweb/backkoshub/middleware"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API is running",
		})
	})

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
			auth.POST("/forgot-password", controllers.ForgotPassword)
			auth.POST("/reset-password", controllers.ResetPassword)
			
			authProtected := auth.Group("")
			authProtected.Use(middleware.RequireAuth)
			{
				authProtected.POST("/choose-role", controllers.ChooseRole)
				authProtected.POST("/profile/setup", controllers.SetupProfile)
				authProtected.GET("/profile", controllers.GetProfile)
				authProtected.PUT("/profile", controllers.UpdateProfile)
				authProtected.DELETE("/account", controllers.DeleteResidentAccount)
				
			}
		}
		ownerGroup := v1.Group("/owner")
		ownerGroup.Use(middleware.RequireAuth, middleware.RoleMiddleware("pemilik"))
		{
			ownerGroup.POST("/room-types", controllers.CreateRoomType)
			ownerGroup.GET("/room-types", controllers.GetRoomTypes)
			ownerGroup.PUT("/room-types/:id", controllers.UpdateRoomType)
			ownerGroup.DELETE("/room-types/:id", controllers.DeleteRoomType)

			ownerGroup.POST("/rooms", controllers.CreateRoom)
			ownerGroup.GET("/rooms", controllers.GetRooms)
			ownerGroup.PUT("/rooms/:id", controllers.UpdateRoom)
			ownerGroup.DELETE("/rooms/:id", controllers.DeleteRoom)
			ownerGroup.DELETE("/rooms/:id/end-lease", controllers.EndLease)

			ownerGroup.GET("/dashboard/stats", controllers.GetDashboardStats)
			ownerGroup.GET("/dashboard/unpaid-residents", controllers.GetUnpaidResidents)
			ownerGroup.GET("/residents", controllers.GetActiveResidents)

			ownerGroup.GET("/billings", controllers.GetBillings)
			ownerGroup.PUT("/billings/:id", controllers.ConfirmBillingPaid)
			ownerGroup.PUT("/billings/:id/nominal", controllers.UpdateBillingNominal)

			ownerGroup.GET("/complaints", controllers.GetComplaints)
			ownerGroup.PUT("/complaints/:id", controllers.UpdateComplaintStatus)
		}

		residentGroup := v1.Group("/resident")
		residentGroup.Use(middleware.RequireAuth, middleware.RoleMiddleware("penghuni"))
		{
			residentGroup.POST("/rooms/join", controllers.JoinRoom)
			residentGroup.GET("/my-room", controllers.GetMyRoom)

			residentGroup.GET("/my-billings", controllers.GetMyBillings)
			residentGroup.POST("/my-billings/:id/pay", controllers.PayBilling)
			residentGroup.GET("/my-billings/history", controllers.GetMyBillingHistory)

			residentGroup.POST("/my-complaints", controllers.CreateComplaint)
			residentGroup.GET("/my-complaints/history", controllers.GetMyComplaints)

		}
	}
}