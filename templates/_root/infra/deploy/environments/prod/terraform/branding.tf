locals {
  company_name     = "Dinner Done Better"
  company_slug     = "dinner-done-better"
  company_slug_ns  = "dinnerdonebetter"
  public_domain    = "${local.company_slug_ns}.com"
  gcp_project_id   = "${local.company_slug}-prod"
}
