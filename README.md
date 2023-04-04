# Tools for manage yandex cloud

![example workflow](https://github.com/tsimer123/pet-infra-yandex-tools/actions/workflows/main.yml/badge.svg)

## iam-token

Get IAM token with service account key. Yandex documentation for key creating: https://cloud.yandex.ru/docs/iam/operations/authorized-key/create.

Environments:

- `YA_JWT_KEY_ID` - service account key id.
- `YA_JWT_SERVICE_ACCOUNT_ID` - service account id.
- `YA_JWT_KEY_BASE64` - base64 encoded service account key.
- `GITHUB_TOKEN` - github token for repository access.
- `GITHUB_OWNER` - repository owner.
- `GITHUB_REPO` - repository name.
- `GITHUB_SECRET_NAME` - repository secret name.
