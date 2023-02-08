package organizations

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func DataSourcePolicy() * schema.Resource {
	return &schema.Resource{
	
		Schema: map[string]*schema.Schema{
			"arn": {
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: verify.ValidARN,
				ExactlyOneOf: []string{"arn", "policy_id", "name"},
			},
			"aws_managed": {
				Type: schema.TypeBool,
				Computed: true,
			},
			"content": {
				Type: schema.TypeString,
				Computed: true,
			},
			"description": {
				Type: schema.TypeString,
				Computed: true,
			},
			"policy_id": {
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
				ExactlyOneOf: []string{"arn", "policy_id", "name"},
			},
			"name": {
				Type: schema.TypeString,
				Optional: true,
				Computed: true,
				ExactlyOneOf: []string{"arn", "policy_id", "name"},
			},
			"tags": tftags.TagsSchemaComputed(),
			"type": {
				Type: schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).OrganizationsConn()
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig //??

	arn := d.Get("arn").(string)
	name := d.Get("name").(string)
	policy_id := d.Get("policy_id").(string)

	if arn != nil {
		// parse policy_id off of the end arn:${Partition}:organizations::${Account}:policy/o-${OrganizationId}/${PolicyType}/p-${PolicyId}
		slices := Split(arn, "/")
		policy_id = slices[len(slices)-1]
	}

	if name != nil {
		// find and set policy id
	}

	if policy_id != nil {
		input := &conn.DescribePolicyInput{
			PolicyId: aws.String(policy_id)
		}

		output, err := conn.DescribePolicyWithContext(ctx, input)

		if err != nil {
			return nil, err
		}

		d.SetId(output.Policy.PolicySummary.Arn)
		d.Set("aws_managed", output.Policy.PolicySummary.AwsManaged)
		d.Set("content", output.Policy.Content)
		d.Set("description", output.Policy.PolicySummary.Description)
		d.set("policy_id", output.Policy.PolicySummary.Id)
		d.set("name", output.Policy.PolicySummary.Name)
		d.set("type", output.Policy.PolicySummary.Type)
	}
}
