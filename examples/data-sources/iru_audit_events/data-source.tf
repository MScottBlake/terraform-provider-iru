data "iru_audit_events" "recent" {}

output "recent_actions" {
  value = [for e in data.iru_audit_events.recent.results : e.action]
}
