package utils

import "fmt"

type Build uint

const (
  DevBuild Build = iota+1
  Releasebuild
)

func(build Build)String()string{
  return [...]string{
    "Development",
    "Release",
  }[build-1]
}

func Mode(build string) Build {
  switch build {
  case "Development":
    return DevBuild
  case "Release":
    return Releasebuild
  default:
    panic(fmt.Sprintf("Unknown Build Mode Configured: \"%s\"", build))
  }
}
