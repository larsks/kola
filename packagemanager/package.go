package packagemanager

import (
	"fmt"

	operators "github.com/operator-framework/operator-lifecycle-manager/pkg/package-server/apis/operators/v1"
)

type (
	Package struct {
		operators.PackageManifest
	}
)

func (pkg *Package) GetDefaultKeywords() ([]string, error) {
	channel := pkg.GetDefaultChannel()
	return channel.CurrentCSVDesc.Keywords, nil
}

func (pkg *Package) GetChannelByName(name string) (*operators.PackageChannel, error) {
	for _, channel := range pkg.Status.Channels {
		if channel.Name == name {
			return &channel, nil
		}
	}
	return nil, fmt.Errorf("channel %s not found", name)
}

func (pkg *Package) GetDefaultChannelName() string {
	return pkg.Status.DefaultChannel
}

func (pkg *Package) GetDefaultChannel() *operators.PackageChannel {
	channelName := pkg.GetDefaultChannelName()
	channel, _ := pkg.GetChannelByName(channelName)
	return channel
}

func (pkg *Package) GetDefaultInstallModes() []string {
	var installModes []string
	channel := pkg.GetDefaultChannel()
	for _, installMode := range channel.CurrentCSVDesc.InstallModes {
		if installMode.Supported {
			installModes = append(installModes, string(installMode.Type))
		}
	}

	return installModes
}

func (pkg *Package) GetChannelNames() []string {
	var channelNames []string

	for _, channel := range pkg.Status.Channels {
		channelNames = append(channelNames, channel.Name)
	}

	return channelNames
}

func (pkg *Package) GetChannels() []operators.PackageChannel {
	return pkg.Status.Channels
}

func (pkg *Package) GetDefaultDescription() string {
	channel := pkg.GetDefaultChannel()
	return channel.CurrentCSVDesc.LongDescription
}

func (pkg *Package) SupportsInstallMode(mode string) bool {
	channel := pkg.GetDefaultChannel()
	for _, installMode := range channel.CurrentCSVDesc.InstallModes {
		if string(installMode.Type) == mode {
			return installMode.Supported
		}
	}
	return false
}
