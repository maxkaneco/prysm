package blockchain

import "github.com/pkg/errors"

var (
	// errNilJustifiedInStore is returned when a nil justified checkpt is returned from store.
	errNilJustifiedInStore = errors.New("nil justified checkpoint returned from store")
	// errNilBestJustifiedInStore is returned when a nil justified checkpt is returned from store.
	errNilBestJustifiedInStore = errors.New("nil best justified checkpoint returned from store")
	// errNilFinalizedInStore is returned when a nil finalized checkpt is returned from store.
	errNilFinalizedInStore = errors.New("nil finalized checkpoint returned from store")
)
