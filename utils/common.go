package utils

import "github.com/google/uuid"

func UUID() string {
	for {
		uuid, err := uuid.NewUUID()
		if err != nil {
			continue
		}
		return uuid.String()
	}
}

func ToUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
