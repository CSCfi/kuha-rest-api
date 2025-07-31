package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// Database struct to hold multiple connections
type Database struct {
	FIS        *sql.DB
	UTV        *sql.DB
	Auth       *sql.DB
	Tietoevry  *sql.DB
	KAMK       *sql.DB
	KLAB       *sql.DB
	ARCHINISIS *sql.DB
}

func NewSingleDB(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	db, err := connectDB(addr, maxOpenConns, maxIdleConns, maxIdleTime)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func New(fisAddr, utvAddr, authAddr, tietoevryAddr, kamkAddr, klabAddr, archinisisAddr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*Database, error) {
	fisDB, err := connectDB(fisAddr, maxOpenConns, maxIdleConns, maxIdleTime)
	if err != nil {
		return nil, err
	}

	utvDB, err := connectDB(utvAddr, maxOpenConns, maxIdleConns, maxIdleTime)
	if err != nil {
		return nil, err
	}

	authDB, err := connectDB(authAddr, maxOpenConns, maxIdleConns, maxIdleTime)
	if err != nil {
		return nil, err
	}

	tietoevryDB, err := connectDB(tietoevryAddr, maxOpenConns, maxIdleConns, maxIdleTime)
	if err != nil {
		return nil, err
	}

	kamkDB, err := connectDB(kamkAddr, maxOpenConns, maxIdleConns, maxIdleTime)
	if err != nil {
		return nil, err
	}

	klabDB, err := connectDB(klabAddr, maxOpenConns, maxIdleConns, maxIdleTime)
	if err != nil {
		return nil, err
	}

	archinisisDB, err := connectDB(archinisisAddr, maxOpenConns, maxIdleConns, maxIdleTime)
	if err != nil {
		return nil, err
	}

	return &Database{FIS: fisDB, UTV: utvDB, Auth: authDB, Tietoevry: tietoevryDB, KAMK: kamkDB, KLAB: klabDB, ARCHINISIS: archinisisDB}, nil
}

func connectDB(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(time.Duration(duration))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func NewWithGracefulFailure(
	fisAddr, utvAddr, authAddr, tietoevryAddr, kamkAddr, klabAddr, archinisisAddr string,
	maxOpenConns, maxIdleConns int, maxIdleTime string,
) (*Database, map[string]error) {
	databases := &Database{}
	errors := make(map[string]error)

	tryConnect := func(name, addr string) *sql.DB {
		if addr == "" {
			errors[name] = fmt.Errorf("not configured")
			return nil
		}
		db, err := connectToDB(addr, maxOpenConns, maxIdleConns, maxIdleTime)
		if err != nil {
			errors[name] = err
			return nil
		}
		errors[name] = nil
		return db
	}

	databases.FIS = tryConnect("fis", fisAddr)
	databases.UTV = tryConnect("utv", utvAddr)
	databases.Auth = tryConnect("auth", authAddr)
	databases.Tietoevry = tryConnect("tietoevry", tietoevryAddr)
	databases.KAMK = tryConnect("kamk", kamkAddr)
	databases.KLAB = tryConnect("klab", klabAddr)
	databases.ARCHINISIS = tryConnect("archinisis", archinisisAddr)

	return databases, errors
}

func connectToDB(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		db.Close()
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
