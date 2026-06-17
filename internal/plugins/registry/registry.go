package registry

import (
	base "mw-sre-platform/internal/plugins"
	"mw-sre-platform/internal/plugins/opengauss"
	"mw-sre-platform/internal/plugins/pythonenv"
	"mw-sre-platform/internal/plugins/redis"
)

func Default() []base.Plugin {
	return []base.Plugin{
		redis.Plugin{},
		opengauss.Plugin{},
		pythonenv.Plugin{},
	}
}

func Names() []string {
	ps := Default()
	out := make([]string, 0, len(ps))
	for _, p := range ps {
		out = append(out, p.Name())
	}
	return out
}
