package namespacesprivateendpointconnections

import "fmt"

const defaultApiVersion = "2018-01-01-preview"

func userAgent() string {
	return fmt.Sprintf("pandora/namespacesprivateendpointconnections/%s", defaultApiVersion)
}
