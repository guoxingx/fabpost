package main

const (
	SecretObjectType = "secret"
	OwnerObejctType  = "owner"
)

type Secret struct {
	ObjectType string        `json:"docType"`
	Id         string        `json:"id"`
	Content    string        `json:"content"`
	Owner      OwnerRelation `"json:owner"`
}

type Owner struct {
	ObjectType string `json:"docType"`
	Id         string `json:"id"`
	Username   string `json:"username"`
}

type OwnerRelation struct {
	Id       string `"json:id"`
	Username string `"json:username"`
}
