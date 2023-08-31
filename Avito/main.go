package main

import (
    "github.com/gin-gonic/gin"
    "Avito/service/config"
    "Avito/service"
    "strconv"
)

func main() {
    Base, err := service.DataBaseConnect(config.DBhost, config.DBuser, config.DBpassword, config.DBname, config.DBport)
    defer Base.Close()

    if err != nil {
        panic(err)
    } 
    
    service.CreateAccountsTable(Base)
    service.CreateEventsTable(Base)
    service.CreateExpirationTable(Base)

    router := gin.Default()

    router.GET("", func(context *gin.Context) { 
        service.GetUserSlugsHTTP(context, Base)
    })

    router.GET("/report/:name", func(context *gin.Context) {
        service.GetCSV(context, Base)        
    })

    router.GET("/report/", func(context *gin.Context) {
        service.GetURL(context, Base)
    })

    router.POST("", func(context *gin.Context) {
        service.AddSlugHTTP(context, Base, config.Users)
    })

    router.PUT("", func(context *gin.Context) {
        service.AddUserToSlugHTTP(context, Base)
    })

    router.DELETE("", func(context *gin.Context) {
        service.DeleteSlugHTTP(context, Base)
    })

    router.DELETE("/report/", func(context *gin.Context) {
        service.DeleteCSV() 
    })

    go service.TTL(Base)

    router.Run(config.HTTPhost + ":" + strconv.Itoa(config.HTTPport))
}
