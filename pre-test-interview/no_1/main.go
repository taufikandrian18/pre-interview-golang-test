package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// struct untuk menampung data dari setiap worker
type WorkerResult struct {
	WorkerID int
	Sum      int64
	Count    int
}

func main() {
	// konstanta untuk konfigurasi slice dan jumlah worker
	const (
		numIntegers = 1000000 // ukuran slice bisa di rubah sewaktu waktu
		numWorkers  = 4       // jumlah worker
	)

	// membuat fungsi random
	rand.Seed(time.Now().UnixNano())
	// membuat slice numbers dengan ukuran besar dari num integers
	numbers := make([]int, numIntegers)
	for i := 0; i < numIntegers; i++ {
		numbers[i] = rand.Intn(1000000)
		//membuat list slice numbers untuk nantinya akan di sum untuk setiap nomor genap dari list tersebut
	}

	// memanggil fungsi utama untuk membuat concurrency dari setiap wroker yang menghitung sum nomor genap
	sumEvenNumbersConcurrently(numbers, numWorkers)
}

// fungsi utama untuk menghitung sum nomor genap
func calculateEvenSum(numbers []int, workerID int, results chan<- WorkerResult, wg *sync.WaitGroup) {
	// di eksekusi terakhir setelah semua kode di dalam fungsi di eksekusi, yang berguna untuk mengakhiri setiap wait group
	defer wg.Done()

	var sum int64
	var count int

	// melakukan iterasi dari setiap list numbers
	for _, num := range numbers {
		// memilah nomor genap dari setiap list numbers
		if num%2 == 0 {
			sum += int64(num)
			count++
		}
	}

	// menyalurkan hasil dari penghitungan ke worker result
	results <- WorkerResult{
		WorkerID: workerID,
		Sum:      sum,
		Count:    count,
	}
}

// fungsi utama untuk membuat concurrency dari setiap wroker yang menghitung sum nomor genap
func sumEvenNumbersConcurrently(numbers []int, numWorkers int) {
	// inisialisasi channel dengan tipe data struct worker dengan size numWorkers
	results := make(chan WorkerResult, numWorkers)
	// inisialisasi wait group untuk mengetahui kapan semua worker selesai bekerja
	var wg sync.WaitGroup

	// menghitung ukuran chunk untuk setiap worker
	chunkSize := len(numbers) / numWorkers
	// menghitung sisa untuk setiap worker
	sisaSize := len(numbers) % numWorkers

	index := 0
	for i := 0; i < numWorkers; i++ {
		lastIndex := index + chunkSize

		// mendistribusikan sisa ke setiap worker
		if i < sisaSize {
			lastIndex++
		}

		// menambahkan wait group
		wg.Add(1)

		//membuat go routine untuk mengkalkulasi sum nomor genap dari setiap worker
		go calculateEvenSum(numbers[index:lastIndex], i+1, results, &wg)

		index = lastIndex
	}

	// menunggu semua worker selesai dan menutup channel
	go func() {
		wg.Wait()
		close(results)
	}()

	//meng collect hasil dari setiap worker
	var totalSum int64
	fmt.Println("\n=== Worker Results ===")
	for result := range results {
		fmt.Printf("Worker %d: Sum = %d, Even numbers found = %d\n",
			result.WorkerID, result.Sum, result.Count)
		totalSum += result.Sum
	}

	fmt.Println("\n=== Final Results ===")
	fmt.Printf("Total sum of even numbers: %d\n", totalSum)
	fmt.Printf("Total numbers processed: %d\n", len(numbers))
}
