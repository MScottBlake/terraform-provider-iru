data "iru_threats_v2" "example" {
  statuses         = ["quarantined", "released"]
  severities       = ["critical", "high"]
  management_state = "managed"
}

output "threat_names" {
  value = [for r in data.iru_threats_v2.example.results : r.threat_name]
}
