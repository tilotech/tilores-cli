package pkg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tilotech/tilores-cli/templates"
)

func ListUpgrades(version string) ([]string, error) {
	ver, err := parseVersion(version)
	if err != nil {
		return nil, err
	}
	upgradeFolders, err := templates.Upgrades.ReadDir("upgrades")
	if err != nil {
		return nil, err
	}
	upgrades := make([]string, 0, len(upgradeFolders))
	for _, upgradeFolder := range upgradeFolders {
		if !upgradeFolder.IsDir() {
			continue
		}
		upgrade := upgradeFolder.Name()
		upgradeVer, err := parseVersion(upgrade)
		if err != nil {
			return nil, err
		}
		if isLowerVersion(ver, upgradeVer) {
			upgrades = append(upgrades, upgrade)
		}
	}
	return upgrades, nil
}

func LatestUpgradeVersion() (string, error) {
	upgrades, err := ListUpgrades("v0.0.0")
	if err != nil {
		return "", err
	}
	return upgrades[len(upgrades)-1], nil
}

type parsedVersion struct {
	major int
	minor int
	patch int
}

func parseVersion(version string) (*parsedVersion, error) {
	invalidVersionErr := fmt.Errorf("invalid version %v", version)
	if version[0] != "v"[0] {
		return nil, invalidVersionErr
	}
	parts := strings.Split(version[1:], ".")
	if len(parts) != 3 {
		return nil, invalidVersionErr
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, invalidVersionErr
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, invalidVersionErr
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, invalidVersionErr
	}
	return &parsedVersion{
		major: major,
		minor: minor,
		patch: patch,
	}, nil
}

func isLowerVersion(c, u *parsedVersion) bool {
	if c.major < u.major {
		return true
	}
	if c.major > u.major {
		return false
	}
	if c.minor < u.minor {
		return true
	}
	if c.minor > u.minor {
		return false
	}
	return c.patch < u.patch
}
