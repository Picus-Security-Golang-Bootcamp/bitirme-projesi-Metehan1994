// Code generated by go-swagger; DO NOT EDIT.

package api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// CategoryWithoutRequiredName category without required name
//
// swagger:model CategoryWithoutRequiredName
type CategoryWithoutRequiredName struct {

	// description
	Description string `json:"description,omitempty"`

	// id
	ID int64 `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`
}

// Validate validates this category without required name
func (m *CategoryWithoutRequiredName) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this category without required name based on context it is used
func (m *CategoryWithoutRequiredName) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *CategoryWithoutRequiredName) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CategoryWithoutRequiredName) UnmarshalBinary(b []byte) error {
	var res CategoryWithoutRequiredName
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
