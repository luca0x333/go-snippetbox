package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"github.com/golangcollege/sessions"
	"github.com/luca0x333/go-snippetbox/pkg/models/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// Default port 4000
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web_user:password@/snippetbox?parseTime=true", "MySQL data source name")
	secret := flag.String("secret", "z6Nah+pPonzHbI*+9Pk8qNWhTzbpa@ge", "Secret Key")

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

	// Initialize a new session manager.
	// flag.String() returns a pointer.
	session := sessions.New([]byte(*secret))
	// Lifetime sets the maximum length of time that a session is valid for before it expires.
	session.Lifetime = 12 * time.Hour

	// Initialize a new instance of application.
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// Initialize a new tls.Config struct to overwrite the default TLS settings we want to change.
	tlsConfig := &tls.Config{
		// By setting PreferServerCipherSuites to "true" Go's cipher suites are preferred over the user cipher suites.
		PreferServerCipherSuites: true,
		// CurvePreferences field lets us specify which elliptic curves
		// should be given preference during the TLS handshake.
		//
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:      *addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(), // Call routes() method.
		TLSConfig: tlsConfig,
		// Add Idle, Read and Write timeouts to the server.
		// If you set ReadTimeout but don’t set IdleTimeout,
		// then IdleTimeout will default to using the same setting as ReadTimeout.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// flag.String() returns a pointer.
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
