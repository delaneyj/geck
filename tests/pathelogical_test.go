package tests

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/delaneyj/geck"
	"github.com/dustin/go-humanize"
)

const (
	componentCount = 10
	queryCount     = 3
	sampleCount    = 1 << 12
	flipRate       = 0.5
	minEntityCount = 1 << 10
	maxEntityCount = 1 << 16
)

func report(queryName string, sum, min, max time.Duration, entityCount int) {
	sum /= sampleCount

	fmt.Printf("total:\t\tmin:\t%s,\tavg: %s,\tmax: %s (%s)\n", min, sum, max, queryName)

	ecf := float64(entityCount)
	minf := float64(min) / ecf
	sumf := float64(sum) / ecf
	maxf := float64(max) / ecf

	fmt.Printf("per entity:\tmin: %0.2fns,\tavg: %0.2fns,\tmax: %0.2fns\n\n", minf, sumf, maxf)
}

func flipCoin() bool {
	return rand.Float64() < flipRate
}

func BenchmarkPathelogical(b *testing.B) {

	for i := minEntityCount; i <= maxEntityCount; i *= 2 {
		run(b, i)
	}
}

func run(b *testing.B, entityCount int) {
	b.StopTimer()
	w := geck.NewWorld()

	start := time.Now()

	// Create component ids
	tags := w.CreateEntitiesWith(componentCount, nil)

	// record table count befor creating entities
	tableCount := w.ArchetypeCount()

	// Create entities
	entities := geck.NewIDSet()

	cIDsToInclude := geck.NewIDSet()
	for i := 0; i < entityCount; i++ {
		cIDsToInclude.Clear()

		tags.Range(func(tag geck.ID) {
			if flipCoin() {
				cIDsToInclude.Add(tag)
			}
		})
		eIDBits := w.CreateEntitiesWith(1, cIDsToInclude)
		entities.Or(eIDBits)
	}

	fmt.Printf(
		"upsert %s entities with %d components with flip rate of %0.0f%% in %s\n",
		humanize.Comma(int64(entityCount)),
		componentCount,
		flipRate*100,
		time.Since(start),
	)

	fmt.Printf("tables created   : %d\n", w.ArchetypeCount()-tableCount)
	fmt.Printf("querying for %d components taking %d samples\n", queryCount, sampleCount)

	var (
		entitySum        int
		sumD, minD, maxD time.Duration
	)

	// f, err := os.Create("profile.out")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	queryTags := make([]geck.ID, 0, componentCount)
	tags.ConditionalRange(func(id geck.ID, i int) bool {
		queryTags = append(queryTags, id)
		return i < queryCount
	})

	b.StartTimer()
	iter := w.QueryAnd(queryTags...)
	for s := 0; s < sampleCount; s++ {
		start = time.Now()
		iter.Reset()
		for iter.HasNext() {
			iter.Next()
			entitySum++
		}
		d := time.Since(start)
		sumD += d
		minD = min(minD, d)
		maxD = max(maxD, d)
	}

	if entitySum == 0 {
		fmt.Println("no entities found")
	} else {
		report("roaring golang", sumD, minD, maxD, entitySum)
	}

}
