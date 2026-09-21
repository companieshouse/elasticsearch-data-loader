<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
| ---- | ------- |
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.3, < 2.0 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | >= 4.54.0, < 6.0 |
| <a name="requirement_vault"></a> [vault](#requirement\_vault) | >= 4.0, < 5.0 |

## Providers

| Name | Version |
| ---- | ------- |
| <a name="provider_aws"></a> [aws](#provider\_aws) | 5.100.0 |
| <a name="provider_vault"></a> [vault](#provider\_vault) | 4.8.0 |

## Modules

| Name | Source | Version |
| ---- | ------ | ------- |
| <a name="module_ecs-service"></a> [ecs-service](#module\_ecs-service) | git@github.com:companieshouse/terraform-modules//aws/ecs/ecs-service | 1.0.334 |
| <a name="module_iac_tags"></a> [iac\_tags](#module\_iac\_tags) | git@github.com:companieshouse/terraform-modules//aws/tagging/iac | tags/1.0.420 |
| <a name="module_owner_tags"></a> [owner\_tags](#module\_owner\_tags) | git@github.com:companieshouse/terraform-modules//aws/tagging/owner | tags/1.0.420 |
| <a name="module_secrets"></a> [secrets](#module\_secrets) | git@github.com:companieshouse/terraform-modules//aws/ecs/secrets | 1.0.334 |

## Resources

| Name | Type |
| ---- | ---- |
| [aws_iam_policy.task_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_policy) | resource |
| [aws_iam_role.task_role](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_role_policy_attachment.task_policy_attachment](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy_attachment) | resource |
| [aws_ecs_cluster.ecs_cluster](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/ecs_cluster) | data source |
| [aws_iam_policy_document.task_assume](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |
| [aws_iam_policy_document.task_policy](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |
| [aws_iam_role.ecs_cluster_iam_role](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_role) | data source |
| [aws_kms_key.kms_key](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/kms_key) | data source |
| [aws_opensearch_domain.opensearch](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/opensearch_domain) | data source |
| [aws_ssm_parameter.global_secret](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/ssm_parameter) | data source |
| [aws_ssm_parameter.secret](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/ssm_parameter) | data source |
| [aws_ssm_parameters_by_path.global_secrets](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/ssm_parameters_by_path) | data source |
| [aws_ssm_parameters_by_path.secrets](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/ssm_parameters_by_path) | data source |
| [aws_subnets.application](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/subnets) | data source |
| [aws_vpc.vpc](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/vpc) | data source |
| [vault_generic_secret.service_secrets](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/data-sources/generic_secret) | data source |
| [vault_generic_secret.shared_s3](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/data-sources/generic_secret) | data source |
| [vault_generic_secret.stack_secrets](https://registry.terraform.io/providers/hashicorp/vault/latest/docs/data-sources/generic_secret) | data source |

## Inputs

| Name | Description | Type | Default | Required |
| ---- | ----------- | ---- | ------- | :------: |
| <a name="input_aws_profile"></a> [aws\_profile](#input\_aws\_profile) | The AWS profile to use for deployment. | `string` | `"development-eu-west-2"` | no |
| <a name="input_aws_region"></a> [aws\_region](#input\_aws\_region) | The AWS region for deployment. | `string` | `"eu-west-2"` | no |
| <a name="input_cloudwatch_alarms_enabled"></a> [cloudwatch\_alarms\_enabled](#input\_cloudwatch\_alarms\_enabled) | Whether to create a standard set of cloudwatch alarms for the service.  Requires an SNS topic to have already been created for the stack. | `bool` | `false` | no |
| <a name="input_create_service_dashboard"></a> [create\_service\_dashboard](#input\_create\_service\_dashboard) | Defines whether a CloudWatch dashboard is created for the default ECS service (true) or not (false) | `bool` | `true` | no |
| <a name="input_desired_task_count"></a> [desired\_task\_count](#input\_desired\_task\_count) | The desired ECS task count for this service | `number` | `0` | no |
| <a name="input_docker_registry"></a> [docker\_registry](#input\_docker\_registry) | The FQDN of the Docker registry. | `string` | n/a | yes |
| <a name="input_elasticsearch_data_loader_version"></a> [elasticsearch\_data\_loader\_version](#input\_elasticsearch\_data\_loader\_version) | The version of the elasticsearch-data-loader container to run. | `string` | n/a | yes |
| <a name="input_environment"></a> [environment](#input\_environment) | The environment name, defined in environments vars. | `string` | n/a | yes |
| <a name="input_hashicorp_vault_password"></a> [hashicorp\_vault\_password](#input\_hashicorp\_vault\_password) | The password used when retrieving configuration from Hashicorp Vault | `string` | n/a | yes |
| <a name="input_hashicorp_vault_username"></a> [hashicorp\_vault\_username](#input\_hashicorp\_vault\_username) | The username used when retrieving configuration from Hashicorp Vault | `string` | n/a | yes |
| <a name="input_max_task_count"></a> [max\_task\_count](#input\_max\_task\_count) | The maximum number of tasks for this service. | `number` | `1` | no |
| <a name="input_min_task_count"></a> [min\_task\_count](#input\_min\_task\_count) | The minimum number of tasks for this service. | `number` | `0` | no |
| <a name="input_opensearch_domain_names"></a> [opensearch\_domain\_names](#input\_opensearch\_domain\_names) | The names of the OpenSearch domains to which the ECS tasks need access. | `list(string)` | <pre>[<br/>  "alphabetical-search"<br/>]</pre> | no |
| <a name="input_required_cpus"></a> [required\_cpus](#input\_required\_cpus) | The required cpu resource for this service. 1024 here is 1 vCPU | `number` | `256` | no |
| <a name="input_required_memory"></a> [required\_memory](#input\_required\_memory) | The required memory for this service | `number` | `512` | no |
| <a name="input_service_autoscale_enabled"></a> [service\_autoscale\_enabled](#input\_service\_autoscale\_enabled) | Whether to enable service autoscaling, including scheduled autoscaling | `bool` | `false` | no |
| <a name="input_service_autoscale_scale_in_cooldown"></a> [service\_autoscale\_scale\_in\_cooldown](#input\_service\_autoscale\_scale\_in\_cooldown) | Cooldown in seconds for ECS Service scale in | `number` | `300` | no |
| <a name="input_service_autoscale_scale_out_cooldown"></a> [service\_autoscale\_scale\_out\_cooldown](#input\_service\_autoscale\_scale\_out\_cooldown) | Cooldown in seconds for ECS Service scale out | `number` | `120` | no |
| <a name="input_service_autoscale_target_value_cpu"></a> [service\_autoscale\_target\_value\_cpu](#input\_service\_autoscale\_target\_value\_cpu) | Target CPU percentage for the ECS Service to autoscale on | `number` | `80` | no |
| <a name="input_service_scaledown_schedule"></a> [service\_scaledown\_schedule](#input\_service\_scaledown\_schedule) | The schedule to use when scaling down the number of tasks to zero. | `string` | `""` | no |
| <a name="input_service_scaleup_schedule"></a> [service\_scaleup\_schedule](#input\_service\_scaleup\_schedule) | The schedule to use when scaling up the number of tasks to their normal desired level. | `string` | `""` | no |
| <a name="input_ssm_version_prefix"></a> [ssm\_version\_prefix](#input\_ssm\_version\_prefix) | String to use as a prefix to the names of the variables containing variables and secrets version. | `string` | `"SSM_VERSION_"` | no |
| <a name="input_use_capacity_provider"></a> [use\_capacity\_provider](#input\_use\_capacity\_provider) | Whether to use a capacity provider instead of setting a launch type for the service | `bool` | `true` | no |
| <a name="input_use_fargate"></a> [use\_fargate](#input\_use\_fargate) | If true, sets the required capabilities for all containers in the task definition to use FARGATE, false uses EC2 | `bool` | `true` | no |
| <a name="input_use_set_environment_files"></a> [use\_set\_environment\_files](#input\_use\_set\_environment\_files) | Toggle default global and shared  environment files | `bool` | `false` | no |

## Outputs

No outputs.
<!-- END_TF_DOCS -->