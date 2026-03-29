package models

import "github.com/google/uuid"

type Folder struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type AddFolder struct {
	Name string `json:"name" validate:"required,max=100"`
}

type FolderEdit struct {
	Name *string `json:"name,omitempty" validate:"omitempty,max=100"`
}
