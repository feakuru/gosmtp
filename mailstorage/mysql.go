package mailstorage

import (
    "log"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

// emails: id | from_id | text
// users: id | names
// recipients: email_id | user_id

var db *sql.DB
var err error

func Init() {
    db, err = sql.Open("mysql", "gosmtp:passgosmtp@tcp(127.0.0.1:3306)/gosmtp")
    if err != nil {
        log.Fatal(err)
    }
    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }
}

func Exit() {
    db.Close()
}

func CreateUser(name string) {
    stmt, err := db.Prepare("INSERT INTO users(name) VALUES(?)")
    if err != nil {
        log.Fatal(err)
    }
    res, err := stmt.Exec(name)
    if err != nil {
        log.Fatal(err)
    }
    lastId, err := res.LastInsertId()
    if err != nil {
        log.Fatal(err)
    }
    rowCnt, err := res.RowsAffected()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
}

func PutEmailToDB(from string, rcpt []string, text string) {
    var from_id int
    from_rows, err := db.Query("select id from users where name = ?", from)
    if err != nil {
        log.Fatal(err)
    }
    defer from_rows.Close()
    for from_rows.Next() {
        err := from_rows.Scan(&from_id)
        if err != nil {
            log.Fatal(err)
        }
        log.Println("FROM", from_id, from)
    }
    err = from_rows.Err()
    if err != nil {
        log.Fatal(err)
    }

    rcpt_ids := make([]int, len(rcpt))
    for rcpt_num := 0; rcpt_num < len(rcpt); rcpt_num++ {
        var rcpt_id int
        rcpt_rows, err := db.Query("select id from users where name = ?", rcpt[rcpt_num])
        if err != nil {
            log.Fatal(err)
        }
        defer rcpt_rows.Close()
        for rcpt_rows.Next() {
            err := rcpt_rows.Scan(&rcpt_id)
            if err != nil {
                log.Fatal(err)
            }
            rcpt_ids[rcpt_num] = rcpt_id
            log.Println("TO", rcpt_id, rcpt[rcpt_num])
        }
        err = rcpt_rows.Err()
        if err != nil {
            log.Fatal(err)
        }
    }

    var emailId int64
    stmt, err := db.Prepare("INSERT INTO emails(from_id, text) VALUES(?, ?)")
    if err != nil {
        log.Fatal(err)
    }
    res, err := stmt.Exec(from_id, text)
    if err != nil {
        log.Fatal(err)
    }
    lastId, err := res.LastInsertId()
    if err != nil {
        log.Fatal(err)
    }
    rowCnt, err := res.RowsAffected()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
    emailId = lastId

    for rcpt_num := 0; rcpt_num < len(rcpt_ids); rcpt_num++ {
        stmt, err := db.Prepare("INSERT INTO recipients(email_id, user_id) VALUES(?, ?)")
        if err != nil {
            log.Fatal(err)
        }
        res, err := stmt.Exec(emailId, rcpt_ids[rcpt_num])
        if err != nil {
            log.Fatal(err)
        }
        lastId, err := res.LastInsertId()
        if err != nil {
            log.Fatal(err)
        }
        rowCnt, err := res.RowsAffected()
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)
    }
}
