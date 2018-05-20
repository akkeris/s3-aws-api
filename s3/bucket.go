package s3

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	structs "s3-aws-api/structs"
	utils "s3-aws-api/utils"
	"strings"
)

func Createbucket(bucketname string, plan string) string {
	svc := s3.New(session.New(&aws.Config{
		Region: aws.String(utils.Region),
	}))

	params := &s3.CreateBucketInput{
		Bucket: aws.String(bucketname),
	}
	resp, err := svc.CreateBucket(params)

	if err != nil {
		fmt.Println(err.Error())
	}

	if plan == "versioned" {
		setVersioning(bucketname)
		setBucketLifecycle(bucketname)
	}
	bucketlocation := *resp.Location
	bucketlocation = strings.Replace(bucketlocation, "http://", "", -1)
	bucketlocation = strings.Replace(bucketlocation, "/", "", -1)
	return bucketlocation
}

func Addbucketpolicy(bucketname string, s3user structs.S3user) {
	userarn := s3user.Arn

	var mypolicy structs.Bucketpolicy
	mypolicy.Version = "2012-10-17"
	mypolicy.ID = "Policy47474747"
	var statements []structs.Statement
	var statement structs.Statement
	statement.Sid = "Stmt47474747"
	statement.Effect = "Allow"
	statement.Principal.AWS = userarn
	statement.Action = "s3:*"
	statement.Resource = "arn:aws:s3:::" + bucketname + "/*"
	statements = append(statements, statement)
	mypolicy.Statement = statements
	str, err := json.Marshal(mypolicy)
	if err != nil {
		fmt.Println("Error preparing request")
	}
	jsonStr := (string(str))

	svc := s3.New(session.New(&aws.Config{
		Region: aws.String(utils.Region),
	}))

	params := &s3.PutBucketPolicyInput{
		Bucket: aws.String(bucketname),
		Policy: aws.String(jsonStr),
	}
         
	_, err= svc.PutBucketPolicy(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}


}

func gettags(bucketname string) []*s3.Tag {

	svc := s3.New(session.New(&aws.Config{
		Region: aws.String(utils.Region),
	}))

	params := &s3.GetBucketTaggingInput{
		Bucket: aws.String(bucketname), // Required
	}
	resp, err := svc.GetBucketTagging(params)

	if err != nil {
		fmt.Println(err.Error())
	}

	return resp.TagSet

}

func Tagbucket(bucketname string, name string, value string) {
	existingtags := gettags(bucketname)
	var new s3.Tag
	var newkey string
	newkey = name
	var newvalue string
	newvalue = value
	new.Key = &newkey
	new.Value = &newvalue
	existingtags = append(existingtags, &new)

	svc := s3.New(session.New(&aws.Config{
		Region: aws.String(utils.Region),
	}))

	params := &s3.PutBucketTaggingInput{
		Bucket: aws.String(bucketname),
		Tagging: &s3.Tagging{
			TagSet: existingtags,
		},
	}

	_, err := svc.PutBucketTagging(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}


}

func Deletebucket(bucketname string) {

	svc := s3.New(session.New(&aws.Config{
		Region: aws.String(utils.Region),
	}))

	params := &s3.DeleteBucketInput{
		Bucket: aws.String(bucketname), // Required
	}
	_, err := svc.DeleteBucket(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}


}

func setVersioning(bucketname string) {

	svc := s3.New(session.New(&aws.Config{
		Region: aws.String(utils.Region),
	}))

	params := &s3.PutBucketVersioningInput{
		Bucket: aws.String(bucketname),
		VersioningConfiguration: &s3.VersioningConfiguration{
			Status: aws.String("Enabled"),
		},
	}
	_, err := svc.PutBucketVersioning(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}


}

func setBucketLifecycle(bucketname string) {

	svc := s3.New(session.New(&aws.Config{
		Region: aws.String(utils.Region),
	}))

	params := &s3.PutBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucketname),
		LifecycleConfiguration: &s3.BucketLifecycleConfiguration{
			Rules: []*s3.LifecycleRule{
				{
					Prefix: aws.String(""),
					Status: aws.String("Enabled"),
					ID:     aws.String("versioned"),
					NoncurrentVersionExpiration: &s3.NoncurrentVersionExpiration{
						NoncurrentDays: aws.Int64(180),
					},
					NoncurrentVersionTransitions: []*s3.NoncurrentVersionTransition{
						{
							NoncurrentDays: aws.Int64(30),
							StorageClass:   aws.String("STANDARD_IA"),
						},
					},
					Transitions: []*s3.Transition{
						{
							Days:         aws.Int64(30),
							StorageClass: aws.String("STANDARD_IA"),
						},
					},
				},
			},
		},
	}
	_, err := svc.PutBucketLifecycleConfiguration(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
