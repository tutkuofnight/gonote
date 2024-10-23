// This file only includes useful functions that i have to go back to look

func HexStringToRGB(hexstr string) (int, int, int, error) {
	c, err := hex.DecodeString(hexstr) // expected something like: FF0000

	if err != nil {
		return 0, 0, 0, err
	}

	if len(c) != 3 {
		return 0, 0, 0, errors.New("invalid hex string")
	}

	return int(c[0]), int(c[1]), int(c[2]), nil
}
