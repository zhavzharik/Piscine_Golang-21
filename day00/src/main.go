package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	_ "math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func ReadInput() []int {
	array := make([]int, 0)
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		str := scanner.Text()
		str = strings.TrimSpace(str)
		num, err := strconv.Atoi(str)
		if err != nil || num < -100000 || num > 100000 {
			fmt.Println("Wrong argument(s)!")

			continue
		}
		array = append(array, num)
	}
	return array
}

func CalcMean(array []int) float64 {
	var total float64 = 0
	for _, value := range array {
		total += float64(value)
	}
	mean := total / float64(len(array))

	return mean
}

func CalcMedian(array []int) float64 {
	var median float64 = 0
	if len(array) > 0 {
		mid := len(array) / 2

		if len(array)%2 != 0 {
			median += float64(array[mid])
		} else {
			median += float64(array[mid]+array[mid-1]) / float64(2)
		}
		return median
	}
	return math.NaN()
}

func CalcMode(array []int) (int, float64) {
	if len(array) > 0 {
		countMap := make(map[int]int)
		for _, value := range array {
			countMap[value] += 1
		}
		max := 0
		mode := 0
		check := 42.42
		for _, key := range array {
			freq := countMap[key]
			if freq > max {
				mode = key
				max = freq
			}
		}
		return mode, check
	}
	return 0, math.NaN()
}

func CalcStd(array []int) float64 {
	if len(array) > 0 {
		var sd float64 = 0
		mean := CalcMean(array)
		var sumSquare float64 = 0
		var disp float64 = 0
		for _, value := range array {
			sumSquare += (float64(value) - mean) * (float64(value) - mean)
		}
		disp += sumSquare / float64(len(array)-1)
		sd += math.Sqrt(disp)
		return sd
	}
	return math.NaN()
}

func main() {
	useMean := flag.Bool("Mean", false, "display Mean")
	useMedian := flag.Bool("Median", false, "display Median")
	useMode := flag.Bool("Mode", false, "display Mode")
	useSd := flag.Bool("SD", false, "display SD")

	flag.Parse()
	array := ReadInput()
	sort.Ints(array)
	if *useMean {
		mean := CalcMean(array)
		fmt.Printf("Mean: %.2f\n", mean)
	}
	if *useMedian {
		median := CalcMedian(array)
		fmt.Printf("Median: %.2f\n", median)
	}
	if *useMode {
		mode, check := CalcMode(array)
		if math.IsNaN(check) {
			fmt.Printf("Mode: %f\n", check)
		} else {
			fmt.Printf("Mode: %d\n", mode)
		}
	}
	if *useSd {
		sd := CalcStd(array)
		fmt.Printf("SD: %.2f\n", sd)
	}
	if !*useMean && !*useMode && !*useMedian && !*useSd {
		mean := CalcMean(array)
		fmt.Printf("Mean: %.2f\n", mean)

		median := CalcMedian(array)
		fmt.Printf("Median: %.2f\n", median)

		mode, check := CalcMode(array)
		if math.IsNaN(check) {
			fmt.Printf("Mode: %f\n", check)
		} else {
			fmt.Printf("Mode: %d\n", mode)
		}

		sd := CalcStd(array)
		fmt.Printf("SD: %.2f\n", sd)
		return
	}
}
