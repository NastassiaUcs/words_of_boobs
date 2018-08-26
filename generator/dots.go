package generator

import "math/rand"

type dotsManager struct {
	dots map[int]map[int]bool
	count int
}

type point struct {
	x, y int
}

func (self *dotsManager) addDot(x int, y int) {
	if _, ok := self.dots[x]; !ok {
		self.dots[x] = make(map[int]bool)
	}
	if _, v := self.dots[x][y]; !v {
		self.dots[x][y] = true
		self.count++
	}
}

func (self *dotsManager) removeDot(x int, y int) {
	if _, ok := self.dots[x]; ok {
		if self.dots[x][y] {
			self.dots[x][y] = false
			self.count--
		}
	}
}

func (self *dotsManager) checkDot(x int, y int) bool {
	if _, ok := self.dots[x]; ok {
		return self.dots[x][y]
	}
	return false
}

func (self *dotsManager) getList(shuffle bool) []point {
	points := make([]point, self.count)
	count := 0
	for x, column := range self.dots {
		for y, v := range column {
			if v {
				points[count] = point{x: x, y: y}
				count++
			}
		}
	}
	if shuffle {
		rand.Shuffle(count, func(i, j int) {
			points[i], points[j] = points[j], points[i]
		})
	}
	return points
}

func (self *dotsManager) getRandomDot() point {
	list := self.getList(false)
	return list[rand.Intn(self.count)]
}

func createDots() *dotsManager {
	var d = dotsManager{}
	d.dots = make(map[int]map[int]bool)
	d.count = 0
	return &d
}