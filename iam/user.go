package iam

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	structs "s3-aws-api/structs"
	utils "s3-aws-api/utils"
)

func Createuser(username string) structs.S3user {

	svc := iam.New(session.New(&aws.Config{
		Region: aws.String(utils.Region),
	}))

	params := &iam.CreateUserInput{
		UserName: aws.String(username),
	}
	resp, err := svc.CreateUser(params)

	if err != nil {
		fmt.Println(err.Error())

	}

	arn := *resp.User.Arn

	paramskey := &iam.CreateAccessKeyInput{
		UserName: aws.String(username),
	}
	respkey, err := svc.CreateAccessKey(paramskey)

	if err != nil {
		fmt.Println(err.Error())
	}

	accesskey := *respkey.AccessKey.AccessKeyId
	secretkey := *respkey.AccessKey.SecretAccessKey
	var s3user structs.S3user
	s3user.Username = username
	s3user.Arn = arn
	s3user.Accesskey = accesskey
	s3user.Secretkey = secretkey
	return s3user

}

func Createuserpolicy(username string, bucketname string) structs.SimpleUserPolicy {

	var userpolicy structs.UserPolicy
	userpolicy.Version = "2012-10-17"
	var statements []structs.UserPolicyStatement
	var statement structs.UserPolicyStatement
	statement.Effect = "Allow"
	var resources []string
	resources = append(resources, "arn:aws:s3:::"+bucketname+"/*")
	resources = append(resources, "arn:aws:s3:::"+bucketname)
	statement.Resource = resources
	var actions []string
	actions = append(actions, "s3:*")
	statement.Action = actions
	statements = append(statements, statement)
	userpolicy.Statement = statements
	str, err := json.Marshal(userpolicy)
	if err != nil {
		fmt.Println("Error preparing request")
	}
	jsonStr := (string(str))

	svc := iam.New(session.New(&aws.Config{
		Region: aws.String(utils.Region),
	}))

	params := &iam.CreatePolicyInput{
		PolicyDocument: aws.String(jsonStr),
		PolicyName:     aws.String(username + "policy"),
	}
	resp, err := svc.CreatePolicy(params)

	if err != nil {
		fmt.Println(err.Error())
	}

	policyarn := *resp.Policy.Arn
	policyname := *resp.Policy.PolicyName
	var simpleuserpolicy structs.SimpleUserPolicy
	simpleuserpolicy.PolicyName = policyname
	simpleuserpolicy.Arn = policyarn
	return simpleuserpolicy
}

func Attachuserpolicy(username string, simpleuserpolicy structs.SimpleUserPolicy) {
	svc := iam.New(session.New(&aws.Config{
		Region: aws.String(utils.Region),
	}))

	params := &iam.AttachUserPolicyInput{
		PolicyArn: aws.String(simpleuserpolicy.Arn),
		UserName:  aws.String(username),
	}
	_, err := svc.AttachUserPolicy(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func Deleteuser(bucketname string) {

	svc := iam.New(session.New(&aws.Config{
		Region: aws.String(utils.Region),
	}))

	params := &iam.DeleteUserInput{
		UserName: aws.String(bucketname), // Required
	}
	_, err := svc.DeleteUser(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}


}

func Detachuserpolicy(bucketname string) {

	svc := iam.New(session.New(&aws.Config{
		Region: aws.String(utils.Region),
	}))

	policyarn := getpolicyarn(bucketname)
	deparams := &iam.DetachUserPolicyInput{
		PolicyArn: aws.String(policyarn),  // Required
		UserName:  aws.String(bucketname), // Required
	}

	_, err := svc.DetachUserPolicy(deparams)

	if err != nil {
		fmt.Println(err.Error())
		return
	}


	deleteuserpolicy(policyarn)

}

func getpolicyarn(bucketname string) string {

	svc := iam.New(session.New(&aws.Config{
		Region: aws.String(utils.Region),
	}))
	params := &iam.ListAttachedUserPoliciesInput{
		UserName: aws.String(bucketname), // Required
	}
	resp, err := svc.ListAttachedUserPolicies(params)

	if err != nil {
		fmt.Println(err.Error())
	}

	policyarn := *resp.AttachedPolicies[0].PolicyArn
	return policyarn
}

func getaccesskeyid(bucketname string) string {

	svc := iam.New(session.New(&aws.Config{
		Region: aws.String(utils.Region),
	}))

	params := &iam.ListAccessKeysInput{
		UserName: aws.String(bucketname),
	}
	resp, err := svc.ListAccessKeys(params)

	if err != nil {
		fmt.Println(err.Error())
	}

	accesskeyid := *resp.AccessKeyMetadata[0].AccessKeyId
	return accesskeyid

}

func Deleteaccesskey(bucketname string) {
	accesskeyid := getaccesskeyid(bucketname)

	svc := iam.New(session.New(&aws.Config{
		Region: aws.String(utils.Region),
	}))

	params := &iam.DeleteAccessKeyInput{
		AccessKeyId: aws.String(accesskeyid), // Required
		UserName:    aws.String(bucketname),
	}
	_, err := svc.DeleteAccessKey(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}


}

func deleteuserpolicy(policyarn string) {

	svc := iam.New(session.New(&aws.Config{
		Region: aws.String(utils.Region),
	}))

	params := &iam.DeletePolicyInput{
		PolicyArn: aws.String(policyarn), // Required
	}
	_, err := svc.DeletePolicy(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}


}
