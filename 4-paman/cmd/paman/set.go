package main

type SetCmd struct {
	ID string `arg`
}

func (cmd *SetCmd) Run() error {
	return nil
}
