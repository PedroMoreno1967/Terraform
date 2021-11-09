package costmanagement

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/validate"
	subscriptionValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/subscription/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SubscriptionCostManagementExportResource struct {
	base costManagementExportBaseResource
}

var _ sdk.Resource = SubscriptionCostManagementExportResource{}

func (r SubscriptionCostManagementExportResource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"subscription_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: subscriptionValidate.SubscriptionID,
		},
	}
	return r.base.arguments(schema)
}

func (r SubscriptionCostManagementExportResource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.attributes()
}

func (r SubscriptionCostManagementExportResource) ModelObject() interface{} {
	return nil
}

func (r SubscriptionCostManagementExportResource) ResourceType() string{
	return "azurerm_cost_management_export_subscription"
}

func (r SubscriptionCostManagementExportResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.SubscriptionCostManagementExportID
}

func (r SubscriptionCostManagementExportResource) Create() sdk.ResourceFunc {
	return r.base.createFunc(r.ResourceType(), "subscription_id")
}

func (r SubscriptionCostManagementExportResource) Read() sdk.ResourceFunc {
	return r.base.readFunc("subscription_id")
}

func (r SubscriptionCostManagementExportResource) Delete() sdk.ResourceFunc {
	return r.base.deleteFunc()
}

func (r SubscriptionCostManagementExportResource) Update() sdk.ResourceFunc {
	return r.base.updateFunc()
}