package main

//ManagerBots основная структура менеджера ботов
type ManagerBots struct {
	ListBot map[string]*Bot
}

func newManagerBots() (*ManagerBots, error) {
	mb := &ManagerBots{}
	mb.ListBot = make(map[string]*Bot)
	return mb, nil
}
