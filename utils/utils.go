package utils

import "github.com/prebid/openrtb/v20/openrtb2"

func GenerateBidResponse(bidRequest *openrtb2.BidRequest) *openrtb2.BidResponse {
	// Example: Create a BidResponse with a sample bid
	bid := openrtb2.Bid{
		ID:    "1",
		ImpID: bidRequest.Imp[0].ID,
		Price: 0.10,
		AdM:   "<html><body><h1>Example Ad</h1></body></html>",
	}

	return &openrtb2.BidResponse{
		ID: bidRequest.ID,
		SeatBid: []openrtb2.SeatBid{
			{
				Bid: []openrtb2.Bid{bid},
			},
		},
	}
}
