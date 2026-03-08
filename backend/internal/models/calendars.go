package models

import "github.com/google/uuid"

type Permission string

const (
	PermissionRead  Permission = "read"
	PermissionWrite Permission = "write"
)

type Visibility string

const (
	VisibilityPrivate Visibility = "private"
	VisibilityShared  Visibility = "shared"
	VisibilityPubic   Visibility = "public"
)

type Calendar struct {
	Id         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	OwnerId    uuid.UUID  `json:"ownerId"`
	Visibility Visibility `json:"visibility"`
}

type AddCalendar struct {
	Name string `json:"name" validate:"required,max=100"`
}

type EditCalendar struct {
	Name *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
}

type ShareCalendar struct {
	UserId     uuid.UUID  `json:"userId" validate:"required,uuid"`
	Permission Permission `json:"permissions" validate:"required,oneof=read write"`
}

type BatchShareCalendar struct {
	Items []ShareCalendar `json:"items" validate:"required,min=1,dive"`
}

type BatchRemoveUserShare struct {
	Ids []uuid.UUID `json:"ids" validate:"required,min=1,dive,required,uuid"`
}
