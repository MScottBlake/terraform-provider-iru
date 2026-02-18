data "iru_tags" "example" {}

output "tag_list" {
  value = data.iru_tags.example.tags
}
