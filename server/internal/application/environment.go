package application

import "strings"

type Environment string

const (
	EnvironmentDevelopment Environment = "development"
	EnvironmentProduction  Environment = "production"
	EnvironmentUnknown     Environment = "unknown"
)

func NewEnvironment(value string) Environment {
	switch strings.ToLower(value) {
	case "development":
		return EnvironmentDevelopment
	case "production":
		return EnvironmentProduction
	default:
		return EnvironmentUnknown
	}
}

func (e Environment) String() string {
	return string(e)
}

func (e Environment) IsDevelopment() bool {
	return e == EnvironmentDevelopment
}

func (e Environment) IsProduction() bool {
	return e == EnvironmentProduction
}
