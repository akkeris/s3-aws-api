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
			fmt.Println("Error: Unable to run migration scripts, could not load create.sql.")
			os.Exit(1)
		}
	}
	defer db.Close()
	_, err = db.Exec(string(buf))
	if err != nil {
		fmt.Println("Error: Unable to run migration scripts, execution failed.")
		os.Exit(1)
	}
}

func Init(brokerdburl string) {
	db, dberr := sql.Open("postgres", brokerdburl)
	if dberr != nil {
		fmt.Println(dberr)
	}
	db.SetMaxIdleConns(4)
	db.SetMaxOpenConns(20)
	CreateDB(db)
	brokerdb = db

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

func Retrieve(bucketname string) (l string, a string, s string) {
	getDBStats()
	stmt, err := brokerdb.Prepare("select  location, accesskey_enc, secretkey_enc from provision where bucketname = $1 ")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(bucketname)
	defer rows.Close()
	var location string
	var accesskey_enc string
	var secretkey_enc string
	for rows.Next() {
		err := rows.Scan(&location, &accesskey_enc, &secretkey_enc)
		if err != nil {
			fmt.Println(err)
		}
	}
	getDBStats()
	return location, utils.Decrypt(accesskey_enc), utils.Decrypt(secretkey_enc)

}

func getDBStats() {
	fmt.Printf("Number of Open Connections: %d\n", brokerdb.Stats().OpenConnections)

}
