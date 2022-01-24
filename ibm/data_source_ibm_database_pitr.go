// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/cloud-databases-go-sdk/clouddatabasesv5"
)

func dataSourceIBMDatabasePitr() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMDatabasePitrRead,

		Schema: map[string]*schema.Schema{
			"resource_instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Deployment ID.",
			},
			"earliest_point_in_time_recovery_time": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIBMDatabasePitrRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cloudDatabasesClient, err := meta.(ClientSession).CloudDatabasesV5()
	if err != nil {
		return diag.FromErr(err)
	}

	getPitrDataOptions := &clouddatabasesv5.GetPitRdataOptions{}

	getPitrDataOptions.SetID(d.Get("resource_instance_id").(string))

	pointInTimeRecoveryData, response, err := cloudDatabasesClient.GetPitRdataWithContext(context, getPitrDataOptions)
	if err != nil {
		log.Printf("[DEBUG] GetPitrDataWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetPitrDataWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIBMDatabasePitrID(d))

	if err = d.Set("earliest_point_in_time_recovery_time", pointInTimeRecoveryData.EarliestPointInTimeRecoveryTime); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting earliest_point_in_time_recovery_time: %s", err))
	}

	return nil
}

// dataSourceIBMDatabasePitrID returns a reasonable ID for the list.
func dataSourceIBMDatabasePitrID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
