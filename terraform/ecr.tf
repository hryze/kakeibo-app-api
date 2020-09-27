resource "aws_ecr_repository" "kakeibo_ecr_repository" {
  name                 = "kakeibo-ecr-repository"
  image_tag_mutability = "MUTABLE"
  tags                 = merge(local.default_tags, map("Name", "kakeibo-ecr-repository"))

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_lifecycle_policy" "kakeibo_lifecycle_policy" {
  repository = aws_ecr_repository.kakeibo_ecr_repository.name

  policy = <<EOF
{
    "rules": [
        {
            "rulePriority": 1,
            "description": "Keep last 30 images",
            "selection": {
                "tagStatus": "tagged",
                "tagPrefixList": ["v"],
                "countType": "imageCountMoreThan",
                "countNumber": 30
            },
            "action": {
                "type": "expire"
            }
        }
    ]
}
EOF
}

resource "aws_ecr_repository_policy" "kakeibo_repository_policy" {
  repository = aws_ecr_repository.kakeibo_ecr_repository.name
  policy     = data.aws_iam_policy_document.kakeibo_ecr_policy_document.json
}

data "aws_iam_policy_document" "kakeibo_ecr_policy_document" {
  statement {
    sid = "ecr policy"
    principals {
      type        = "*"
      identifiers = ["*"]
    }
    effect = "Allow"
    actions = [
      "ecr:GetAuthorizationToken",
      "ecr:BatchGetImage",
      "ecr:BatchCheckLayerAvailability",
      "ecr:GetDownloadUrlForLayer",
      "ecr:InitiateLayerUpload",
      "ecr:CompleteLayerUpload",
      "ecr:UploadLayerPart",
      "ecr:DescribeImages",
      "ecr:PutImage",
      "ecr:DescribeRepositories",
      "ecr:GetRepositoryPolicy",
      "ecr:ListImages",
      "ecr:BatchDeleteImage",
    ]
  }
}
