package model

import "github.com/google/uuid"

//CreateFormInput model
type CreateFormInput struct {
	Name        *string    `json:"name" validate:"required"`
	Description *string    `json:"description" validate:"required"`
	Type        *string    `json:"type" validate:"required"`
	Structure   *string    `json:"structure" validate:"required"`
	Reference   *uuid.UUID `json:"reference,omitempty"`
}

//ValidateInput model
type ValidateInput struct {
	ID             *uuid.UUID              `json:"id,omitempty"`
	Version        int                     `json:"version" validate:"required"`
	TypeForm       string                  `json:"type" validate:"required"`
	InputedForm    map[string]interface{}  `json:"values" validate:"required"`
	AdditionalData *map[string]interface{} `json:"additional_data,omitempty"`
}

//GetLatestInput model
type GetLatestInput struct {
	TypeForm string    `json:"type" httpurl:"type_form" validate:"required"`
	ID       uuid.UUID `json:"id" httpurl:"id"`
}

//PublishInput model
type PublishInput struct {
	ID uuid.UUID `json:"id" httpurl:"id" validate:"required"`
}
