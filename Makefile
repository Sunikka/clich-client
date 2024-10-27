
BINDIR := bin
APP := $(BINDIR)/clich

build:
	go build -o $(APP) ./cmd


# Build and run executables
run: build
	$(APP)

# NOT IMPLEMENTED YET
# Add debug log window
run-dev: build
	$(BINDIR)/clich --dev




# Delete executable binaries
clean:
	rm -rf $(BINDIR)
