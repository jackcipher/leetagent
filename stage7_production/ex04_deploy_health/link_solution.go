//go:build solution

package ex04_deploy_health

import "leetagent/stage7_production/ex04_deploy_health/solution"

const DockerfilePath = solution.DockerfilePath

type (
	Checker = solution.Checker
	Health  = solution.Health
	Env     = solution.Env
	Config  = solution.Config
)

var (
	LoadConfig      = solution.LoadConfig
	CheckDockerfile = solution.CheckDockerfile
)
