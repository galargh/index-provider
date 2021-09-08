package suppliers

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-car/v2/blockstore"
	"github.com/stretchr/testify/require"
)

func Test_drain(t *testing.T) {
	tests := []struct {
		name    string
		carPath string
	}{
		{
			"CARv1ReturnsExpectedCIDs",
			"../../testdata/sample-v1.car",
		},
		{
			"CARv2ReturnsExpectedCIDs",
			"../../testdata/sample-wrapped-v2.car",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cidIter, err := newCarCidIterator(tt.carPath)
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, cidIter.Close()) })

			// Open ReadOnly blockstore used to provide wanted case for testing
			want, err := blockstore.OpenReadOnly(tt.carPath)
			require.NoError(t, err)

			// Wait at most 3 seconds for iteration over wanted CIDs.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

			// Fail immediately if error is encountered while iterating over CIDs.
			ctx = blockstore.WithAsyncErrorHandler(ctx, func(err error) { require.Fail(t, "expected no error", "%v", err) })
			t.Cleanup(cancel)

			// Instantiate wanted CIDs channel
			keysChan, err := want.AllKeysChan(ctx)
			require.NoError(t, err)

			gotCids, err := drain(cidIter)
			require.NoError(t, err)

			// Assert CIDs are consistent with the drained iterator
			var i int
			for wantCid := range keysChan {
				gotCid := gotCids[i]
				require.Equal(t, wantCid, gotCid)
				i++
			}

			// Assert drain has fully drained the iterator
			gotCid, gotErr := cidIter.Next()
			require.Equal(t, io.EOF, gotErr)
			require.Equal(t, cid.Undef, gotCid)
		})
	}
}
