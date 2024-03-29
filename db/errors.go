package db

import "errors"

var (
	ErrNoRow              = errors.New("data not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrBalanceNotFound    = errors.New("balance not found")
	ErrinsuficientBalance = errors.New("insuficient balance")
)
