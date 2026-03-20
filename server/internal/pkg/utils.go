package pkg

func GetOrDefault4String(src string, defaultVal string) string {
	if src != "" {
		return src
	}
	return defaultVal
}
