# Go params
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: build
build: 
	$(GOBUILD) -v
	@echo 'done'

test:
	$(GOTEST) -v -race ./... #run tests (TestXxx) excluding benchmarks
	@echo 'done'

bench:
	$(GOTEST) -v -race -run=XXX -bench=. ./... #run all benchmarks (BenchmarkXxx)
	@echo 'done'

test-bench:
	$(GOTEST) -v -race -bench=. ./... #run all tests and benchmarks (TestXxx and BenchmarkXxx)
	@echo 'done'
