package v1_test

import (
	"testing"

	"github.com/prysmaticlabs/prysm/beacon-chain/state"
	testtmpl "github.com/prysmaticlabs/prysm/beacon-chain/state/testing"
	v1 "github.com/prysmaticlabs/prysm/beacon-chain/state/v1"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

func TestBeaconState_ValidatorAtIndexReadOnly_HandlesNilSlice(t *testing.T) {
	testtmpl.VerifyBeaconState_ValidatorAtIndexReadOnly_HandlesNilSlice(t, func() (state.BeaconState, error) {
		return v1.InitializeFromProtoUnsafe(&ethpb.BeaconState{
			Validators: nil,
		})
	})
}
