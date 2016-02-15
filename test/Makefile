JAVA_HOME = "/opt/jdk"

.PHONY: all build_for_benchflow 

all: build_for_benchflow

build: build_for_benchflow

build_jar:
	cd camunda/wfmsbenchmark && \
	JAVA_HOME=$(JAVA_HOME) make all && \
	cd ../.. && \
	cp camunda/wfmsbenchmark/build/wfmsbenchmark.jar ./wfmsbenchmark.jar

build_for_benchflow:
	mkdir -p wfmsbenchmark/drivers
	cd camunda/wfmsbenchmark && \
	JAVA_HOME=$(JAVA_HOME) make all && \
	cd ../.. && \
	cp camunda/wfmsbenchmark/build/wfmsbenchmark.jar wfmsbenchmark/drivers/wfmsbenchmark.jar && \
	zip -r wfmsbenchmark.zip wfmsbenchmark
