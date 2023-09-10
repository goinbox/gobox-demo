package demo

type Controller struct {
}

func (c *Controller) Name() string {
	return "Demo"
}

func (c *Controller) IndexAction() *indexAction {
	return newIndexAction()
}

func (c *Controller) AddAction() *addAction {
	return newAddAction()
}

func (c *Controller) EditAction() *editAction {
	return newEditAction()
}

func (c *Controller) DelAction() *delAction {
	return newDelAction()
}
