package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

type seed [19]int

var checkZomGen = false // 检查出怪分布是否正确
var totalRun int
var allZom seed

func main() {
	rand.Seed(time.Now().Unix())

	// 计算9990起手的破阵率，arg1是总模拟次数，arg2是目标冲关f数
	// doSim(10000, 10000)

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

	// 显示出怪详细
	if checkZomGen {
		for i := range allZom {
			fmt.Printf("%.6f ", float64(allZom[i])/float64(totalRun))
		}
	}

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

	// 记录此次出怪
	if checkZomGen {
		totalRun++
		for i := range seed {
			allZom[i] += seed[i]
		}
	}

	// 以下为多元线性回归求得的阳光变化公式
	return int(2466 - 403*seed[2] - 789*seed[5] - 442*seed[7] - 604*seed[11] - 491*seed[17] - 2478*seed[18])
}

// 随机决定出怪
func rollSeed() seed {
	var seed seed

	// 决定选路障或读报
	flag := rand.Intn(5)
	if flag == 0 {
		seed[3] = 1
	} else {
		seed[0] = 1
	}

	// 通过对僵尸列表shuffle，随机选出剩下僵尸
	// 路障 撑杆 铁桶 读报 铁门 橄榄 舞王 潜水 冰车 海豚 小丑 气球 矿工 跳跳 蹦极 扶梯 投篮 白眼 红眼 旗帜 僵王
	// 其中，旗帜/僵王即使被选中，最终也会被忽略
	zombies := make([]int, 0)
	for i := 0; i < 21; i++ {
		if flag == 0 && i == 3 {
			continue
		}
		if flag != 0 && i == 0 {
			continue
		}
		zombies = append(zombies, i)
	}
	rand.Shuffle(len(zombies), func(i, j int) { zombies[i], zombies[j] = zombies[j], zombies[i] })
	for i := 0; i < 9; i++ {
		if zombies[i] >= 19 {
			continue
		}
		seed[zombies[i]] = 1
	}
	return seed
}
