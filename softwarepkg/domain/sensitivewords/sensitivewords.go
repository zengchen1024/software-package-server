package sensitivewords

type SensitiveWords interface {
	CheckSensitiveWords(string) error
}
