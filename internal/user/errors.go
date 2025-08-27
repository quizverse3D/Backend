package user

import "errors"

var (
	ErrUserNotFound                 = errors.New("user not found")
	ErrUserParamsNotFound           = errors.New("user client params not found")
	ErrUserParamsInvalidLangCode    = errors.New("lang_code is invalid")
	ErrUserParamsInvalidSoundVolume = errors.New("sound_volume is invalid")
)
