package user

type Filter struct {
	Limit                     int32
	Page                      int32
	Keyword                   string
	ProjectID                 string
	ServiceDataKeyResourceIds []string
}
