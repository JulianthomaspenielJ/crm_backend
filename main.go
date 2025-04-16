package main

import (
	signup "auth/authsignup"
	"crm/queue"
	"crm/routes"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	go checkTasksEveryMinute()
	q := queue.NewMemoryQueue()
	go q.Consume()

	r := gin.Default()
	r.Use(cors.Default())
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		c.Writer.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.JSON(http.StatusOK, nil)
			return
		}
		c.Next()
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/auth/google/callback", signup.GoogleCallbackHandler)
	r.POST("meetgoogle", routes.MeetGoogle)

	r.GET("/approveuser", routes.ApproveUserList)
	r.GET("/getcomponentlist", routes.FetchComponents)
	r.GET("/getcomponents", routes.GetComponents)
	r.GET("/getproductlist", routes.GetProductList)
	r.GET("/getproduct", routes.GetProduct)
	r.GET("/productmainlist", routes.GetProductMainList)
	r.GET("/productdeleteddata", routes.GetDeletedData)
	r.GET("/userslist", routes.UserList)
	r.GET("/userdeleteddata", routes.GetDeletedUsersData)
	r.GET("/usernamelist", routes.UserNameList)

	r.POST("/restoreuser", routes.RestoreUserManagementData)
	r.POST("removeuser", routes.RemoveUser)
	r.POST("/updateuser", routes.UpdateUser)
	r.POST("/updateusermanagementdetails", routes.UpdateUserManagementData)
	r.POST("/userlogin", routes.LoginHandler)
	r.GET("respectivesalesperson", routes.RespectiveSalesPerson)
	r.GET("getteamlist", routes.GetTeam)

	r.POST("/componentlist", routes.ComponentList)
	r.POST("/getproductdata", routes.GetProductsData)
	r.POST("/updatecomponent", routes.UpdateComponent)
	r.POST("/postcomponentsdata", routes.PostComponentsData)
	r.POST("/deletecomponent", routes.DeleteComponent)
	r.POST("/updateuserdetails", routes.UpdateUserDetails)
	r.POST("/recovercomponent", routes.RecoverComponent)
	r.POST("/updatecomponentlist", routes.UpdateComponentList)
	r.POST("/deletecomponentlist", routes.DeleteComponentList)
	r.POST("getprocessordata", routes.GetProcessorData)

	r.POST("/pricing", routes.Pricing)
	r.POST("/productsdata", routes.Products)
	r.POST("/postproduct", routes.PostProduct)
	r.POST("/productID", routes.ProductID)
	r.POST("/updateproduct", routes.UpdateProductData)
	r.POST("/removeproduct", routes.RemoveProduct)
	r.POST("productsreference", routes.ProductReference)

	r.POST("/auditdata", routes.AuditData)
	r.POST("/audit", routes.AuditTable)
	r.POST("/createproductmainlist", routes.CreateProductMainList)

	// Companies Handler
	r.POST("addcompany", routes.AddCompany)
	r.POST("getcompanydata", routes.GetCompanyData)
	r.POST("deletecompany", routes.DeleteContact)
	r.POST("insertcompanycontact", routes.InsertCompanyContact)

	r.GET("getallcompany", routes.GetAllCompany)

	// r.POST("mail", routes.Mail)
	r.POST("addcontact", routes.AddContact)
	r.POST("getcontactdata", routes.GetContactData)
	r.POST("editcontactdata", routes.EditContactData)
	r.POST("editcompanydata", routes.EditCompanyData)
	r.POST("gst", routes.GST)

	r.POST("insertnote", routes.InsertNotes)
	r.POST("inserttask", routes.InsertTask)
	r.POST("insertmeeting", routes.InsertMeeting)
	r.POST("insertcalls", routes.InsertCalls)

	r.POST("addmargin", routes.AddMargin)
	r.POST("addmarginvalues", routes.AddMarginValues)
	r.POST("updatemargin", routes.UpdateMargin)
	r.POST("deletemargin", routes.DeleteMargin)
	r.POST("deletemarginvalues", routes.DeleteMarginValues)
	r.POST("updatemarginvalues", routes.UpdateMarginValues)
	r.POST("recovermargin", routes.RecoverMargin)
	r.POST("recovermarginvalues", routes.RecoverMarginValues)
	r.POST("getmarginvaluesdata", routes.GetMarginValuesData)

	r.GET("getmargintabledata", routes.GetMarginTableData)
	r.GET("getmargin", routes.GetMargin)
	r.GET("getmarginvaluestabledata", routes.GetMarginValuesTableData)

	r.POST("adddeal", routes.AddDeal)
	r.POST("alldeals", routes.GetAllDealsData)
	r.POST("approvedeal", routes.ApproveDeal)
	r.POST("dealdetaileddata", routes.DealDetailedData)
	r.POST("getquotationdeal", routes.GetQuotationDeal)
	r.POST("updatequotation", routes.UpdateQuotation)
	r.POST("updatedealstatus", routes.UpdateDealStatus)
	r.POST("denydeal", routes.DenyDeal)

	// r.POST("updatefinaloutcome", routes.UpdateFinalOutCome)
	// r.POST("updateStatus", routes.UpdateDealFinalStatus)

	r.GET("getuplodeddeal", routes.GetQuotationUploadedData)
	r.GET("getapprovaldeal", routes.GetApprovalDeal)
	r.GET("getattendees", routes.Attendees)
	r.POST("getContactTableData", routes.GetContactTableData)
	r.POST("getcompanytabledata", routes.GetCompanyTableData)

	//dashboard
	r.POST("dashboard", routes.Dashboard)
	r.POST("notificationdata", routes.GetNotificationsNumber)
	r.POST("updategemnotification", routes.UpdateGemNotification)
	r.POST("markasread", routes.MarkAsRead)

	r.GET("usernamelistdashboard", routes.UserNameListDashboards)
	r.GET("getuserdata", routes.GetUsersData)

	r.POST("addchat", routes.InsertChat)
	r.POST("getchat", routes.GetChat)

	r.POST("addpartner", routes.AddPartner)
	r.POST("addcontactperson", routes.AddContactPerson)
	r.POST("addgem", routes.AddGem)
	r.POST("getdetailedgem", routes.GetGem)
	r.POST("uploadpricinggem", routes.UploadQuotationGem)
	r.POST("approvegem", routes.ApproveGeM)
	r.POST("denygem", routes.DenyGeM)
	r.POST("mygem", routes.GetMyGem)
	r.POST("updategemstatus", routes.UpdateStatus)
	r.POST("partnerdata", routes.GetPartner)
	r.POST("getproduct", routes.GetProductImage)

	r.GET("getpartner", routes.GetPartners)
	r.GET("getpartnerlist", routes.GetPartnersList)
	r.GET("getcontactpersonlist", routes.GetDataContactPerson)
	r.GET("getawaitinggem", routes.GetAwaitingPricingGEMList)
	r.GET("getupdatedpricinggem", routes.GetUpdatedPricingGEMList)
	r.GET("waitingforapprovalgem", routes.GetGemWaitingForApproval)

	//Export Data
	r.POST("exportdealdata", routes.ExportData)
	r.POST("sendgmail", routes.SendCustomerMail)
	// r.POST("/refreshToken", routes.RefreshTokenHandler)
	r.POST("save-token", routes.StoreFCMToken)

	r.POST("compliance", routes.GetComplianceData)
	r.POST("adddatasheet", routes.AddDataSheet)

	r.POST("getdatasheetdata", routes.GetDataSheetData)
	r.POST("addprocessor", routes.AddProcessor)
	r.POST("getdatasheetdatas", routes.GetDataSheetDatas)

	r.GET("getprocessor", routes.GetProcessor)
	r.GET("getdatasheetname", routes.GetDataSheetName)

	//Demo
	r.POST("adddemo", routes.AddDemo)
	r.POST("updatedemo", routes.UpdateDemo)
	r.POST("getdemo", routes.GetDemo)
	r.POST("detaileddemo", routes.GetDetailedDemo)

	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()
	go queue.ConsumeMailQueue()
	go queue.StartEmailWorker()

	select {}
}

func checkTasksEveryMinute() {
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ticker.C:
			routes.CheckAndNotifyDueTasks()
		}
	}
}
