package main

import (
	"arena"
	"fmt"
	"math"
	"runtime"
	"testing"
)

type S struct {
	byte
}

func BenchmarkLargeSliceWithoutArena(b *testing.B) {
	for i := 0; i <= 3; i++ { // 1 to 1000 MiB
		n := int(math.Pow(10, float64(i)))
		b.Run(fmt.Sprintf("%dMiB", n), func(b *testing.B) {
			// 2^17 * i * 100 * 8 bytes = i * 100 MiB
			l := 1 << 17 * n
			for j := 0; j < b.N; j++ {
				b.StopTimer()
				func() {
					s := make([]*S, l, l)

					b.StartTimer()
					runtime.GC()
					b.StopTimer()

					b.ReportAllocs()

					runtime.KeepAlive(s)
				}()
				runtime.GC()
			}
			fmt.Printf("--- %dMiB\n\n", n)
		})
	}
}

func BenchmarkLargeSliceWithArena(b *testing.B) {
	for i := 0; i <= 3; i++ { // 1 to 1000 MiB
		n := int(math.Pow(10, float64(i)))
		b.Run(fmt.Sprintf("%d0MiB", n), func(b *testing.B) {
			// 2^17 * i * 100 * 8 bytes = i * 100 MiB
			l := 1 << 17 * n
			for j := 0; j < b.N; j++ {
				b.StopTimer()
				func() {
					a := arena.NewArena()
					s := arena.MakeSlice[*S](a, l, l)

					b.StartTimer()
					runtime.GC()
					b.StopTimer()

					b.ReportAllocs()

					runtime.KeepAlive(s)

					a.Free()
				}()
				runtime.GC()
			}

			fmt.Printf("------- %dMiB\n\n", n)
		})
	}
}

func BenchmarkSliceHasManyItemsWithoutArena(b *testing.B) {
	for i := 1; i <= 6; i++ { // 10 to 1,000,000 items
		n := int(math.Pow(10, float64(i)))
		b.Run(fmt.Sprintf("%d items", n), func(b *testing.B) {
			// fmt.Printf("\n\n------- b.N: %d\n", b.N)
			for j := 0; j < b.N; j++ {
				b.StopTimer()
				func() {
					s := make([]*int, n)
					for i := 0; i < n; i++ {
						s[i] = new(int)
						runtime.KeepAlive(*s[i])
					}

					b.StartTimer()
					runtime.GC()
					b.StopTimer()

					b.ReportAllocs()

					runtime.KeepAlive(s)
				}()
				runtime.GC()
			}

			// fmt.Printf("%d items, NumEarlyReturnForUserArena: %d\n", n, runtime.NumEarlyReturnForUserArena())
		})
	}
}

func BenchmarkSliceHasManyItems(b *testing.B) {
	for i := 1; i <= 6; i++ { // 10 to 1,000,000 items
		n := int(math.Pow(10, float64(i)))
		b.Run(fmt.Sprintf("%d items", n), func(b *testing.B) {
			// fmt.Printf("\n\n------- b.N: %d\n", b.N)
			for j := 0; j < b.N; j++ {
				b.StopTimer()
				func() {
					a := arena.NewArena()
					s := make([]*int, n)
					for i := 0; i < n; i++ {
						s[i] = arena.New[int](a)
						runtime.KeepAlive(*s[i])
					}

					b.StartTimer()
					runtime.GC()
					b.StopTimer()

					b.ReportAllocs()
				}()
				runtime.GC()
			}
			// fmt.Printf("%d items, NumEarlyReturnForUserArena: %d\n", n, runtime.NumEarlyReturnForUserArena())
		})
	}
}

// ryicoh@ryicohs-MacBook-Air src % GOEXPERIMENT=arenas GOMAXPROCS=1 ../bin/go test -bench BenchmarkSliceHasManyItems ./arenabench -benchtime=1s
// goos: darwin
// goarch: arm64
// pkg: arenabench
// BenchmarkSliceHasManyItemsWithoutArena/10_items                    12436             99104 ns/op               0 B/op          0 allocs/op
// BenchmarkSliceHasManyItemsWithoutArena/100_items                   10000            115794 ns/op               0 B/op          0 allocs/op
// BenchmarkSliceHasManyItemsWithoutArena/1000_items                  10000            106620 ns/op               0 B/op          0 allocs/op
// BenchmarkSliceHasManyItemsWithoutArena/10000_items                  6811            175305 ns/op               0 B/op          0 allocs/op
// BenchmarkSliceHasManyItemsWithoutArena/100000_items                 1390            862616 ns/op               0 B/op          0 allocs/op
// BenchmarkSliceHasManyItemsWithoutArena/1000000_items                  85          12900849 ns/op               0 B/op          0 allocs/op
// BenchmarkSliceHasManyItems/10_items                                12380             96816 ns/op               0 B/op          0 allocs/op
// BenchmarkSliceHasManyItems/100_items                               12403             96928 ns/op               0 B/op          0 allocs/op
// BenchmarkSliceHasManyItems/1000_items                              12211             96750 ns/op               0 B/op          0 allocs/op
// BenchmarkSliceHasManyItems/10000_items                             12109             99193 ns/op               0 B/op          0 allocs/op
// BenchmarkSliceHasManyItems/100000_items                            10000            111163 ns/op               1 B/op          0 allocs/op
// BenchmarkSliceHasManyItems/1000000_items                            3920            322053 ns/op              14 B/op          0 allocs/op
// PASS
// ok      arenabench      95.435s
