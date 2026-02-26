package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"matching-engine-lab/go/internal/core"
	"matching-engine-lab/go/internal/enginev1_btree"
	"matching-engine-lab/go/internal/enginev2_treemap"
	"matching-engine-lab/go/internal/enginev3_pool"
)

func main() {
	engineFlag := flag.String("engine", "v1", "Engine: v1, v2, or v3")
	latencyFlag := flag.Bool("latency", false, "Report per-command latency percentiles (p50, p99, p999)")
	flag.Parse()

	var eng core.Engine
	switch *engineFlag {
	case "v1":
		eng = enginev1_btree.New()
	case "v2":
		eng = enginev2_treemap.New()
	case "v3":
		eng = enginev3_pool.New()
	default:
		fmt.Fprintf(os.Stderr, "Unknown engine: %s. Use v1, v2, or v3.\n", *engineFlag)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	count := int64(0)
	start := time.Now()
	var latencies []int64

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		cmd, err := core.ParseCommand(line)
		if err != nil {
			rej := &core.Event{Rejected: &core.RejectedEvent{Reason: "parse error: " + err.Error()}}
			out, _ := core.SerializeEvent(rej)
			fmt.Fprintln(writer, out)
			writer.Flush()
			continue
		}

		cmdStart := time.Now()
		events, err := eng.Submit(*cmd)
		if *latencyFlag {
			latencies = append(latencies, time.Since(cmdStart).Nanoseconds())
		}
		if err != nil {
			rej := &core.Event{Rejected: &core.RejectedEvent{Reason: err.Error()}}
			out, _ := core.SerializeEvent(rej)
			fmt.Fprintln(writer, out)
			writer.Flush()
			continue
		}

		for _, e := range events {
			out, _ := core.SerializeEvent(&e)
			fmt.Fprintln(writer, out)
		}
		writer.Flush()
		count++
	}

	elapsed := time.Since(start)
	secs := elapsed.Seconds()
	opsPerSec := 0.0
	if secs > 0 {
		opsPerSec = float64(count) / secs
	}
	fmt.Fprintf(os.Stderr, "Processed %d commands in %v (%.0f ops/sec)\n", count, elapsed, opsPerSec)

	if *latencyFlag && len(latencies) > 0 {
		sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
		p50 := percentile(latencies, 50)
		p99 := percentile(latencies, 99)
		p999 := percentile(latencies, 99.9)
		fmt.Fprintf(os.Stderr, "Latency (ns): p50=%d p99=%d p999=%d\n", p50, p99, p999)
	}
}

func percentile(sorted []int64, p float64) int64 {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(float64(len(sorted)) * p / 100)
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}
