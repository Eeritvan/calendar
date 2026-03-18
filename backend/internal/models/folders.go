package models

type AddFolder struct {
	Name string `json:"name" validate:"required,max=100"`
}

type FolderEdit struct {
	Name *string `json:"name,omitempty" validate:"omitempty,max=100"`
}
