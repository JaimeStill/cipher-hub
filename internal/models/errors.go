package models

import "errors"

var (
	ErrInvalidID                = errors.New("invalid ID")
	ErrInvalidName              = errors.New("invalid name")
	ErrInvalidServiceID         = errors.New("invalid service ID")
	ErrInvalidParticipantStatus = errors.New("invalid participant status")
	ErrInvalidAlgorithm         = errors.New("invalid algorithm")
	ErrInvalidKeyStatus         = errors.New("invalid key status")
	ErrInvalidVersion           = errors.New("invalid version")
)
