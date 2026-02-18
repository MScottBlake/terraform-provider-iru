data "iru_behavioral_detections" "example" {}

output "detections" {
  value = data.iru_behavioral_detections.example.results
}
