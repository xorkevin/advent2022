.PHONY: bench

bench:
	find . -maxdepth 1 -type d -name 'day*' | sort | xargs -I{} sh -c 'cd {} && pwd && make bench' 2>/dev/null | stdbuf -oL grep '\(Benchmark\|Time\)' | sed 's/Benchmark.*day/day/'
