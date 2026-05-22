package main

type GenerateCmd struct {
	Chars  string `default:"a-zA-Z0-9!-/:-@[-"`
	ID     string `arg`
	Length int    `arg`
}

func (cmd *GenerateCmd) Run() error {
	return nil
}
