
resource "aws_s3_bucket_public_access_block" "prowler-poc-bucket_access_block" {
  bucket = "prowler-poc-bucket"

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls    = true
  restrict_public_buckets = true
}
