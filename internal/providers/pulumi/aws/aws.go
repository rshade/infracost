package aws

import (
	"strings"

	"github.com/infracost/infracost/internal/schema"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

var DefaultProviderRegion = "us-east-1"

func GetDefaultRefIDFunc(d *schema.ResourceData) []string {
	defaultRefs := []string{d.RawValues.Get("urn").String()}
	return defaultRefs
}

func GetSpecialContext(d *schema.ResourceData) map[string]interface{} {

	specialContexts := make(map[string]interface{})

	if strings.HasPrefix(d.Get("region").String(), "cn-") {
		specialContexts["isAWSChina"] = true
	}

	return specialContexts
}

func GetResourceRegion(resourceType string, v gjson.Result) string {
	// If a region key exists in the values use that

	return v.Get("region").String()

}

func ParseTags(resourceType string, v gjson.Result) map[string]string {
	tags := make(map[string]string)

	for k, v := range v.Get("tags").Map() {
		if k == "__defaults" {
			continue
		}
		tags[k] = v.String()
	}
	log.Debugf("tags %s", tags)
	return tags
}

func GetAWSResourceTypes() map[string]string {
	resourceTypes := map[string]string{
		"aws_ec2_eip":             "aws_eip",
		"aws_ec2_natgateway":      "aws_nat_gateway",
		"aws_rds_instance":        "aws_db_instance",
		"aws_rds_clusterinstance": "aws_rds_cluster_instance",
		"aws_eks_nodegroup":       "aws_eks_node_group",
		"aws_ebs_snapshotcopy":    "aws_ebs_snapshot_copy",
	}
	return resourceTypes
}