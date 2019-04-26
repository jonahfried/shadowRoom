package main

func (p *Agent) moveHandler(msg string) {
	switch msg {
	case "w":
		p.Acc.Y += 5
	case "s":
		p.Acc.Y -= 5
	case "a":
		p.Acc.X -= 5
	case "d":
		p.Acc.X += 5
	}
}
