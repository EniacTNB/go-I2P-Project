package routers

import (
	"github.com/gin-gonic/gin"

	"gin-vue/middleware/cors"
	"gin-vue/middleware/myjwt"
	"gin-vue/pkg/e"
	"gin-vue/pkg/setting"
	v1 "gin-vue/routers/api/v1"
	v2 "gin-vue/routers/api/v2"
	i2p "gin-vue/routers/i2p"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(cors.CorsHandler())
	r.Use(gin.Recovery())
	gin.SetMode(setting.RunMode)
	var authMiddleware = myjwt.GinJWTMiddlewareInit(&myjwt.AllUserAuthorizator{})
	r.POST("/login", authMiddleware.LoginHandler)
	//404 handler
	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		code := e.PAGE_NOT_FOUND
		c.JSON(404, gin.H{"code": code, "message": e.GetMsg(code)})
	})
	auth := r.Group("/auth")
	{
		// Refresh time can be longer than token timeout
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	api := r.Group("/user")
	api.Use(authMiddleware.MiddlewareFunc())
	{
		api.GET("/info", v1.GetUserInfo)
		api.POST("/logout", v1.Logout)
	}

	var adminMiddleware = myjwt.GinJWTMiddlewareInit(&myjwt.AdminAuthorizator{})
	apiv1 := r.Group("/api/v1")
	//使用AdminAuthorizator中间件，只有admin权限的用户才能获取到接口
	apiv1.Use(adminMiddleware.MiddlewareFunc())
	{
		//vue获取table信息
		apiv1.GET("/table/list", v2.GetArticles)
		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)

		//获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiv1.POST("/articles", v1.AddArticle)
		//更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	apiI2p := r.Group("/api/i2p")
	apiI2p.Use()
	{
		apiI2p.GET("filelists", i2p.GetFileList)
		apiI2p.GET("filecontent", i2p.GetFileContent)
		apiI2p.POST("saveluafile", i2p.SaveLuaFile)
		apiI2p.POST("runExper", i2p.RunLuaFile)
		apiI2p.GET("testKafkaParse", i2p.KafkaConsumeSample)
		apiI2p.GET("metadata", i2p.GetMetaData)
		apiI2p.GET("metadata_content", i2p.GetMetaDataContent)
		apiI2p.GET("exper/list", i2p.GetExperList)
		apiI2p.GET("exper/filecontent", i2p.GetFileContent)
		apiI2p.GET("exper/getLog", i2p.GetDockerLogs)
		apiI2p.GET("exper/getDockerState", i2p.GetDockerStateTest)
		apiI2p.GET("exper/getFileChange", i2p.GetNetDbChangeInDocker)
		apiI2p.GET("exper/getFileList", i2p.GetNetDbFileList)

	}

	var testMiddleware = myjwt.GinJWTMiddlewareInit(&myjwt.TestAuthorizator{})
	apiv2 := r.Group("/api/v2")
	apiv2.Use(testMiddleware.MiddlewareFunc())
	{
		//获取文章列表
		apiv2.GET("/articles", v2.GetArticles)
	}

	return r
}
