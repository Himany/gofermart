package utils

func SetStringIfUnset(envSet map[string]bool, envKey string, cfgValue *string, flagValue string) {
	if envSet[envKey] {
		return
	}
	*cfgValue = flagValue
}
