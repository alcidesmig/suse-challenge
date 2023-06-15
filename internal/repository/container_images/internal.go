package container_images

import "regexp"

func findImageValues(input string) []string {
	regexPattern := `image:\s*(.*?)\n`
	regex := regexp.MustCompile(regexPattern)
	matches := regex.FindAllStringSubmatch(input, -1)

	var imageValues []string
	for _, match := range matches {
		if len(match) >= 2 {
			imageValues = append(imageValues, match[1])
		}
	}

	return imageValues
}
