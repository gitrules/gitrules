package etc

import "github.com/gitrules/gitrules/proto"

var (
	EtcNS      = proto.RootNS.Append("etc")
	SettingsNS = EtcNS.Append("settings.json")
)

type Settings struct {
}

var DefaultSettings = Settings{}
