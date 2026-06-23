data "iru_behavioral_detections_v2" "example" {
  statuses = ["informational", "detected"]
}

output "detection_descriptions" {
  value = [for r in data.iru_behavioral_detections_v2.example.results : r.description]
}
