package service

import (
    "fmt"
    "log"
    "time"
    "strconv"
    "database/sql"
    "github.com/lib/pq"
)

func AccountsExists(db *sql.DB, table string) bool {
    query, _ := db.Query("SELECT EXISTS ( SELECT 1 FROM information_schema.tables WHERE table_name = '" + table + "');")

    var flag bool
    for query.Next() {
        err := query.Scan(&flag)

        if err != nil {
            log.Printf("Error occured: %v", err)
        }
    }

    return flag
}

func DataBaseConnect(host, user, password, dbname string, port int) (*sql.DB, error) {
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

    return sql.Open("postgres", psqlInfo)
}

func CreateAccountsTable(db *sql.DB) error {
    _, err := db.Query("CREATE TABLE accounts (slug VARCHAR(50), id INT[])")
    return err
}   

func CreateEventsTable(db *sql.DB) error {
    _, err := db.Query("CREATE TABLE events (id INT, slug VARCHAR(50), event_name VARCHAR(10), date_time TIMESTAMP)")
    
    return err
}

func CreateExpirationTable(db *sql.DB) error {
    _, err := db.Query("CREATE TABLE expiration (id INT, slug VARCHAR(50), date_time TIMESTAMP)")
    
    return err
}

func CreateSlug(db *sql.DB, slug_tag string) error {
    _, err := db.Query("INSERT INTO accounts(slug, id) VALUES ('" + slug_tag + "', '{}')")

    return err
}

func AddUserToSlug(db *sql.DB, slug_tag, id string, date_time time.Time) error {

    err := AddEvent(db, id, slug_tag, "added", date_time.Format("2006-01-02 15:04:05"))
    
    if err != nil {
        log.Printf("Error occured: %v", err)
    }

    _, err = db.Query("UPDATE accounts SET id = array_append(id, " + id + ") WHERE slug = '" + slug_tag + "'")
    
    return err
}

func AddUsersToSlug(db *sql.DB, slug_tag string, id []string) error {
    var ids string

    for i, _ := range id {
        ids = ids + id[i]
        if i != len(id) - 1 {
            ids = ids + ", "
        } 
    }

    _, err := db.Query("UPDATE accounts SET id = id || '{" + ids + "}' WHERE slug = '" + slug_tag + "'") 
    return err
}

func DeleteSlug(db *sql.DB, slug_tag string) error {      
    var users []int32
 
    query, err := db.Query("SELECT id FROM accounts WHERE slug = '" + slug_tag + "'")

    if err != nil {
        log.Printf("Error occured: %v", err)
    }

    _, err = db.Query("DELETE FROM accounts WHERE slug = '" + slug_tag + "'")   
    
    if err != nil {
        log.Printf("Error occured: %v", err)
    }

    for query.Next() {
        query.Scan(pq.Array(&users))
    }
    
    query, _ = db.Query("SELECT id FROM expiration WHERE slug = '" + slug_tag + "'")

    idMap := make(map[int]bool)

    for query.Next() {
        var val int

        query.Scan(&val)

        idMap[val] = true
    }

    usersStr := make([]string, 0)
    for i, _ := range users {
        if _, exists := idMap[int(users[i])]; !exists {
            usersStr = append(usersStr, strconv.Itoa(int(users[i])))
        }
    }

    AddEvents(db, slug_tag, "deleted", time.Now().Format("2006-01-02 15:04:05"), usersStr)  

    return err
}

func DeleteUserFromSlug(db *sql.DB, slug_tag, id string) error {
    
    err := AddEvent(db, id, slug_tag, "deleted", time.Now().Format("2006-01-02 15:04:05"))
    
    if err != nil {
        log.Printf("Error occured: %v", err)
    }
 
    _, err = db.Query("UPDATE accounts SET id = (SELECT ARRAY_REMOVE(id, " + id + ") FROM accounts WHERE slug = '"+ slug_tag + "') WHERE slug = '" + slug_tag + "'")
    
    if err != nil {
        log.Printf("Error occured: %v", err)
    }

    return err
}

func GetUserSlugs(db *sql.DB, id string) (*sql.Rows, error) {
    return db.Query("SELECT slug FROM accounts WHERE " + id + " = ANY(id)")
}

func GetUserExpSlugs(db *sql.DB, id string) (*sql.Rows, error) {
    return db.Query("SELECT slug FROM expiration WHERE id = " + id)
}

func AddPercentToSlug(db *sql.DB, slug_tag string, users []int) {
    strUsers := make([]string, len(users), len(users))

    for i, _ := range users {
        strUsers[i] = strconv.Itoa(users[i])
    }

    AddEvents(db, slug_tag, "added", time.Now().Format("2006-01-02 15:04:05"), strUsers)
    AddUsersToSlug(db, slug_tag, strUsers)
}

func UserInSlug(db *sql.DB, slug, id string) bool {
    query, _ := db.Query("SELECT COUNT(*) FROM accounts WHERE " + id + " = ANY(id) AND slug = '" + slug + "'")
    
    var exist bool

    for query.Next() {
        query.Scan(&exist)
    }

    return exist
}

func SlugExists(db *sql.DB, slug string) bool {
    query, _ := db.Query("SELECT COUNT(*) FROM accounts WHERE slug = '" + slug + "'" )

    var exist bool

    for query.Next() {
        query.Scan(&exist)
    }

    return exist
}

func AddEvent(db *sql.DB, id, slug, event_name, date_time string) error {
    _, err := db.Query("INSERT INTO events(id, slug, event_name, date_time) VALUES (" + id + ", '" + slug + "', '" + event_name + "', '" + date_time + "')")

    return err
}

func AddEvents(db *sql.DB, slug, event_name, date_time string, id []string) {
    query := "INSERT INTO events(id, slug, event_name, date_time) VALUES " 

    for i, _ := range id {
        query = query + "(" + id[i] + ", '" + slug + "', '" + event_name + "', '" + date_time + "')"

        if i != len(id) - 1 {
            query = query + ", "
        }
    }
    
    db.Query(query)
}

func AddExpDate(db *sql.DB, id, slug_tag string, date time.Time) error {
    _, err := db.Query("INSERT INTO expiration(id, slug, date_time) VALUES (" + id + ", '" + slug_tag + "', '" + date.Format("2006-01-02 15:04:05") + "')")

    return err
}
