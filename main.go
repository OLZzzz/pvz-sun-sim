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
var allZom seed
var checkSun = false // 检查阳光变化
var sunLimit = 9990
var survivedFlags int
var lowestSun []int

func main() {
	rand.Seed(time.Now().Unix())

	// 计算9990起手的破阵率，arg1是总模拟次数，arg2是目标冲关f数
	doSim(10000, 10000)

	// 计算对于不同阳光起手的破阵率（0到9990），arg1是总模拟次数，arg2是目标冲关f数
	//drawSim(1000, 100)
}

func doSim(repeat int, limit int) {
	f, _ := os.Create("out.txt")
	defer f.Close()
	fail := 0
	succ := 0

	// 进行repeat次模拟
	for i := 0; i < repeat; i++ {
		if i%(repeat/20) == 0 {
			fmt.Printf("No. %d\n", i)
		}
		if checkSun {
			io.WriteString(f, "Start: \n")
		}
		sun := sunLimit
		failed := false
		minSun := sunLimit
		for j := 0; j < limit/2; j++ {
			if sun < minSun {
				minSun = sun
			}
			sun += getSunChange()
			if sun > sunLimit {
				sun = sunLimit
			}
			if checkSun {
				io.WriteString(f, fmt.Sprintf("%d\t", sun))
			}
			if sun < 0 {
				failed = true
				survivedFlags += j
				break
			}
		}
		if checkSun {
			io.WriteString(f, "\n")
		}
		if failed {
			fail++
		} else {
			succ++
			survivedFlags += limit / 2
		}
		lowestSun = append(lowestSun, minSun)
	}

	// 输出结果
	fmt.Printf("Goal: Survive for %d flags\n", limit)
	fmt.Printf("Failure: %d / %d (%f)\n", fail, repeat, float64(fail)/float64(repeat))
	fmt.Printf("Success: %d / %d (%f)\n", succ, repeat, float64(succ)/float64(repeat))
	fmt.Printf("Average survived flags: %.1ff\n", 2*float64(survivedFlags)/float64(repeat))
	output(lowestSun, repeat)

	// 显示出怪详细
	if checkZomGen {
		for i := range allZom {
			fmt.Printf("%.6f ", float64(allZom[i])/float64(survivedFlags))
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
		for i := range seed {
			allZom[i] += seed[i]
		}
	}

	// 以下为多元线性回归求得的阳光变化公式

	// 5.5花两仪，五花两仪+岸路底线向日葵/阳光菇
	return int(2466 + 1246/2 - 403*seed[2] - 789*seed[5] - 442*seed[7] - 604*seed[11] - 312*seed[12] - 491*seed[17] - 2478*seed[18])

	// 6花两仪，来自5花两仪反推
	// return int(3711  - 403*seed[2] - 789*seed[5] - 442*seed[7] - 604*seed[11] - 625*seed[12] - 491*seed[17] - 2478*seed[18])

	// 6花两仪
	// return int(2106 - 237*seed[2] - 1083*seed[5] - 255*seed[12] - 769*seed[17] - 2597*seed[18])
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

func output(minSun []int, repeat int) {
	//fmt.Println(minSun)
	min := 9990
	sum := 0
	for _, j := range minSun {
		if j < min {
			min = j
		}
		sum += j
	}
	fmt.Printf("Average lowest sun: %.2f (min: %d)", float64(sum)/float64(repeat), min)
}
