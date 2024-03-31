package pmp_0

import (
	"fmt"

	"github.com/gitrules/gitrules/materials"
)

var Welcome = fmt.Sprintf(
	`

This project is managed by [Gov4Git](%s), a decentralized governance system for collaborative git projects.
To participate in governance, __install the [Gov4Git desktop app](%s)__.
	`,
	materials.Gov4GitWebsiteURL,
	materials.Gov4GitDesktopAppInstall,
)
