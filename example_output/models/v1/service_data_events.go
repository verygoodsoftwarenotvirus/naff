package models

type ServiceDataEvent string

const (
	Create  ServiceDataEvent = "create"
	Update  ServiceDataEvent = "update"
	Archive ServiceDataEvent = "delete"
)
