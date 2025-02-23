package config

type ConfigOption func(l *Loader)

// 设置Host字段
func WithHost(host string) ConfigOption {
	return func(l *Loader) {
		l.Host = host
	}
}

// 设置Paths字段
func WithPaths(paths []string) ConfigOption {
	return func(l *Loader) {
		l.Paths = append(l.Paths, paths...)
	}
}

// EnableEnv字段
func WithEnableEnv(enable bool) ConfigOption {
	return func(l *Loader) {
		l.EnableEnv = enable
	}
}

// 设置EnvPrefix字段
func WithEnvPrefix(prefix string) ConfigOption {
	return func(l *Loader) {
		l.EnvPrefix = prefix
	}
}
