package rest

import (
	"fmt"
	"net/url"
	"path"
)

// MetricsLatest gets the latest metric for entity.
// metric can be host, guest, storage, pool, and cluster
// resolution can be 20, 900, 3600, or 21600
func (client *Client) MetricsLatest(metric string, entity string, resolution uint) ([]byte, error) {
	path := path.Join("metrics", metric, entity, "latest")
	query := url.Values{}
	if resolution > 0 {
		query.Set("resolution", fmt.Sprintf("%d", resolution))
	}
	return client.request("GET", path+"?"+query.Encode(), nil)
}

//MetricsExport export metrics as json or csv
// output can be json or csv
func (client *Client) MetricsExport(metric string, entity string, resolution uint, output string) ([]byte, error) {
	path := path.Join("metrics", metric, entity, "export")
	query := url.Values{}
	if resolution > 0 {
		query.Set("resolution", fmt.Sprintf("%d", resolution))
	}
	if output == "" {
		output = "json"
	}
	query.Set("output", output)
	return client.request("GET", path+"?"+query.Encode(), nil)
}
