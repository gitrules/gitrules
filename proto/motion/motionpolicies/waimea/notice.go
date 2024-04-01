package waimea

import (
	"fmt"

	"github.com/gitrules/gitrules/materials"
)

var Welcome = fmt.Sprintf(
	`

This project is managed by [GitRules](%s), a decentralized governance system for collaborative git projects.
To participate in governance, __install the [GitRules desktop app](%s)__.
	`,
	materials.GitRulesWebsiteURL,
	materials.GitRulesDesktopAppInstall,
)
