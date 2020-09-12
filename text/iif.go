/*
 * Copyright (c) 2020. Brickman Source.
 */

package text

func Iif(condition bool, retIf1, retIf2 string) string {
	if condition {
		return retIf1
	}
	return retIf2
}
