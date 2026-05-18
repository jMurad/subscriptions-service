package query

func ParseOptionalString(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}
