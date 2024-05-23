package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type IndexQueryParameters struct {
	Cve              string `json:"cve"`
	Alias            string `json:"alias"`
	Iava             string `json:"iava"`
	LastModStartDate string `json:"lastModStartDate"`
	LastModEndDate   string `json:"lastModEndDate"`
	PubStartDate     string `json:"pubStartDate"`
	PubEndDate       string `json:"pubEndDate"`
	ThreatActor      string `json:"threat_actor"`
	MitreId          string `json:"mitre_id"`
	MispId           string `json:"misp_id"`
	Ransomware       string `json:"ransomware"`
	Botnet           string `json:"botnet"`
	Limit            int    `json:"limit"`
	Sort             string `json:"sort"`
	Order            string `json:"order"`
	Page             int    `json:"page"`
	Cursor           string `json:"cursor"`
	NextCursor       string `json:"next_cursor"`
	PrevCursor       string `json:"prev_cursor"` // this may not exist
}

type IndexMeta struct {
	Timestamp      string                `json:"timestamp"`
	Index          string                `json:"index"`
	Limit          int                   `json:"limit"`
	TotalDocuments int                   `json:"total_documents"`
	Sort           string                `json:"sort"`
	Parameters     []IndexMetaParameters `json:"parameters"`
	Order          string                `json:"order"`
	Page           int                   `json:"page"`
	TotalPages     int                   `json:"total_pages"`
	MaxPages       int                   `json:"max_pages"`
	FirstItem      int                   `json:"first_item"`
	LastItem       int                   `json:"last_item"`
}

type IndexMetaParameters struct {
	Name   string `json:"name"`
	Format string `json:"format"`
}

type IndexResponse struct {
	Benchmark float64       `json:"_benchmark"`
	Meta      IndexMeta     `json:"_meta"`
	Data      []interface{} `json:"data"`
}

// add method to set query parameters
func setIndexQueryParameters(query url.Values, queryParameters ...IndexQueryParameters) {
	for _, queryParameter := range queryParameters {
		if queryParameter.Cve != "" {
			query.Add("cve", queryParameter.Cve)
		}
		if queryParameter.Alias != "" {
			query.Add("alias", queryParameter.Alias)
		}
		if queryParameter.Iava != "" {
			query.Add("iava", queryParameter.Iava)
		}
		if queryParameter.LastModStartDate != "" {
			query.Add("lastModStartDate", queryParameter.LastModStartDate)
		}
		if queryParameter.LastModEndDate != "" {
			query.Add("lastModEndDate", queryParameter.LastModEndDate)
		}
		if queryParameter.PubStartDate != "" {
			query.Add("pubStartDate", queryParameter.PubStartDate)
		}
		if queryParameter.PubEndDate != "" {
			query.Add("pubEndDate", queryParameter.PubEndDate)
		}
		if queryParameter.ThreatActor != "" {
			query.Add("threat_actor", queryParameter.ThreatActor)
		}
		if queryParameter.MitreId != "" {
			query.Add("mitre_id", queryParameter.MitreId)
		}
		if queryParameter.MispId != "" {
			query.Add("misp_id", queryParameter.MispId)
		}
		if queryParameter.Ransomware != "" {
			query.Add("ransomware", queryParameter.Ransomware)
		}
		if queryParameter.Botnet != "" {
			query.Add("botnet", queryParameter.Botnet)
		}
		if queryParameter.Limit != 0 {
			query.Add("limit", fmt.Sprintf("%d", queryParameter.Limit))
		}
		if queryParameter.Sort != "" {
			query.Add("sort", queryParameter.Sort)
		}
		if queryParameter.Order != "" {
			query.Add("order", queryParameter.Order)
		}
		if queryParameter.Page != 0 {
			query.Add("page", fmt.Sprintf("%d", queryParameter.Page))
		}
		if queryParameter.Cursor != "" {
			query.Add("cursor", queryParameter.Cursor)
		}
		if queryParameter.NextCursor != "" {
			query.Add("next_cursor", queryParameter.NextCursor)
		}
		if queryParameter.PrevCursor != "" {
			query.Add("prev_cursor", queryParameter.PrevCursor)
		}
	}
}

// https://docs.vulncheck.com/api/indice
func (c *Client) GetIndex(index string, queryParameters ...IndexQueryParameters) (responseJSON *IndexResponse, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", c.GetUrl()+"/v3/index/"+url.QueryEscape(index), nil)
	if err != nil {
		panic(err)
	}

	c.SetAuthHeader(req)

	query := req.URL.Query()
	setIndexQueryParameters(query, queryParameters...)
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var metaError MetaError
		_ = json.NewDecoder(resp.Body).Decode(&metaError)

		return nil, fmt.Errorf("error: %v", metaError.Errors)
	}

	_ = json.NewDecoder(resp.Body).Decode(&responseJSON)

	return responseJSON, nil
}

// Strings representation of the response
func (r IndexResponse) String() string {
	return fmt.Sprintf("Benchmark: %f\nMeta: %v\nData: %v\n", r.Benchmark, r.Meta, r.Data)
}

// GetData - Returns the data from the response
func (r IndexResponse) GetData() []interface{} {
	return r.Data
}
