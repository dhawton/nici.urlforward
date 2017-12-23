package main

import (
  "net/http"
  "html/template"
  "github.com/julienschmidt/httprouter"
  "os"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "log"
  "fmt"
  "regexp"
)

var db *sql.DB

type RedirectVars struct {
  URL string
}

type Forward struct {
  url_forward string
  url_forward_mask int
}

func getDomain(domain string) string {
  r := regexp.MustCompile("^www.")
  domain = r.ReplaceAllString(domain, "")
  r = regexp.MustCompile(":[0-9]+")
  domain = r.ReplaceAllString(domain, "")
  return domain
}

func redirect(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
  var forward Forward
  fmt.Println(getDomain(r.Host))

  row := db.QueryRow("SELECT url_forward, url_forward_mask FROM `domains` WHERE `domain`='" + getDomain(r.Host) + "' LIMIT 1")
  err := row.Scan(&forward.url_forward, &forward.url_forward_mask)
  switch err {
  case sql.ErrNoRows:
  case nil:
  default:
    fmt.Println(err)
    return
  }

  if forward.url_forward_mask == 1 {
    http.Redirect(w, r, forward.url_forward, 303)
    return
  }

  // No redirect, show template
  p := RedirectVars{URL: forward.url_forward}
  t, _ := template.ParseFiles("template.html")
  t.Execute(w, p)
}

func getDSN() string {
  return os.Getenv("MYSQL_USERNAME") + ":" + os.Getenv("MYSQL_PASSWORD") +
    "@tcp(" + os.Getenv("MYSQL_HOST") + ")/" + os.Getenv("MYSQL_DATABASE")
}

func main () {
  router := httprouter.New()
  router.GET("/:whatever", redirect)
  var err error
  db, err = sql.Open("mysql", getDSN())
  if err != nil {
    log.Fatal(err)
  }

  log.Fatal(http.ListenAndServe(":8080", router))
}
