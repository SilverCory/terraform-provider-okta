package sdk

import "encoding/json"

type UserSchemaAttribute struct {
	Description       string                           `json:"description,omitempty"`
	Enum              []interface{}                    `json:"enum,omitempty"`
	ExternalName      string                           `json:"externalName,omitempty"`
	ExternalNamespace string                           `json:"externalNamespace,omitempty"`
	Items             *UserSchemaAttributeItems        `json:"items,omitempty"`
	Master            *UserSchemaAttributeMaster       `json:"master,omitempty"`
	MaxLength         int64                            `json:"-"`
	MaxLengthPtr      *int64                           `json:"maxLength,omitempty"`
	MinLength         int64                            `json:"-"`
	MinLengthPtr      *int64                           `json:"minLength,omitempty"`
	Mutability        string                           `json:"mutability,omitempty"`
	OneOf             []*UserSchemaAttributeEnum       `json:"oneOf,omitempty"`
	Pattern           *string                          `json:"pattern,omitempty"`
	Permissions       []*UserSchemaAttributePermission `json:"permissions,omitempty"`
	Required          *bool                            `json:"required,omitempty"`
	Scope             string                           `json:"scope,omitempty"`
	Title             string                           `json:"title,omitempty"`
	Type              string                           `json:"type,omitempty"`
	Union             string                           `json:"union,omitempty"`
	Unique            string                           `json:"unique,omitempty"`
}

func (a *UserSchemaAttribute) MarshalJSON() ([]byte, error) {
	type Alias UserSchemaAttribute
	type local struct {
		*Alias
	}
	result := local{Alias: (*Alias)(a)}
	if a.MaxLength != 0 {
		result.MaxLengthPtr = Int64Ptr(a.MaxLength)
	}
	if a.MinLength != 0 {
		result.MinLengthPtr = Int64Ptr(a.MinLength)
	}
	return json.Marshal(&result)
}

func (a *UserSchemaAttribute) UnmarshalJSON(data []byte) error {
	type Alias UserSchemaAttribute

	result := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	if result.MaxLengthPtr != nil {
		a.MaxLength = *result.MaxLengthPtr
		a.MaxLengthPtr = result.MaxLengthPtr
	}
	if result.MinLengthPtr != nil {
		a.MinLength = *result.MinLengthPtr
		a.MinLengthPtr = result.MinLengthPtr
	}
	return nil
}
