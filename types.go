package myvmware

import "net/http"

type Client struct {
	HttpClient *http.Client
}

// TODO: for now we keep the types as interface{}, because we don't really know yet

type Products interface{}

type AccountInfo interface{}

// type AccountInfo struct {
// 	UserType    string   `json:"userType"`
// 	AccountList []string `json:"accntList"`
// }

type CurrentUser interface{}

// type CurrentUser struct {
// 	FirstName string `json:"firstname"`
// 	LastName  string `json:"lastname"`
// }
