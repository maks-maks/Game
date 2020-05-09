package main

import (
	"fmt"
	"math/rand"
	"time"
)

func old_game_main() {
	rand.Seed(int64(time.Now().Nanosecond()))

	// squad1 := Squad{unit1, unit2, unit3}
	squad1 := createSquad_("SANDALI")
	// fmt.Println(squad1)
	squad2 := createSquad_("GEROI")
	// for squad1.units[2].health > 0 && squad2.units[2].health > 0 {
	// 	time.Sleep(100 * time.Millisecond)
	// 	battleSquad(&squad1, &squad2)
	// 	battleSquad(&squad2, &squad1)
	// 	fmt.Println(squad1)
	// 	fmt.Println(squad2)
	// }

	i := 0
	d := 100 * time.Millisecond
	for true {
		time.Sleep(d)
		// fmt.Printf("\rtick %d ", i)
		i = i + 1

		battleSquad(d, &squad1, &squad2)
		battleSquad(d, &squad2, &squad1)
		fmt.Printf("\r %d %v %v", i, squad1, squad2)
		// fmt.Println(squad2)

		if squad1.units[2].health <= 0 || squad2.units[2].health <= 0 {
			break
		}
	}
}

func battleSquad(d time.Duration, a, b *Squad) {
	for _, unit := range a.units {
		if unit.health <= 0 {
			continue // can't hit because dead
		}

		if unit.damage > 0 {
			target := chooseTarget(*b)
			if target == nil {
				return
			}
			battle(d, unit, target)
		} else if unit.heal > 0 {
			target := chooseTarget(*a)
			if target == nil {
				return
			}
			heal(d, unit, target)
		}
	}
}

func chooseTarget(s Squad) *Unit {
	for i := 0; i < len(s.units); i++ {
		if s.units[i].health > 0 {
			return s.units[i]
		}
	}
	return nil
}

type Squad struct {
	name  string
	units []*Unit
}

func createSquad_(n string) Squad {
	squad := Squad{
		name: n,
	}

	tank := &Unit{
		name:        "Frederik",
		squad:       &squad,
		health:      500,
		maxHealth:   500,
		damage:      500,
		cooldown:    6000000 * time.Millisecond,
		stamina:     10,
		staminaCost: 90,
		dodge:       10,
	}
	ranger := &Unit{
		name:        "Robin",
		squad:       &squad,
		health:      200,
		maxHealth:   200,
		damage:      5,
		cooldown:    100 * time.Millisecond,
		stamina:     100,
		staminaCost: 10,
		dodge:       30,
	}
	healer := &Unit{
		name:        "Angel",
		squad:       &squad,
		health:      100,
		maxHealth:   100,
		heal:        50,
		cooldown:    4000 * time.Millisecond,
		stamina:     100,
		staminaCost: 60,
	}

	squad.units = append(squad.units, tank, ranger, healer)

	return squad
}

type Unit struct {
	name        string
	squad       *Squad
	health      int
	maxHealth   int
	damage      int
	cooldown    time.Duration
	reload      time.Duration
	stamina     int
	staminaCost int
	dodge       int
	heal        int
}

func heal(d time.Duration, a *Unit, b *Unit) {
	a.reload = a.reload - d
	if a.reload > 0 {
		return
	}
	a.reload = a.cooldown
	if a.stamina < a.staminaCost {
		a.stamina = a.stamina + 50

		return
	}
	if b.health < b.maxHealth {
		a.stamina = a.stamina - a.staminaCost
		b.health = b.health + a.heal
		if b.health > b.maxHealth {
			b.health = b.maxHealth
		}
		fmt.Printf(" %s(%s) HEAL %s(%s) by %d\n", a.name, a.squad.name, b.name, b.squad.name, a.heal)
	}
}
func battle(d time.Duration, a *Unit, b *Unit) {
	a.reload = a.reload - d
	if a.reload > 0 {
		return
	}
	a.reload = a.cooldown

	if a.stamina < a.staminaCost {
		a.stamina = a.stamina + 50

		return
	}

	a.stamina = a.stamina - a.staminaCost

	randPercent := rand.Intn(100)
	if b.dodge > randPercent {
		// fmt.Printf("!!!   %s(%s) dodge from %s(%s) hit\n", b.name, b.squad.name, a.name, a.squad.name)
		return
	}
	b.health = b.health - a.damage
	// fmt.Printf(" %s(%s) HIT %s(%s) by %d\n", a.name, a.squad.name, b.name, b.squad.name, a.damage)
}

func (u Unit) String() string {
	return fmt.Sprintf("{%s, %d, %d}", u.name, u.health, u.stamina)
}

func (s Squad) String() string {
	str := s.name + " "
	for _, unit := range s.units {
		str += unit.String()
	}
	return str
}
