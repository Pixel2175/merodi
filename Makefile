BINARY := merodi
PREFIX := $(if $(filter 0,$(shell id -u)),/usr/bin,$(HOME)/.local/bin)
TARGET := $(PREFIX)/$(BINARY)

.PHONY: all build install uninstall clean test

all: build

build:
	@printf "\033[1;34m==>\033[0m Compiling $(BINARY)...\n"
	@go build -o $(BINARY) ./src/main.go
	@printf "\033[1;32m==>\033[0m Built \033[1m./$(BINARY)\033[0m\n"

install: build
	@printf "\033[1;34m==>\033[0m Installing $(BINARY) to \033[1m$(PREFIX)\033[0m...\n"
	@mkdir -p $(PREFIX)
	@install -m 755 $(BINARY) $(TARGET)
	@printf "\033[1;32m==>\033[0m Installed \033[1m$(TARGET)\033[0m\n"

uninstall:
	@printf "\033[1;34m==>\033[0m Removing $(TARGET)...\n"
	@rm -f $(TARGET)
	@printf "\033[1;32m==>\033[0m Uninstalled $(BINARY)\n"

clean:
	@printf "\033[1;34m==>\033[0m Removing local binary...\n"
	@rm -f $(BINARY)
	@printf "\033[1;32m==>\033[0m Clean\n"

test: build
	@clear
	@./$(BINARY) $(w)
