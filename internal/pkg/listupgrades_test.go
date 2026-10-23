package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The names of the upgrade folders sort v0.100.0 before v0.24.0, so the
// upgrades must be ordered by their version.
func TestListUpgradesIsOrderedByVersion(t *testing.T) {
	upgrades, err := ListUpgrades("v0.0.0")
	require.NoError(t, err)
	require.NotEmpty(t, upgrades)

	for i := 1; i < len(upgrades); i++ {
		previous, err := parseVersion(upgrades[i-1])
		require.NoError(t, err)
		current, err := parseVersion(upgrades[i])
		require.NoError(t, err)
		assert.True(t, isLowerVersion(previous, current), "%v must come after %v", upgrades[i-1], upgrades[i])
	}
}

func TestListUpgradesAfterAVersion(t *testing.T) {
	upgrades, err := ListUpgrades("v0.96.0")
	require.NoError(t, err)
	assert.Equal(t, []string{"v0.100.0"}, upgrades)
}

func TestLatestUpgradeVersion(t *testing.T) {
	latest, err := LatestUpgradeVersion()
	require.NoError(t, err)
	assert.Equal(t, "v0.100.0", latest)
}
