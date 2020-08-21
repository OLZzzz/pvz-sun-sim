package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

// 出怪总数为9种、10种、11种的概率分布
var totalNumProb = []float64{0.1895, 0.5211, 0.2895}

func main() {
	rand.Seed(time.Now().Unix())

	// 计算9990起手的破阵率，arg1是总模拟次数，arg2是目标冲关f数
	doSim(10000, 100)

	// 计算对于不同阳光起手的破阵率（0到9990），arg1是总模拟次数，arg2是目标冲关f数
	drawSim(1000, 100)
}

func doSim(repeat int, limit int) {
	fail := 0
	succ := 0
	for i := 0; i < repeat; i++ {
		sun := 9990
		failed := false
		for j := 0; j < limit/2; j++ {
			sun += getSunChange()
			if sun > 9990 {
				sun = 9990
			}
			if sun < 0 {
				failed = true
				break
			}
		}
		if failed {
			fail++
		} else {
			succ++
		}
	}
	fmt.Printf("Failure: %d / %d (%f)\n", fail, repeat, float64(fail)/float64(repeat))
	fmt.Printf("Success: %d / %d (%f)\n", succ, repeat, float64(succ)/float64(repeat))
}

func drawSim(repeat int, limit int) {
	f, err := os.Create("out.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	for i := 0; i <= 10000; i += 25 {
		fail := 0
		startSun := i
		if i == 10000 {
			startSun = 9990
		}
		for j := 0; j < repeat; j++ {
			sun := startSun
			for j := 0; j < limit/2; j++ {
				sun += getSunChange()
				if sun > 9990 {
					sun = 9990
				}
				if sun < 0 {
					fail++
					break
				}
			}
		}
		io.WriteString(f, fmt.Sprintf("%d\t%.4f\n", startSun, float64(fail)/float64(repeat)))
	}
}

// 随机出阳光变化
func getSunChange() int {
	seed := rollSeed()

	// 以下为多元线性回归求得的阳光变化公式
	return int(2695 - 608*seed[2] - 886*seed[5] - 580*seed[7] - 582*seed[11] - 536*seed[17] - 2785*seed[18])
}

// 随机决定出怪
func rollSeed() [19]int {
	var seed [19]int

	// 决定出怪总数
	totalNum := rollTotalNum()

	// 决定选路障或读报
	flag := roll(0.8) // 若为true，选择路障；若为false，选择读报
	if flag {
		seed[0] = 1
	} else {
		seed[3] = 1
	}

	// 通过对僵尸列表shuffle，随机选出剩下僵尸
	zombies := make([]int, 0)
	for i := 0; i < 19; i++ {
		if flag && i == 0 {
			continue
		}
		if !flag && i == 3 {
			continue
		}
		zombies = append(zombies, i)
	}
	rand.Shuffle(len(zombies), func(i, j int) { zombies[i], zombies[j] = zombies[j], zombies[i] })
	for i := 0; i < totalNum-1; i++ {
		seed[zombies[i]] = 1
	}
	return seed
}

// 随机决定出怪总数
func rollTotalNum() int {
	flag := roll(totalNumProb[0])
	if flag {
		return 8
	}
	flag = roll(totalNumProb[1] / (totalNumProb[1] + totalNumProb[2]))
	if flag {
		return 9
	}
	return 10
}

// 对于概率为prob的事件，随机决定此事件是否发生
func roll(prob float64) bool {
	f := 1 - rand.Float64() // 0.0 < f <= 1.0
	if prob >= f {
		return true
	}
	return false
}
