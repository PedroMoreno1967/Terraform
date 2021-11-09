package costmanagement_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CostManagementExportSubscription struct {
}

func TestAccCostManagementExportSubscription_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cost_management_export_subscription", "test")
	r := CostManagementExportSubscription{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCostManagementExportSubscription_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cost_management_export_subscription", "test")
	r := CostManagementExportSubscription{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t CostManagementExportSubscription) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.CostManagementExportID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.CostManagement.ExportClient.Get(ctx, id.Scope, id.Name, "")
	if err != nil {
		return nil, fmt.Errorf("retrieving (%s): %+v", *id, err)
	}

	return utils.Bool(resp.ExportProperties != nil), nil
}

func (CostManagementExportSubscription) basic(data acceptance.TestData) string {
	start := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	end := time.Now().AddDate(0, 0, 2).Format("2006-01-02")

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "test" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cm-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_cost_management_export_subscription" "test" {
  name                    = "accrg%d"
  subscription_id         = data.azurerm_subscription.test.id
  recurrence_type         = "Monthly"
  recurrence_period_start = "%sT00:00:00Z"
  recurrence_period_end   = "%sT00:00:00Z"

  export_data_storage_location {
    storage_account_id = azurerm_storage_account.test.id
    container_name     = "acctestcontainer"
    root_folder_path   = "/root"
  }

  export_data_definition {
    type       = "Usage"
    time_frame = "TheLastMonth"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, start, end)
}

func (CostManagementExportSubscription) update(data acceptance.TestData) string {
	start := time.Now().AddDate(0, 3, 0).Format("2006-01-02")
	end := time.Now().AddDate(0, 4, 0).Format("2006-01-02")

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "test" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cm-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_cost_management_export_subscription" "test" {
  name                    = "accrg%d"
  subscription_id         = data.azurerm_subscription.test.id
  recurrence_type         = "Monthly"
  recurrence_period_start = "%sT00:00:00Z"
  recurrence_period_end   = "%sT00:00:00Z"

  export_data_storage_location {
    storage_account_id = azurerm_storage_account.test.id
    container_name     = "acctestcontainer"
    root_folder_path   = "/root/updated"
  }

  export_data_definition {
    type       = "Usage"
    time_frame = "WeekToDate"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger, start, end)
}

