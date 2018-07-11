package structs

type Tagspec struct {
	Resource string `json:"resource"`
	Name     string `json:"name"`
	Value    string `json:"value"`
}

type Provisionspec struct {
	Plan        string `json:"plan"`
	Billingcode string `json:"billingcode"`
}

type SimpleUserPolicy struct {
	PolicyName string
	Arn        string
}

type UserPolicyStatement struct {
	Resource []string `json:"Resource"`
	Action   []string `json:"Action"`
	Effect   string   `json:"Effect"`
}

type UserPolicy struct {
	Statement []UserPolicyStatement `json:"Statement"`
	Version   string                `json:"Version"`
}

type S3user struct {
	Username  string
	Arn       string
	Accesskey string
	Secretkey string
}
type Statement struct {
	Sid       string `json:"Sid"`
	Effect    string `json:"Effect"`
	Principal struct {
		AWS string `json:"AWS"`
	} `json:"Principal"`
	Action   string `json:"Action"`
	Resource string `json:"Resource"`
}

type Bucketpolicy struct {
	Version   string      `json:"Version"`
	ID        string      `json:"Id"`
	Statement []Statement `json:"Statement"`
}

type S3spec struct {
	Location  string `json:"S3_LOCATION"`
        Bucket    string `json:"S3_BUCKET"`
	Accesskey string `json:"S3_ACCESS_KEY"`
	Secretkey string `json:"S3_SECRET_KEY"`
	Region string `json:"S3_REGION"`
}
