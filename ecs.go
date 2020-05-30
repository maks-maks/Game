package main

import (
	"fmt"
	"reflect"

	"github.com/bytearena/ecs"
)

var ecsManager *ECSManager

type ECSManager struct {
	ecs.Manager
	events       EventBus
	components   []*ecs.Component
	componentMap map[string]*ecs.Component
	nameMap      map[*ecs.Component]string
	type2comp    map[reflect.Type]*ecs.Component
}

func NewECSManager() *ECSManager {
	return &ECSManager{
		Manager:      *(ecs.NewManager()),
		componentMap: make(map[string]*ecs.Component),
		nameMap:      make(map[*ecs.Component]string),
		type2comp:    make(map[reflect.Type]*ecs.Component),
	}
}

func (m *ECSManager) RegisterComponent(name string, componentDataSample interface{}) *ecs.Component {
	c := m.NewComponent()
	m.components = append(m.components, c)
	m.componentMap[name] = c
	m.nameMap[c] = name

	t := reflect.TypeOf(componentDataSample)
	if t.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("only pointers could be registered as component data in ECSManager (got %s)", t))
	}
	m.type2comp[t] = c

	return c
}

func (m *ECSManager) ComponentName(c *ecs.Component) string {
	return m.nameMap[c]
}

func (m *ECSManager) AllComponentDataForEntity(e *ecs.Entity) map[*ecs.Component]interface{} {
	r := make(map[*ecs.Component]interface{})
	for _, v := range m.componentMap {
		d, ok := e.GetComponentData(v)
		if ok {
			r[v] = d
		}
	}
	return r
}

func (m *ECSManager) AddComponent(entity *ecs.Entity, componentData interface{}) error {
	component, err := m.componentByData(componentData)
	if err != nil {
		return err
	}
	entity.AddComponent(component, componentData)
	return nil
}

func (m *ECSManager) componentByData(componentData interface{}) (*ecs.Component, error) {
	c, ok := m.type2comp[reflect.TypeOf(componentData)]
	if !ok {
		panic(fmt.Sprintf("no component registered for type `%s` in ECSManager", reflect.TypeOf(componentData)))
	}
	return c, nil
}

type System interface {
	Update(dt float32)
}
