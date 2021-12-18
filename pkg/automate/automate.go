package automate

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

func HandleTick(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	if err := CleanUpResourceGroups(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CleanUpResourceGroups() error {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Printf("Error while loading Azure credentials %s", err)
		return err
	}
	s, err := getSubscriptionId()
	if err != nil {
		return err
	}

	client := armresources.NewResourceGroupsClient(s, cred, nil)

	pager := client.List(nil)
	for pager.NextPage(context.Background()) {
		if err := pager.Err(); err != nil {
			log.Printf("Failed to load next page %s", err)
			return err
		}
		updates := make([]string, 0)
		removals := make([]string, 0)
		for _, rg := range pager.PageResponse().ResourceGroupListResult.Value {
			log.Printf("\nLooking at Resource Group %s\n", *rg.Name)

			if hasKeeperTag(rg.Tags) {
				log.Printf("Resource Group '%s' has keeper tag. Will leave it as it is.", *rg.Name)
				continue
			} else if !hasExpirationTag(rg.Tags) {
				log.Printf("Resource Group '%s' will be marked for expiration.", *rg.Name)
				updates = append(updates, *rg.Name)
				continue
			} else if isExpired(rg.Tags) {
				log.Printf("Resource Group '%s' is expired... Will delete it.", *rg.Name)
				removals = append(removals, *rg.Name)
				continue
			} else {
				log.Printf("Resource Group '%s' already marked for expiration. Will check again at next run...", *rg.Name)
			}
		}
		applyUpdates(updates, client)
		applyRemovals(removals, client)
	}
	return nil
}

func applyRemovals(removals []string, client *armresources.ResourceGroupsClient) error {
	for _, r := range removals {
		p, err := client.BeginDelete(context.Background(), r, nil)
		if err != nil {
			log.Printf("Error while trying to begin removal of Resource Group '%s'", r)
			return err
		}
		_, err = p.PollUntilDone(context.Background(), 15*time.Second)
		if err != nil {
			log.Printf("Error while removing Resource Group '%s'", r)
			return err
		}
	}
	return nil
}

func applyUpdates(updates []string, client *armresources.ResourceGroupsClient) error {
	for _, u := range updates {
		r, err := client.Update(context.Background(), u, armresources.ResourceGroupPatchable{
			Tags: map[string]*string{
				getExpirationTagName(): to.StringPtr(getExpiration()),
			}},
			nil)
		if err != nil {
			log.Printf("Error while updating Resource Group '%s': %s", u, err)
			return err
		}
		log.Printf("Resource Group '%s' updated successfully", *r.Name)
	}
	return nil
}

func getSubscriptionId() (string, error) {
	sub, ok := os.LookupEnv("SUBSCRIPTION_ID")
	if !ok {
		return "", errors.New("the Azure Subscription ID not found. Please set env var SUBSCRIPTION_ID")
	}
	return sub, nil
}
