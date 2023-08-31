package service

import (
    "github.com/gin-gonic/gin"
    "Avito/service/structs"
    "Avito/service/config"
    "database/sql"
    "net/http"
    "strconv"
    "bufio"
    "time"
    "os"
)

func GetUserSlugsHTTP(context *gin.Context, db *sql.DB) error {
    var user structs.User

    if err := context.BindJSON(&user); err != nil {
        panic(err)
    }

    query, err := GetUserSlugs(db, user.ID)

    if err != nil {
        return err
    }

    var slugs []string
    for query.Next() {
        var slug string

        query.Scan(&slug)

        slugs = append(slugs, slug)
    }

    context.IndentedJSON(http.StatusOK, slugs)

    return err
}

func AddUserToSlugHTTP(context *gin.Context, db *sql.DB) {
    var newUser structs.UserToSlug

    if err := context.BindJSON(&newUser); err != nil {
        panic(err)
    }

    DateTime := time.Now()

    for _, del := range newUser.ToDelete {
        if !SlugExists(db, del) {
            continue
        }

        if UserInSlug(db, del, newUser.ID) {
            DeleteUserFromSlug(db, del, newUser.ID)
        }
    }

    for i, add := range newUser.ToAdd {
        if !SlugExists(db, add) {
            continue
        }

        if !UserInSlug(db, add, newUser.ID) {
            if newUser.Exp[i].Y != "" && newUser.Exp[i].Mo != "" && newUser.Exp[i].D != "" && newUser.Exp[i].H != "" && newUser.Exp[i].Mi != "" {
                if newUser.Exp[i].Y != "" {
                    newUser.Exp[i].Y = "0"
                }

                if newUser.Exp[i].Mo != "" {
                    newUser.Exp[i].Mo = "0"                                                                                                                                                                  
                }

                if newUser.Exp[i].D != "" {
                    newUser.Exp[i].D = "0"
                }

                if newUser.Exp[i].H != "" {
                    newUser.Exp[i].H = "0"
                }

                if newUser.Exp[i].Mi != "" {
                    newUser.Exp[i].Mi = "0"                                                                                                                                                                                 }

                Year,    _ := strconv.Atoi(newUser.Exp[i].Y)
                Month,   _ := strconv.Atoi(newUser.Exp[i].Mo)
                Days,    _ := strconv.Atoi(newUser.Exp[i].D)
                Hours,   _ := strconv.Atoi(newUser.Exp[i].H)
                Minutes, _ := strconv.Atoi(newUser.Exp[i].Mi)

                ExpDateTime := DateTime.AddDate(Year, Month, Days)
                duration := time.Duration(Hours) * time.Hour + time.Duration(Minutes) * time.Minute
                ExpDateTime = ExpDateTime.Add(duration)

                AddExpDate(db, newUser.ID, add, ExpDateTime)
            }

            AddUserToSlug(db, add, newUser.ID, DateTime)
        }
    }
}

func AddSlugHTTP(context *gin.Context, db *sql.DB, CountOfUsers int) error {
    var newSlug structs.Slug_ser
   
    if err := context.BindJSON(&newSlug); err != nil {
        panic(err)
    }

    if SlugExists(db, newSlug.Slug) {
        return nil
    }

    err := CreateSlug(db, newSlug.Slug)
    
    if err != nil {
        panic(err)
    }

    Ids := GetPercentOfUsers(CountOfUsers, newSlug.Percent)
    
    AddPercentToSlug(db, newSlug.Slug, Ids)

    return err
}

func DeleteSlugHTTP(context *gin.Context, db *sql.DB) error {
    var newSlug structs.Slug_ser
      
    if err := context.BindJSON(&newSlug); err != nil {
        panic(err)
    }

    return DeleteSlug(db, newSlug.Slug)
}

func GetURL(context *gin.Context, db *sql.DB) {
    var data structs.UserEvents

    context.BindJSON(&data)

    query, _ := db.Query("SELECT * FROM events WHERE EXTRACT(MONTH FROM date_time) = " + data.Month + " AND EXTRACT(YEAR FROM date_time) = " + data.Year + " AND id = " + data.ID)

    name := "http://" + config.HTTPhost + ":" + strconv.Itoa(config.HTTPport) + "/report/" +  MakeCSVFromSQLData(query)
    
    context.IndentedJSON(http.StatusOK, structs.CSV_link{URL: name})
}

func GetCSV(context *gin.Context, db *sql.DB) {
    fileName := context.Param("name")
    
    path := "/app/"

    file, _ := os.Open(path + fileName + ".csv")
    defer file.Close()
    
    str := make([]string, 0);
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        str = append(str, scanner.Text())
    }

    context.IndentedJSON(http.StatusOK, str)
}
