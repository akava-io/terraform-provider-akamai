provider "akamai" {
  edgerc = "~/.edgerc"
  appsec_section = "global"
}

/*
data "akamai_appsec_configuration" "appsecconfig" {
  name = "Akamai Tools"
  
}

output "configs" {
  value = data.akamai_appsec_configuration.appsecconfig.config_id
}
*/

data "akamai_appsec_configuration" "appsecconfigedge" {
  name = "Example for EDGE"
  version = 11 
}

resource "akamai_appsec_configuration_clone" "appsecconfigurationclone" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    create_from_version = 29 //data.akamai_appsec_configuration.appsecconfigedge.latest_version 
    rule_update  = false
   }

/*
data "akamai_appsec_selectable_hostnames" "appsecselectablehostnames" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version   
}*/
/*
output "selectablehostnames" {
  value = data.akamai_appsec_selectable_hostnames.appsecselectablehostnames.host_names
}
*/
output "configsedge" {
  value = data.akamai_appsec_configuration.appsecconfigedge.config_id
}

output "configsedgelatestversion" {
  value = data.akamai_appsec_configuration.appsecconfigedge.latest_version 
}

output "configsedgeconfiglist" {
  value = data.akamai_appsec_configuration.appsecconfigedge.config_list
}

output "configsedgeconfigversion" {
  value = data.akamai_appsec_configuration.appsecconfigedge.version
}

/*
data "akamai_appsec_export_configuration" "export" {
  config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
  version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
  
}

output "exportconfig" {
  value = data.akamai_appsec_export_configuration.export.json
}*/


resource "akamai_appsec_selected_hostnames" "appsecselectedhostnames" {
    //config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    //version = akamai_appsec_configuration_clone.appsecconfigurationclone.version //data.akamai_appsec_configuration.appsecconfigedge.latest_version 
    config_id = akamai_appsec_configuration_clone.appsecconfigurationclone.config_id
    version = akamai_appsec_configuration_clone.appsecconfigurationclone.version
    //hostnames = ["*.example.net","example.com","m.example.com"]  
    host_names = ["rinaldi.sandbox.akamaideveloper.com","sujala.sandbox.akamaideveloper.com"]  
   // hostnames = ["rinaldi.sandbox.akamaideveloper.com"]  
    depends_on = ["akamai_appsec_configuration_clone.appsecconfigurationclone"]
}

/*
resource "akamai_appsec_security_policy_clone" "appsecsecuritypolicyclone" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version 
    
    create_from_security_policy = "LNPD_76189"
    policy_name = "PL Cloned Test for Launchpad"
    policy_prefix = "PL" 
    depends_on = ["akamai_appsec_configuration_clone.appsecconfigurationclone"]
   }

output "secpolicyclone" {
  value = akamai_appsec_security_policy_clone.appsecsecuritypolicyclone.policy_id
}
*/
/*
data "akamai_contract" "contract" {
}

data "akamai_group" "group" {
}
*/

resource "akamai_appsec_match_targets" "appsecmatchtargets" {
    //config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    //version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    config_id = akamai_appsec_configuration_clone.appsecconfigurationclone.config_id
    version = akamai_appsec_configuration_clone.appsecconfigurationclone.version
    type =  "website"
    //json =  file("${path.module}/match_targets.json")
    sequence =  1
    is_negative_path_match =  false
    is_negative_file_extension_match =  true
    default_file = "NO_MATCH" //"BASE_MATCH" //NO_MATCH
    host_names =  ["example.com","www.example.net","n.example.com"]
    file_paths =  ["/sssi/*","/cache/aaabbc*","/price_toy/*"]
    file_extensions = ["wmls","jpeg","pws","carb","pdf","js","hdml","cct","swf","pct"]
    security_policy = "f1rQ_106946"
    
    //bypass_network_lists = ["888518_ACDDCKERS","1304427_AAXXBBLIST"]
    depends_on = ["akamai_appsec_configuration_clone.appsecconfigurationclone"]
}

