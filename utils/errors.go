package utils

import "errors"

func GetNoEmptyError(err error) error {
	if err == nil {
		return errors.New("")
	} else {
		return err
	}
}
