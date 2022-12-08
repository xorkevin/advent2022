.PHONY: bench

bench:
	find . -maxdepth 1 -type d -name 'day*' | sort | xargs -I{} sh -c 'cd {} && pwd && make bench'
