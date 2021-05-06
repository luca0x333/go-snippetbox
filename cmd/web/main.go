package main

import (
	"database/sql"
	"flag"
	"github.com/luca0x333/go-snippetbox/pkg/models/mysql"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// Default port 4000
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web_user:password@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	// New creates a new Logger. The out variable sets the
	// destination to which log data will be written.
	// The prefix appears at the beginning of each generated log line.
	// The flag argument defines the logging properties.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.LUTC)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize a new templateCache
	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a new instance of application.
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(), // Call routes() method.
	}

	// flag.String() returns a pointer.
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// openDB() wraps sql.Open() and return *sql.DB or an error
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
