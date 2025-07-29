package api

//go:generate go tool oapi-codegen -config openapi_config.yaml openapi.yaml

func StringPtr(s string) *string {
	return &s
}
