package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/nu7hatch/gouuid"
	octdb "s3-aws-api/db"
	octiam "s3-aws-api/iam"
	octs3 "s3-aws-api/s3"
	structs "s3-aws-api/structs"
	utils "s3-aws-api/utils"
        "strconv"
	"strings"
	"time"
	"os"
)

var region string

func main() {
	brokerdb := utils.Init()
	octdb.Init(brokerdb)
        region = os.Getenv("REGION")
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Get("/v1/s3/plans", plans)
	m.Post("/v1/s3/instance", binding.Json(structs.Provisionspec{}), provision)
	m.Post("/v1/tag", binding.Json(structs.Tagspec{}), tag)
	m.Get("/v1/s3/url/:name", url2)
	m.Delete("/v1/s3/instance/:name", delete)
	m.Run()

}

func url2(params martini.Params, r render.Render) {
	rlocation, raccesskey, rsecretkey := octdb.Retrieve(params["name"])

	r.JSON(200, map[string]string{"S3_BUCKET": params["name"], "S3_LOCATION": rlocation, "S3_ACCESS_KEY": raccesskey, "S3_SECRET_KEY": rsecretkey,"S3_REGION":region})
}

func createname() string {

	u, _ := uuid.NewV4()
	newusername := "u" + strings.Split(u.String(), "-")[0]
	newusername = os.Getenv("NAME_PREFIX")+"-" + newusername
	return newusername
}

func tag(spec structs.Tagspec, berr binding.Errors, r render.Render) {
	if berr != nil {
		errorout := make(map[string]interface{})
		errorout["error"] = berr
		r.JSON(500, errorout)
		return
	}
	octs3.Tagbucket(spec.Resource, spec.Name, spec.Value)
	r.JSON(201, map[string]interface{}{"response": "tag added"})
}
func plans(params martini.Params, r render.Render) {
	plans := make(map[string]interface{})
	plans["basic"] = "Simple Bucket - no versioning - no georeplication"
	plans["versioned"] = "Bucket - with versioning"
	r.JSON(200, plans)
}
func url(bucketname string) structs.S3spec {
	var s3spec structs.S3spec
	return s3spec
}
func provision(spec structs.Provisionspec, err binding.Errors, r render.Render) {
	var s3spec structs.S3spec
	basename := createname()
	s3user := octiam.Createuser(basename)
	bucketlocation := octs3.Createbucket(s3user.Username,spec.Plan)
	octs3.Tagbucket(basename, "billingcode", spec.Billingcode)
        sleeptimestring := os.Getenv("S3_BUCKET_POLICY_WAIT_SECONDS")
        if len(sleeptimestring) == 0 {
               sleeptimestring="10"
        }
        sleeptime, cerr := strconv.ParseInt(sleeptimestring, 10, 64)
        if cerr != nil {
            fmt.Println("Unable to get sleeptime")
        }
        time.Sleep( time.Second * time.Duration(sleeptime))
	octs3.Addbucketpolicy(s3user.Username, s3user)
	simpleuserpolicy := octiam.Createuserpolicy(s3user.Username, s3user.Username)
	octiam.Attachuserpolicy(s3user.Username, simpleuserpolicy)
	s3spec.Location = bucketlocation
        s3spec.Bucket = basename
	s3spec.Accesskey = s3user.Accesskey
	s3spec.Secretkey = s3user.Secretkey
        s3spec.Region = region
	octdb.Store(s3user.Username, s3spec.Location, s3spec.Accesskey, s3spec.Secretkey, spec.Billingcode)
	r.JSON(200, s3spec)
}

func delete(params martini.Params, r render.Render) {

	bucketname := params["name"]
	octs3.Deletebucket(bucketname)
	octiam.Detachuserpolicy(bucketname)
	octiam.Deleteaccesskey(bucketname)
	octiam.Deleteuser(bucketname)
	octdb.Delete(bucketname)

}

func runtest(billingcode string) {
	basename := createname()
	s3user := octiam.Createuser(basename)
	_ = octs3.Createbucket(s3user.Username, "basic")
	octs3.Tagbucket(basename, "billingcode", billingcode)
	time.Sleep(5000 * time.Millisecond)
	octs3.Addbucketpolicy(s3user.Username, s3user)
	simpleuserpolicy := octiam.Createuserpolicy(s3user.Username, s3user.Username)
	octiam.Attachuserpolicy(s3user.Username, simpleuserpolicy)

}
