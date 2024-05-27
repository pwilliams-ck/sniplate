package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	// Import the pq driver so that it can register itself with the database/sql
	// package. We alias this import to the blank identifier, to stop the Go
	// compiler complaining that the package isn't being used.
	_ "github.com/lib/pq"
)

const version = "0.1.5"

type config struct {
	port   int
	env    string
	useTLS bool
	useLog bool
	db     struct {
		dsn string
	}
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	var cfg config

	// port defines the port number for the API server.
	// Defaults to 4200 if not provided via CLI.
	flag.IntVar(&cfg.port, "port", 4200, "API server port")

	// env represents the current environment the application is running in.
	// Valid values are: development, staging, production.
	// Defaults to "development" if not set via CLI.
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	// Boolean useTLS gives the option to enable TLS.
	// Defaults to false, use true for production.
	flag.BoolVar(&cfg.useTLS, "tls", false, "Enable TLS (true|false)")

	// Boolean useLog gives the option to enable logging to a file, as well as the usual stdout.
	// Defaults to false, use true for production.
	flag.BoolVar(&cfg.useLog, "log", false, "Enable log file (true|false)")

	// Read the DSN value from the db-dsn command-line flag into the config struct. We
	// default to using our development DSN if no flag is provided.
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:password@postgres/snips?sslmode=disable", "PostgreSQL DSN")
	// We need to parse all CLI flags in order to use them as well.
	flag.Parse()

	// Logging setup
	logWriter := setupLogger(cfg.useLog)

	// Create a new logger that writes to standard output (os.Stdout).
	// Logger is configured with a text handler that formats log records as plain text.
	// nil argument specifies that no additional handler options are provided.
	logger := slog.New(slog.NewTextHandler(logWriter, nil))

	// Call the openDB() helper function (see below) to create the connection pool,
	// passing in the config struct. If this returns an error, we log it and exit the
	// application immediately.
	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Defer db.Close() so that the connection pool is closed before the main() function exits.
	defer db.Close()

	logger.Info("database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
	}

	// TLS Config is set up for modern web, maybe remove some of these settings if needed.
	// TLS 1.3 remains unaffected by all of this, as all of its connections are considered
	// safe while writing this for Go 1.22.
	tlsConfig := &tls.Config{
		MinVersion:       tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	// `srv` initializes an HTTP server with defined configuration for address, handlers,
	// TLS settings, timeouts, and error logging.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env, "tls", cfg.useTLS, "log", cfg.useLog)

	// Start server with or without TLS.
	if cfg.useTLS {
		srv.TLSConfig = tlsConfig
		err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	} else {
		err := srv.ListenAndServe()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}
}

// The openDB() function returns a sql.DB connection pool.
func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config
	// struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use PingContext() to establish a new connection to the database, passing in the
	// context we created above as a parameter. If the connection couldn't be
	// established successfully within the 5 second deadline, then this will return an
	// error. If we get this error, or any other, we close the connection pool and
	// return the error.
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	// Return the sql.DB connection pool.
	return db, nil
}
