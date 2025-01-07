package models

import (
	"time"
)

type AdminLog struct {
	Id           string
	Action       adminActions
	ResourceType resourceTypes
	ResourceName string
	Username     string
	CreatedAt    time.Time
}

type AdminLogResponse struct {
	Action       adminActions
	Username     string
	ResourceId   string
	ResourceName string
	ResourceType resourceTypes
	CreatedAt    time.Time
}

type resourceTypes string

const (
	User         resourceTypes = "USER"
	DbConnection resourceTypes = "DB_CONNECTION"
)

type adminActions string

const (
	CreateAdminAction            adminActions = "CREATE"
	UpdateSettingsAdminAction    adminActions = "UPDATE_SETTINGS"
	UpdateCredentialsAdminAction adminActions = "UPDATE_CREDENTIALS"
	DeleteAdminAction            adminActions = "DELETE"
)
