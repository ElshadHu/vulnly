package osv

import (
	"context"

	"golang.org/x/sync/errgroup"
)

const maxConcurrentRequests = 25

// HydratedResult contains full vulnerability details for a single query
type HydratedResult struct {
	Vulns []*Vulnerability
}

// HydratedBatchedResponse contains hydrated results for all queries
type HydratedBatchedResponse struct {
	Results []HydratedResult
}

func Hydrate(resp *BatchedResponse) (*HydratedBatchedResponse, error) {
	result := &HydratedBatchedResponse{
		Results: make([]HydratedResult, len(resp.Results)),
	}

	// Preallocate
	for i := range result.Results {
		result.Results[i].Vulns = make([]*Vulnerability, len(resp.Results[i].Vulns))
	}

	g, ctx := errgroup.WithContext(context.Background())
	g.SetLimit(maxConcurrentRequests)

	for batchIdx, response := range resp.Results {
		for resultIdx, minVuln := range response.Vulns {
			id := minVuln.ID
			batchIdx, resultIdx := batchIdx, resultIdx // capture loop vars

			g.Go(func() error {
				if ctx.Err() != nil {
					return nil
				}
				vuln, err := Get(id)
				if err != nil {
					return err
				}
				result.Results[batchIdx].Vulns[resultIdx] = vuln
				return nil
			})
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return result, nil
}
