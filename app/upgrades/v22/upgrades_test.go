package v22_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	v22 "github.com/cosmos/gaia/v21/app/upgrades/v22"
	testutil "github.com/cosmos/interchain-security/v6/testutil/keeper"
	providertypes "github.com/cosmos/interchain-security/v6/x/ccv/provider/types"
)

func TestSetDefaultConsumerInfractionParams(t *testing.T) {
	t.Helper()
	inMemParams := testutil.NewInMemKeeperParams(t)
	pk, ctx, ctrl, _ := testutil.GetProviderKeeperAndCtx(t, inMemParams)
	defer ctrl.Finish()

	// Add consumer chains
	initConsumerID := pk.FetchAndIncrementConsumerId(ctx)
	pk.SetConsumerChainId(ctx, initConsumerID, "init-1")
	pk.SetConsumerPhase(ctx, initConsumerID, providertypes.CONSUMER_PHASE_INITIALIZED)
	launchedConsumerID := pk.FetchAndIncrementConsumerId(ctx)
	pk.SetConsumerChainId(ctx, launchedConsumerID, "launched-1")
	pk.SetConsumerPhase(ctx, launchedConsumerID, providertypes.CONSUMER_PHASE_LAUNCHED)
	stoppedConsumerID := pk.FetchAndIncrementConsumerId(ctx)
	pk.SetConsumerChainId(ctx, stoppedConsumerID, "stopped-1")
	pk.SetConsumerPhase(ctx, stoppedConsumerID, providertypes.CONSUMER_PHASE_STOPPED)

	activeConsumerIds := pk.GetAllActiveConsumerIds(ctx)
	require.Equal(t, 2, len(activeConsumerIds))

	for _, consumerId := range activeConsumerIds {
		_, err := pk.GetInfractionParameters(ctx, consumerId)
		require.Error(t, err)
	}

	err := v22.SetConsumerInfractionParams(ctx, pk)
	require.NoError(t, err)

	defaultInfractionParams := v22.DefaultInfractionParams()
	for _, consumerId := range activeConsumerIds {
		infractionParams, err := pk.GetInfractionParameters(ctx, consumerId)
		require.NoError(t, err)
		require.Equal(t, defaultInfractionParams, infractionParams)
	}
}
