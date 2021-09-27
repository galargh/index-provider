package adminserver

import "github.com/ipfs/go-cid"

// ErrorRes represents a response that captures information about failure to handle a request.
type ErrorRes struct {
	// The human-readable message that provides hints about the failure cause.
	Message string `json:"message"`
}

type (
	// ConnectReq request to connect to a given multiaddr.
	ConnectReq struct {
		Maddr string `json:"maddr"`
	}
	// ConnectRes represents successful response to ConnectReq request.
	ConnectRes struct { // Empty placeholder used to return an empty JSON object in body.
	}
)

type (
	// ImportCarReq represents a request for importing a CAR file.
	ImportCarReq struct {
		// The path to the CAR file
		Path string `json:"path"`
		// The optional lookup key associated to the CAR. If not provided, one will be generated.
		Key []byte `json:"key"`
		// The optional metadata.
		Metadata []byte `json:"metadata"`
	}
	// ImportCarRes represents the response to an ImportCarReq.
	ImportCarRes struct {
		// The lookup Key associated to the imported CAR.
		Key []byte `json:"key"`
		// The CID of the advertisement generated as a result of import.
		AdvId cid.Cid `json:"adv_id"`
	}
)