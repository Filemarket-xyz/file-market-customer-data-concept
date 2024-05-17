package dune

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/internal/config"
	"github.com/Filemarket-xyz/file-market-customer-data-concept/backend/pkg/logger"
)

type dune struct {
	apiKey string

	client *http.Client
	logger logger.Logger
}

func NewDune(cfg *config.DuneConfig, logger logger.Logger) Dune {
	return &dune{
		apiKey: cfg.ApiKey,
		client: &http.Client{},
		logger: logger,
	}
}

func (d *dune) ExecuteQuery(ctx context.Context, queryId string, req *Req) (*ExecuteResponse, error) {
	subUrl := fmt.Sprintf("/query/%s/execute", queryId)
	data, err := d.makeRequest(ctx, subUrl, http.MethodPost, "", req)
	if err != nil {
		return nil, fmt.Errorf("ExecuteQuery/makeQueryRequest: %w", err)
	}
	var res ExecuteResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("ExecuteQuery/Unmarshal: %w", err)
	}
	return &res, nil
}

func (d *dune) GetLatestQueryResult(ctx context.Context, queryId, filter string) (*QueryResultResponse, error) {
	subUrl := fmt.Sprintf("/query/%s/results", queryId)
	data, err := d.makeRequest(ctx, subUrl, http.MethodGet, filter, nil)
	if err != nil {
		return nil, fmt.Errorf("GetLatestQueryResultRequest/makeQueryRequest: %w", err)
	}
	var res QueryResultResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("GetLatestQueryResultRequest/Unmarshal: %w", err)
	}
	return &res, nil
}

func (d *dune) GetExecutionResult(ctx context.Context, executionId, filter string) (*QueryResultResponse, error) {
	subUrl := fmt.Sprintf("/execution/%s/results", executionId)
	data, err := d.makeRequest(ctx, subUrl, http.MethodGet, filter, nil)
	if err != nil {
		return nil, fmt.Errorf("GetExecutionResultRequest/makeQueryRequest: %w", err)
	}
	var res QueryResultResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("GetExecutionResultRequest/Unmarshal: %w", err)
	}
	return &res, nil
}

func (d *dune) makeRequest(ctx context.Context, subUrl, method, filter string, body interface{}) ([]byte, error) {
	urls := fmt.Sprintf("https://api.dune.com/api/v1%s", subUrl)

	// Create query parameters
	params := url.Values{}
	// params.Set("limit", 10)
	if filter != "" {
		params.Set("filters", filter)
	}
	// params.Set("columns", "tx_from,tx_to,tx_hash,amount_usd")
	// params.Set("sort_by", "amount_usd desc, block_time")

	// Add parameters to URL
	fullURL := fmt.Sprintf("%s?%s", urls, params.Encode())

	var b io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("makeRequest/json.Marshal: %w", err)
		}
		b = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		fullURL,
		b,
	)
	if err != nil {
		return nil, fmt.Errorf("makeRequest/NewRequestWithContext: %w", err)
	}

	req.Header.Add("X-DUNE-API-KEY", d.apiKey)

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("makeRequest/Do: %w", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("makeRequest/ReadAll: %w", err)
	}

	d.logger.Info("DUNE RESP: ", string(data))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("makeRequest ERROR status: ", resp.StatusCode)
	}

	return data, nil
}

func (d *dune) LayerZeroUserFees(ctx context.Context, address string) (*QueryResultResponse, error) {
	subUrl := fmt.Sprintf("/execution/%s/results", "3722138")
	data, err := d.makeRequest(ctx, subUrl, http.MethodGet, fmt.Sprintf("from = '%s'", address), nil)
	if err != nil {
		return nil, fmt.Errorf("LayerZeroUserFees/makeQueryRequest: %w", err)
	}
	var res QueryResultResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("LayerZeroUserFees/Unmarshal: %w", err)
	}
	return &res, nil
}
