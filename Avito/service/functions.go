package service

import (
    "Avito/service/structs"
    "path/filepath"
    "encoding/csv"
    "database/sql"
    "math/rand"
    "strings"
    "time"
    "fmt"
    "os"
)

func GetPercentOfUsers(CountOfUsers int, Percent float32) []int {       //remake it for O(n^2)
    rand.Seed(time.Now().UnixNano())
    n := int(float32(CountOfUsers) * Percent / 100.0)

    usersSet := make(map[int]structs.Void, n)
    users := make([]int, n, n)

    if n == CountOfUsers {
        for i, _ := range users {
            users[i] = i
        }

        return users
    }

    for i, _ := range users {
        
        exists := true
        var rnd int 

        for {
            rnd = rand.Intn(CountOfUsers)
            _, exists = usersSet[rnd]

            if !exists {
                break
            }
        }

        users[i] = rnd
        usersSet[rnd] = structs.Void{}
    }

    return users
}

func NameGenerator(n int) string {
    rand.Seed(time.Now().UnixNano())

    name := ""

    for i := 0; i < n; i++ {
        name = name + fmt.Sprintf("%c", 65 +  rand.Intn(26))
    }

    return name
}

func MakeCSVFromSQLData(query *sql.Rows) string {
    var id, slug, event, date string

    path := "/app/"

    name := NameGenerator(6)

    file, _ := os.Create(path + name + ".csv")
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    for query.Next() {
        query.Scan(&id, &slug, &event, &date)
        date = strings.ReplaceAll(date, "T", " ")
        date = strings.ReplaceAll(date, "Z", " ")
        writer.Write([]string{id, slug, event, date})
    }

    return name
}

func DeleteCSV() {
    path := "/app/"

    pattern := filepath.Join(path, "*.csv")

    files, _ := filepath.Glob(pattern)

    for _, val := range files {
        os.Remove(val)
    }
}

func TTL(db *sql.DB) {
    ticker := time.NewTicker(time.Minute)

    for ticktime := range ticker.C {
        query, _ := db.Query("SELECT id, slug FROM expiration WHERE date_time <= '" + ticktime.Format("2006-01-02 15:04:05") + "'")

        for query.Next() {
            var id, slug string

            query.Scan(&id, &slug)

            DeleteUserFromSlug(db, slug, id)
        }

        db.Query("DELETE FROM expiration WHERE date_time <= '" + ticktime.Format("2006-01-02 15:04:05") + "'")
    }
}
