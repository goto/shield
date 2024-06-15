package user

type Filter struct {
	Limit                     int32
	Page                      int32
	Keyword                   string
	Project                   string
	ServiceDataKeyResourceIds []string
}
