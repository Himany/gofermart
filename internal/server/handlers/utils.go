package handlers

func validateLuhn(num string) bool {
	if len(num) == 0 {
		return false
	}
	sum := 0
	parity := len(num) % 2
	for index, value := range num {
		if value < '0' || value > '9' {
			return false
		}
		d := int(value - '0')
		if index%2 == parity {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
	}
	return sum%10 == 0
}
