package cube

type Config struct {
	Procedure func(proc *Procedure) error
	Filename  string
	Source    string
}

func Compile(config *Config) error {
	return (&parseContext{
		config: config,
		lexer:  NewLexer(config.Source),
	}).parse()
}
