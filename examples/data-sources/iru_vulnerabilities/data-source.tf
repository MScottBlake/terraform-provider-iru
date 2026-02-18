data "iru_vulnerabilities" "critical" {}

output "vulnerability_ids" {
  value = [
    for v in data.iru_vulnerabilities.critical.results : v.cve_id 
    if v.severity == "CRITICAL"
  ]
}
