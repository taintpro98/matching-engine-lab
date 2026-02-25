package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"matching-engine-lab/go/internal/core"
	"matching-engine-lab/go/internal/enginev1_btree"
	"matching-engine-lab/go/internal/enginev2_treemap"
	"matching-engine-lab/go/internal/enginev3_pool"
)

func main() {
	engineFlag := flag.String("engine", "v1", "Engine: v1, v2, or v3")
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

		events, err := eng.Submit(*cmd)
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
}
