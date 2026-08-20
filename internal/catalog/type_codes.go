package catalog

import "fmt"

func TypeName(code uint16) string {
	for _, m := range metadata {
		if m.Code == code {
			return m.Name
		}
	}
	return fmt.Sprintf("TYPE%d", code)
}
