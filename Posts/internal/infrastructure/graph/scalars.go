package graph

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"time"
)

// MarshalUUID marshals a UUID to a string.
func MarshalUUID(v uuid.UUID) graphql.Marshaler {
	return graphql.MarshalString(v.String())
}

// UnmarshalUUID unmarshals a UUID from a string.
func UnmarshalUUID(v interface{}) (uuid.UUID, error) {
	id, err := uuid.Parse(v.(string))
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

// MarshalDateTime marshals a time.Time to a string.
func MarshalDateTime(v time.Time) graphql.Marshaler {
	return graphql.MarshalString(v.String())
}

// UnmarshalDateTime unmarshals a time.Time from a string.
func UnmarshalDateTime(v interface{}) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, v.(string))
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
