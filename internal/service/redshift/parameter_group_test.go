package redshift_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/redshift"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfredshift "github.com/hashicorp/terraform-provider-aws/internal/service/redshift"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccRedshiftParameterGroup_basic(t *testing.T) {
	var v redshift.ClusterParameterGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_redshift_parameter_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshift.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckParameterGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccParameterGroupConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckParameterGroupExists(resourceName, &v),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "redshift", fmt.Sprintf("parametergroup:%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(resourceName, "family", "redshift-1.0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "parameter.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "parameter.*", map[string]string{
						"name":  "require_ssl",
						"value": "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "parameter.*", map[string]string{
						"name":  "query_group",
						"value": "example",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "parameter.*", map[string]string{
						"name":  "enable_user_activity_logging",
						"value": "true",
					}),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRedshiftParameterGroup_disappears(t *testing.T) {
	var v redshift.ClusterParameterGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_redshift_parameter_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshift.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckParameterGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccParameterGroupConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterGroupExists(resourceName, &v),
					acctest.CheckResourceDisappears(acctest.Provider, tfredshift.ResourceParameterGroup(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccRedshiftParameterGroup_update(t *testing.T) {
	var v redshift.ClusterParameterGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_redshift_parameter_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshift.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckParameterGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccParameterGroupConfig_noParameters(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckParameterGroupExists(resourceName, &v),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "redshift", fmt.Sprintf("parametergroup:%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(resourceName, "family", "redshift-1.0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "parameter.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccParameterGroupConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckParameterGroupExists(resourceName, &v),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "redshift", fmt.Sprintf("parametergroup:%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "description", "Managed by Terraform"),
					resource.TestCheckResourceAttr(resourceName, "family", "redshift-1.0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "parameter.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "parameter.*", map[string]string{
						"name":  "require_ssl",
						"value": "true",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "parameter.*", map[string]string{
						"name":  "query_group",
						"value": "example",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "parameter.*", map[string]string{
						"name":  "enable_user_activity_logging",
						"value": "true",
					}),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccRedshiftParameterGroup_tags(t *testing.T) {
	var v redshift.ClusterParameterGroup
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_redshift_parameter_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, redshift.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckParameterGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccParameterGroupConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterGroupExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "description", "Test parameter group for terraform"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccParameterGroupConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterGroupExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccParameterGroupConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckParameterGroupExists(resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckParameterGroupDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).RedshiftConn()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_redshift_parameter_group" {
			continue
		}

		_, err := tfredshift.FindParameterGroupByName(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("Redshift Parameter Group %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccCheckParameterGroupExists(n string, v *redshift.ClusterParameterGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Redshift Parameter Group ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).RedshiftConn()

		output, err := tfredshift.FindParameterGroupByName(conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccParameterGroupConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_redshift_parameter_group" "test" {
  name   = %[1]q
  family = "redshift-1.0"

  parameter {
    name  = "require_ssl"
    value = "true"
  }

  parameter {
    name  = "query_group"
    value = "example"
  }

  parameter {
    name  = "enable_user_activity_logging"
    value = "true"
  }
}
`, rName)
}

func testAccParameterGroupConfig_noParameters(rName string) string {
	return fmt.Sprintf(`
resource "aws_redshift_parameter_group" "test" {
  name   = %[1]q
  family = "redshift-1.0"
}
`, rName)
}

func testAccParameterGroupConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_redshift_parameter_group" "test" {
  name        = %[1]q
  family      = "redshift-1.0"
  description = "Test parameter group for terraform"

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccParameterGroupConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_redshift_parameter_group" "test" {
  name        = %[1]q
  family      = "redshift-1.0"
  description = "Test parameter group for terraform"

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}
