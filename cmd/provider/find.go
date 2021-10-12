package main

import (
	"fmt"

	"github.com/filecoin-project/go-indexer-core"
	httpfinderclient "github.com/filecoin-project/storetheindex/api/v0/finder/client/http"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
	"github.com/urfave/cli/v2"
)

var FindCmd = &cli.Command{
	Name:   "find",
	Usage:  "Query an indexer for indexed content",
	Flags:  findFlags,
	Action: findCommand,
}

func findCommand(cctx *cli.Context) error {
	cli, err := httpfinderclient.New(cctx.String("indexer"))
	if err != nil {
		return err
	}

	mhArgs := cctx.StringSlice("mh")
	cidArgs := cctx.StringSlice("cid")
	mhs := make([]multihash.Multihash, 0, len(mhArgs)+len(cidArgs))
	for i := range mhArgs {
		m, err := multihash.FromB58String(mhArgs[i])
		if err != nil {
			return err
		}
		mhs = append(mhs, m)
	}
	for i := range cidArgs {
		c, err := cid.Decode(cidArgs[i])
		if err != nil {
			return err
		}
		mhs = append(mhs, c.Hash())
	}

	resp, err := cli.FindBatch(cctx.Context, mhs)
	if err != nil {
		return err
	}

	if len(resp.MultihashResults) == 0 {
		fmt.Println("index not found")
		return nil
	}

	fmt.Println("Content providers:")
	for i := range resp.MultihashResults {
		fmt.Println("   Multihash:", resp.MultihashResults[i].Multihash.B58String(), "==>")
		for _, pr := range resp.MultihashResults[i].ProviderResults {
			v := indexer.Value{
				Metadata: pr.Metadata,
			}
			protocol, data, err := v.GetData()
			if err != nil {
				return fmt.Errorf("cannot decode metadata: %s", err)
			}
			fmt.Println("       Provider:", pr.Provider, "Proto:", protocol, "Data:", string(data))
		}
	}

	return nil
}