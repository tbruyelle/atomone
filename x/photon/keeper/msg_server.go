package keeper

import (
	"context"

	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	govtypes "github.com/atomone-hub/atomone/x/gov/types"
	"github.com/atomone-hub/atomone/x/photon/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// MintPhoton implements the MsgServer.MintPhoton method.
// TODO add logs & events
func (k msgServer) MintPhoton(goCtx context.Context, msg *types.MsgMintPhoton) (*types.MsgMintPhotonResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)
	if params.MintDisabled {
		return nil, types.ErrMintDisabled
	}

	// Ensure burned amount denom is bond denom (uatone)
	bondDenom := k.stakingKeeper.BondDenom(ctx)
	if msg.Amount.Denom != bondDenom {
		return nil, types.ErrBurnInvalidDenom
	}
	// Compute photons to mint
	var (
		atoneSupply    = k.bankKeeper.GetSupply(ctx, bondDenom).Amount.ToLegacyDec()
		photonSupply   = k.bankKeeper.GetSupply(ctx, "uphoton").Amount.ToLegacyDec()
		conversionRate = k.conversionRate(ctx, atoneSupply, photonSupply)
		atoneToBurn    = msg.Amount
		photonToMint   = atoneToBurn.Amount.ToLegacyDec().Mul(conversionRate)
	)
	// If no photon to mint, do not burn atoneToBurn, returns an error
	if photonToMint.IsZero() {
		return nil, types.ErrNoMintablePhotons
	}
	// If photonToMint + photonSupply > photonMaxSupply, returns an error
	if photonSupply.Add(photonToMint).GT(sdk.NewDec(PhotonMaxSupply)) {
		return nil, types.ErrNotEnoughPhotons
	}

	// Burn/Mint phase:
	// 1) move ATONEs from msg signer address to this module address
	// 2) burn ATONEs from this module address
	// 3) mint PHOTONs into this module address
	// 4) move PHOTONs from this module address to msg signer address
	var (
		coinsToBurn = sdk.NewCoins(atoneToBurn)
		coinsToMint = sdk.NewCoins(sdk.NewCoin("uphoton", photonToMint.RoundInt()))
	)
	// 1) Send atone to photon module for burn
	to, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return nil, err
	}
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, to, types.ModuleName, coinsToBurn); err != nil {
		return nil, err
	}
	// 2) Burn atone
	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, coinsToBurn); err != nil {
		return nil, err
	}

	// 3) Mint photons
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coinsToMint); err != nil {
		return nil, err
	}
	// 4) Send minted photon to account
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, to, coinsToMint); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintPhoton,
			sdk.NewAttribute(types.AttributeKeyBurned, coinsToBurn.String()),
			sdk.NewAttribute(types.AttributeKeyMinted, coinsToMint.String()),
		),
	})

	return &types.MsgMintPhotonResponse{
		Minted:         coinsToMint[0],
		ConversionRate: conversionRate.String(),
	}, nil
}

// UpdateParams implements the MsgServer.UpdateParams method.
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}