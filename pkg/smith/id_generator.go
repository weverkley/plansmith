package smith

import (
	"fmt"
)

func GenerateID(prefix string, id int) string {
	return fmt.Sprintf("[%s_%d]", prefix, id)
}
