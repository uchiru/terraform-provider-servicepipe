


```Golang
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	v1 "terraform-provider-servicepipe/internal/pkg/sdkv1"
	"terraform-provider-servicepipe/internal/pkg/sdkv1/l7origin"
	l7resource "terraform-provider-servicepipe/internal/pkg/sdkv1/l7resource"
)

const (
	resourceName = "testdomain.xyz"
	firstOrigin  = "190.90.160.30"
)

func main() {
	token := os.Getenv("SERVICEPIPE_API_TOKEN")

	if token == "" {
		fmt.Println("Error: Environment variable 'SERVICEPIPE_API_TOKEN' is not set or is empty.")
		os.Exit(1)
	}

	// Initialize the Domains API V1 client
	client := v1.NewClientV1WithDefaultEndpoint(token)

	fmt.Println("================= Step 1: check existing resource =================")
	listResources, _, err := l7resource.List(context.Background(), client)
	if err != nil {
		log.Fatal(err)
	}

	_, ok := findItemByName(listResources, resourceName)
	if !ok {
		fmt.Println("================= Step 2: create =================")
		createOpts := &l7resource.CreateOpts{
			L7ResourceName: resourceName,
			OriginData:     firstOrigin,
		}

		// Create domain
		result, _, err := l7resource.Create(context.Background(), client, createOpts)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Added new l7 resource: %+v\n", result.Data.Result.L7ResourceName)
	}

	fmt.Println("================= Step 3: check existing resource =================")
	resources, _, err := l7resource.List(context.Background(), client)
	if err != nil {
		log.Fatal(err)
	}

	rs, ok := findItemByName(resources, resourceName)
	if !ok {
		fmt.Printf("Resource not found by name %s", resourceName)
		return
	}

	printStruct(rs)

	fmt.Println("================= Step 4: update =================")
	rs.Wwwredir = 1

	updateResult, _, err := l7resource.Update(context.Background(), client, rs)
	if err != nil {
		log.Fatal(err)
	}

	printStruct(updateResult.Data.Result)

	fmt.Println("================= Step 5: add new origin =================")
	originListOpts := &l7origin.ListOpts{
		L7ResourceID: rs.L7ResourceID,
	}

	or, _, err := l7origin.List(context.Background(), client, originListOpts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(or)

	newOrigin := "190.90.160.99"
	_, ok = findOriginByIP(or, newOrigin)
	if !ok {
		originCreateOpts := &l7origin.CreateOpts{
			L7ResourceID: rs.L7ResourceID,
			IP:           newOrigin,
			Weight:       1,
			Mode:         "primary",
		}

		origCreateResult, _, err := l7origin.Create(context.Background(), client, originCreateOpts)
		if err != nil {
			log.Fatal(err)
		}

		printStruct(origCreateResult.Data.Result)
	}

	time.Sleep(20 * time.Second)
	fmt.Println("================= Step 6: update existing origin =================")
	origins, _, err := l7origin.List(context.Background(), client, originListOpts)
	if err != nil {
		log.Fatal(err)
	}

	origRs, ok := findOriginByIP(origins, firstOrigin)
	if !ok {
		fmt.Printf("Origin not found %s", firstOrigin)
		return
	}

	origRs.L7ResourceID = rs.L7ResourceID
	origRs.IP = "190.90.160.30"
	origRs.Mode = "backup"
	origRs.Weight = 40

	origUpdateResult, _, err := l7origin.Update(context.Background(), client, origRs)
	if err != nil {
		log.Fatal(err)
	}

	printStruct(origUpdateResult.Data.Result)

	time.Sleep(20 * time.Second)
	fmt.Println("================= Step 7: delete origin =================")
	deleteOriginOpts := &l7origin.DeleteOpts{
		L7ResourceID: rs.L7ResourceID,
		ID:           origRs.ID,
	}
	origDelResult, _, err := l7origin.Delete(context.Background(), client, deleteOriginOpts)
	if err != nil {
		printStruct(origDelResult)
		log.Fatal(err)
	}

	if origDelResult.Data.Result == "ok" {
		fmt.Println("Resource deleted %s", origRs.IP)
	}

	time.Sleep(20 * time.Second)
	fmt.Println("================= Step 8: delete resource =================")
	deleteOpts := &l7resource.DeleteOpts{
		L7ResourceID: rs.L7ResourceID,
	}
	resul, _, err := l7resource.Delete(context.Background(), client, deleteOpts)
	if err != nil {
		log.Fatal(err)
	}

	if resul.Data.Result == "ok" {
		fmt.Println("Resource deleted %s", resourceName)
	}
}

func printStruct(v interface{}) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fmt.Printf("%s: %v\n", typ.Field(i).Name, field.Interface())
	}
}

func findItemByName(items []*l7resource.Item, name string) (*l7resource.Item, bool) {
	for _, item := range items {
		if item.L7ResourceName == name {
			return item, true
		}
	}
	return &l7resource.Item{}, false // Return an empty Person and false if not found
}

func findOriginByIP(items []*l7origin.Item, IP string) (*l7origin.Item, bool) {
	for _, item := range items {
		if item.IP == IP {
			return item, true
		}
	}
	return &l7origin.Item{}, false // Return an empty Person and false if not found
}

```
