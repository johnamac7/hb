workflow "Build" {
  on       = "release"

  resolves = [
    "release darwin/amd64",
    "release windows/amd64",
    "release linux/amd64",
  ]
}

action "release darwin/amd64" {
  uses    = "ngs/go-release.action@v1.0.2"

  secrets = [
    "GITHUB_TOKEN"
  ]

  env     = {
    GOOS   = "darwin"
    GOARCH = "amd64"
  }
}

action "release windows/amd64" {
  uses    = "ngs/go-release.action@v1.0.2"

  secrets = [
    "GITHUB_TOKEN"
  ]

  env     = {
    GOOS   = "windows"
    GOARCH = "amd64"
  }
}

action "release linux/amd64" {
  uses    = "ngs/go-release.action@v1.0.2"

  secrets = [
    "GITHUB_TOKEN"
  ]

  env     = {
    GOOS   = "linux"
    GOARCH = "amd64"
  }
}
