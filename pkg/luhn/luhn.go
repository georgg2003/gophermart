package luhn

func ValidLuhn(s string) bool {
	var sum int
	double := false

	if len(s) == 0 {
		return false
	}

	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		if c < '0' || c > '9' {
			return false
		}

		d := int(c - '0')

		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}

		sum += d
		double = !double
	}

	return sum%10 == 0
}
