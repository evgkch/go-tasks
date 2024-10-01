//go:build !solution

package retryupdate

import (
	"errors"

	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	var authErr *kvapi.AuthError
	var conflictErr *kvapi.ConflictError

	for {
		var oldValue *string
		var oldVersion uuid.UUID

		getResp, err := c.Get(&kvapi.GetRequest{Key: key})

		switch {
		case err == nil:
			oldValue = &getResp.Value
			oldVersion = getResp.Version

		case errors.Is(err, kvapi.ErrKeyNotFound):
			oldValue = nil
			oldVersion = uuid.UUID{}

		case errors.As(err, &authErr):
			return err

		default:
			continue
		}

		newValue, err := updateFn(oldValue)
		if err != nil {
			return err
		}

		newVerion := uuid.Must(uuid.NewV4())

		for {
			_, err := c.Set(&kvapi.SetRequest{
				Key:        key,
				Value:      newValue,
				OldVersion: oldVersion,
				NewVersion: newVerion,
			})

			switch {
			case err == nil:
				return nil

			case errors.Is(err, kvapi.ErrKeyNotFound):
				oldValue = nil
				oldVersion = uuid.UUID{}

				newValue, err = updateFn(oldValue)
				if err != nil {
					return err
				}
				continue

			case errors.As(err, &authErr):
				return err

			case errors.As(err, &conflictErr):
				if conflictErr.ExpectedVersion == newVerion {
					return nil
				}
				goto outerLoop

			default:
				continue
			}
		}
	outerLoop:
	}
}
