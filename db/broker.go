package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	utils "s3-aws-api/utils"
        "log"
        "io/ioutil"
)

var brokerdb *sql.DB

func Init(brokerdburl string) {
	db, dberr := sql.Open("postgres", brokerdburl)
	if dberr != nil {
		fmt.Println(dberr)
	}
	db.SetMaxIdleConns(4)
	db.SetMaxOpenConns(20)
	brokerdb = db
        createDB(brokerdb)

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

func createDB(db *sql.DB) {
        buf, err := ioutil.ReadFile("./create.sql")
        if err != nil {
                buf, err = ioutil.ReadFile("../create.sql")
                if err != nil {
                        log.Println("Error: Unable to run migration scripts, could not load create.sql.")
                        log.Fatalln(err)
                }
        }
        _, err = db.Query(string(buf))
        if err != nil {
                log.Println("Error: Unable to run migration scripts, execution failed.")
                log.Fatalln(err)
        }
}

