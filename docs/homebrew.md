# Homebrew Setup

The homebrew tap is located at https://github.com/thegreenwebfoundation/homebrew-carbon-aware-tools

## Permissions

The GoReleaser GitHub Action uses the secret `HOMEBREW_TAP_GITHUB_TOKEN` to
authenticate with the GitHub API when publishing to the 
`homebrew-carbon-aware-tools` repo.

To restrict access to just this repo a GitHub App has been created with an
installation token that can write to this repo.

## Setup

- The GitHub App is named `Grid Intensity App O Tron` and can be found in
the settings of the `thegreenwebfoundation` organization.

![GitHub App]](github_app.png)

- By default installation tokens expire after 8 hours. So we need to opt out of
this setting via Optional Features.

![GitHub App optional features]](github_app_optional_features.png)

- To create the installation token we need to create a [private key](https://docs.github.com/en/developers/apps/building-github-apps/authenticating-with-github-apps)
for the app. (This can be found in 1 Password).

![GitHub App private key]](github_app_private_key.png)

- The app needs to be installed in the `thegreenwebfoundation` org for just the
`grid-intensity-go` repository.

![GitHub App installation]](github_app_install.png)
![GitHub App installation for repo]](github_app_install_repo.png)

- Create a JWT token using the linked Ruby script with the private key.
https://docs.github.com/en/developers/apps/building-github-apps/authenticating-with-github-apps#authenticating-as-a-github-app

- Get the installation ID for the app for the grid-intensity-go repo.

```
 curl \
  -H "Accept: application/vnd.github.v3+json" \
  -H "Authorization: Bearer ***JWT TOKEN***" \
  -s https://api.github.com/app/installations | jq '.[0].id'

26614624
```

- Create an installation token that can write to `grid-intensity-go`.

```
curl \
  -X POST \
  -H "Accept: application/vnd.github.v3+json" \
  -H "Authorization: Bearer ***JWT TOKEN***
  https://api.github.com/app/installations/26614624/access_tokens \
  -d '{"repository":"homebrew-carbon-aware-tools","permissions":{"contents":"write"}}'
```

- Add the token as a secret named `HOMEBREW_TAP_GITHUB_TOKEN` in the
`grid-intensity-go` repo.

![GitHub secret with token]](github_repo_secret.png)
