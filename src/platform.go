package main

type Platform string

const (
	TELEGRAM = "tg"
	DISCORD  = "ds"
)

func (p Platform) Prefix() string {
	switch p {
	case TELEGRAM:
		return "tg-"
	case DISCORD:
		return "ds-"
	}
	return "UKNOWN"
}
