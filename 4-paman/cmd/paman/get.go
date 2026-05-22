package main

type GetCmd struct {
	ID string `arg`
}

func (cmd *GetCmd) Run() error {
	return nil
}
