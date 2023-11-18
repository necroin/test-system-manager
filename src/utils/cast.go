package utils

import "fmt"

func CastToString(value any) string {
	return fmt.Sprintf("%v", value)
}
