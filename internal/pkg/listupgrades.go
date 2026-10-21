package pkg

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/tilotech/tilores-cli/templates"
)

// ListUpgrades provides a list of all unapplied upgrade versions, oldest first.
//
// The folders are sorted by their version, because their names sort v0.100.0
// before v0.24.0.
func ListUpgrades(version string) ([]string, error) {
	ver, err := parseVersion(version)
	if err != nil {
		return nil, err
	}
	upgradeFolders, err := templates.Upgrades.ReadDir("upgrades")
	if err != nil {
		return nil, err
	}
	type upgrade struct {
		name    string
		version *parsedVersion
	}
	upgrades := make([]upgrade, 0, len(upgradeFolders))
	for _, upgradeFolder := range upgradeFolders {
		if !upgradeFolder.IsDir() {
			continue
		}
		upgradeVer, err := parseVersion(upgradeFolder.Name())
		if err != nil {
			return nil, err
		}
		if isLowerVersion(ver, upgradeVer) {
			upgrades = append(upgrades, upgrade{name: upgradeFolder.Name(), version: upgradeVer})
		}
	}
	sort.Slice(upgrades, func(i, j int) bool {
		return isLowerVersion(upgrades[i].version, upgrades[j].version)
	})
	names := make([]string, len(upgrades))
	for i, u := range upgrades {
		names[i] = u.name
	}
	return names, nil
}

// LatestUpgradeVersion returns the latest upgrade version.
//
// This version can be used for a fresh init.
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
