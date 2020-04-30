package main

func setupECS() {
	ecsManager = NewECSManager()

	ecsManager.RegisterComponent("position", &PositionComponent{})

	e1 := ecsManager.NewEntity()
	ecsManager.AddComponent(e1, &PositionComponent{X: 1, Y: 2})

	e2 := ecsManager.NewEntity()
	ecsManager.AddComponent(e2, &PositionComponent{X: 10, Y: 15})
}
