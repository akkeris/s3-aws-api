package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	utils "s3-aws-api/utils"

	_ "github.com/lib/pq"
)

var brokerdb *sql.DB

func CreateDB(db *sql.DB) {
	buf, err := ioutil.ReadFile("./create.sql")
	if err != nil {
		buf, err = ioutil.ReadFile("../create.sql")
		if err != nil {
			db.Close()
			fmt.Println("Error: Unable to run migration scripts, could not load create.sql.")
			os.Exit(1)
		}
	}

	_, err = db.Exec(string(buf))
	if err != nil {
		db.Close()
		fmt.Println("Error: Unable to run migration scripts, execution failed.")
		os.Exit(1)
	}
}

func Init(brokerdburl string) *sql.DB {
	db, dberr := sql.Open("postgres", brokerdburl)
	if dberr != nil {
		fmt.Println(dberr)
	}
	db.SetMaxIdleConns(4)
	db.SetMaxOpenConns(20)
	CreateDB(db)
	brokerdb = db
	return brokerdb
}
func Delete(bucketname string) {

	getDBStats()

	stmt, err := brokerdb.Prepare("delete from provision where bucketname=$1")
	if err != nil {
		fmt.Println(err)
	}
	res, err := stmt.Exec(bucketname)
	if err != nil {
		fmt.Println(err)
	}
	affect, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(affect, "rows changed")
	getDBStats()

}
func Store(bucketname string, location string, accesskey string, secretkey string, billingcode string) {
	getDBStats()
	var newname string
	err := brokerdb.QueryRow("INSERT INTO provision(bucketname, location, accesskey_enc, secretkey_enc, billingcode) VALUES($1,$2,$3,$4,$5) returning bucketname;", bucketname, location, utils.Encrypt(accesskey), utils.Encrypt(secretkey), billingcode).Scan(&newname)

	if err != nil {
		fmt.Println(err)
	}
	getDBStats()
}

func Retrieve(bucketname string) (l string, a string, s string, e error) {
	getDBStats()

	var location string
	var accesskey_enc string
	var secretkey_enc string

	err := brokerdb.QueryRow("SELECT location, accesskey_enc, secretkey_enc FROM provision WHERE bucketname = $1", bucketname).Scan(&location, &accesskey_enc, &secretkey_enc)
	if err != nil {
		return "", "", "", err
	}

	getDBStats()
	return location, utils.Decrypt(accesskey_enc), utils.Decrypt(secretkey_enc), nil
}

func getDBStats() {
	fmt.Printf("Number of Open Connections: %d\n", brokerdb.Stats().OpenConnections)
}

func ValidBucket(bucketname string) bool {
	var exists bool
	err := brokerdb.QueryRow("SELECT EXISTS (SELECT FROM provision WHERE bucketname = $1)", bucketname).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}
