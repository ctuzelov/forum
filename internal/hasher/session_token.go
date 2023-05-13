package hasher

import "github.com/gofrs/uuid"

func GenerateToken() (*string, error) {
	token, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	tokenstring := token.String()
	return &tokenstring, nil
}
