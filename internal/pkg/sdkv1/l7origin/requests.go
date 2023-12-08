package l7origin

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

const l7OriginPath = "l7/origin"

// GetByID returns a single resource by its id.
func GetByID(ctx context.Context, client *v1.Client, l7ResourceID int, ID int) (*Data, *v1.ResponseResult, error) {
	url := strings.Join([]string{client.Endpoint, l7OriginPath, strconv.Itoa(l7ResourceID), strconv.Itoa(ID)}, "/")
	responseResult, err := client.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	// Extract a result from the response body.
	result := &Data{}
	err = responseResult.ExtractResult(result)
	if err != nil {
		return nil, responseResult, err
	}

	return result, responseResult, nil
}

// List gets a list of all origins.
func List(ctx context.Context, client *v1.Client) ([]*Item, *v1.ResponseResult, error) {
	url := strings.Join([]string{client.Endpoint, l7OriginPath}, "/")
	responseResult, err := client.DoRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	// Extract origins from the response body.
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
	url := strings.Join([]string{client.Endpoint, l7OriginPath}, "/")
	requestBody, err := json.Marshal(opts)
	if err != nil {
		return nil, nil, err
	}

	if opts.IP == "" {
		return nil, nil, fmt.Errorf("sp-go: origin IP must be not empty")
	}

	responseResult, err := client.DoRequest(ctx, http.MethodPost, url, bytes.NewReader(requestBody))
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	// Extract origin from the response body.
	result := &Data{}
	err = responseResult.ExtractResult(result)
	if err != nil {
		return nil, responseResult, err
	}

	return result, responseResult, nil
}

// Delete deletes a single origin by its id.
func Delete(ctx context.Context, client *v1.Client, opts *DeleteOpts) (*DataDelete, *v1.ResponseResult, error) {
	url := strings.Join([]string{client.Endpoint, l7OriginPath}, "/")
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

	// Extract origin from the response body.
	origin := &DataDelete{}
	err = responseResult.ExtractResult(origin)
	if err != nil {
		return nil, responseResult, err
	}

	return origin, responseResult, nil
}

// Update deletes a single origin by its id.
func Update(ctx context.Context, client *v1.Client, item *Item) (*Data, *v1.ResponseResult, error) {
	url := strings.Join([]string{client.Endpoint, l7OriginPath}, "/")
	requestBody, err := json.Marshal(item)
	if err != nil {
		return nil, nil, err
	}

	responseResult, err := client.DoRequest(ctx, http.MethodPut, url, bytes.NewReader(requestBody))
	if err != nil {
		return nil, nil, err
	}
	if responseResult.Err != nil {
		return nil, responseResult, responseResult.Err
	}

	// Extract origin from the response body.
	result := &Data{}
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
