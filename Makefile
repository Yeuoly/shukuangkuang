# Makefile for running and cleaning sh files in the build folder

# List of all shell scripts in the build folder
SH_FILES := $(wildcard build/*.sh)

# Target to run all shell scripts
run: $(SH_FILES)
	@echo "Running all shell scripts in build folder"
	@for file in $^; do \
	    echo "Running $$file"; \
	    bash $$file; \
	done

# Target to clean compiled files
clean:
	@echo "Cleaning compiled files in build folder"
	@rm -f build/shukuangkuang*

# Phony targets to avoid conflicts with file names
.PHONY: run clean
