package l7resource

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	v1 "terraform-provider-servicepipe/internal/pkg/sdkv1"
)

const l7ResourcePath = "l7/resource"

// GetByID returns a single resource by its id.
func GetByID(ctx context.Context, client *v1.Client, l7ResourceID int) (*Data, *v1.ResponseResult, error) {
	url := strings.Join([]string{client.Endpoint, l7ResourcePath, strconv.Itoa(l7ResourceID)}, "/")
	responseResult, err := client.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	// Extract a l7Resource from the response body.
	l7Resource := &Data{}
	err = responseResult.ExtractResult(l7Resource)
	if err != nil {
		return nil, responseResult, err
	}

	return l7Resource, responseResult, nil
}

// List gets a list of all domains.
func List(ctx context.Context, client *v1.Client) ([]*Item, *v1.ResponseResult, error) {
	url := strings.Join([]string{client.Endpoint, l7ResourcePath}, "/")
	responseResult, err := client.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	// Extract domains from the response body.
	var dataitems *DataItems
	err = responseResult.ExtractResult(&dataitems)
	if err != nil {
		return nil, responseResult, err
	}
	convertedItems := convertToSliceOfPointers(dataitems.DataItems.ResultItems.Items)
	return convertedItems, responseResult, nil
}

// Create requests a creation of a new domain.
func Create(ctx context.Context, client *v1.Client, opts *CreateOpts) (*Data, *v1.ResponseResult, error) {
	url := strings.Join([]string{client.Endpoint, l7ResourcePath}, "/")
	requestBody, err := json.Marshal(opts)
	if err != nil {
		return nil, nil, err
	}

	if opts.OriginData == "" {
		return nil, nil, fmt.Errorf("sp-go: originData must be not empty")
	}

	responseResult, err := client.DoRequest(ctx, http.MethodPost, url, bytes.NewReader(requestBody))
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	// Extract domain from the response body.
	domain := &Data{}
	err = responseResult.ExtractResult(domain)
	if err != nil {
		return nil, responseResult, err
	}

	return domain, responseResult, nil
}

// Delete deletes a single domain by its id.
func Delete(ctx context.Context, client *v1.Client, opts *DeleteOpts) (*DataDelete, *v1.ResponseResult, error) {
	url := strings.Join([]string{client.Endpoint, l7ResourcePath}, "/")
	requestBody, err := json.Marshal(opts)
	if err != nil {
		return nil, nil, err
	}

	responseResult, err := client.DoRequest(ctx, http.MethodDelete, url, bytes.NewReader(requestBody))
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	// Extract domain from the response body.
	result := &DataDelete{}
	err = responseResult.ExtractResult(result)
	if err != nil {
		return nil, responseResult, err
	}

	return result, responseResult, nil
}

func convertToSliceOfPointers(items []Item) []*Item {
	pointers := make([]*Item, len(items))
	for i := range items {
		pointers[i] = &items[i]
	}
	return pointers
}
