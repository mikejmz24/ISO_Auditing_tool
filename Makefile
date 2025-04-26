# Build the application
all: build

build:
	@echo "Building..."
	@templ generate
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Create DB container
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@if [ -z "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
			path="./tests/..."; \
	elif [ -n "$(word 2, $(MAKECMDGOALS))" ]; then \
			path="./tests/$(word 2, $(MAKECMDGOALS))"; \
			if [ -n "$(word 3, $(MAKECMDGOALS))" ]; then \
					path="-run $(word 3, $(MAKECMDGOALS)) $$path"; \
			fi; \
	fi; \
	go test -v $$path 2>&1 | \
	awk ' \
    /--- (F|P)/ {print} \
    /Error: / {flag = 1} \
		/Test: / {print; flag = 0} \
		flag {print}' | \
	awk '\
		BEGIN {\
			test_count = 0;\
			in_error = 0;\
			in_messages = 0;\
		}\
		/--- (F|P)/ {\
			print $$0;\
		}\
		/Error:/ {\
			current_error = $$0;\
			test_count++;\
			errors[test_count] = current_error;\
			in_error = 1;\
			in_messages = 0;\
			next;\
		}\
		/Test:/ {\
			in_error = 0;\
			in_messages = 0;\
			test_text = substr($$0, 23);\
			tests[test_count] = test_text;\
			next;\
		}\
		/Messages:/ {\
			in_error = 0;\
			in_messages = 1;\
			messages[test_count] = $$0;\
			next;\
		}\
		{\
			if (in_error) {\
				errors[test_count] = errors[test_count] "\n" $$0;\
			}\
			if (in_messages && !/---/) {\
				messages[test_count] = messages[test_count] "\n" $$0;\
			}\
		}\
		/--- FAIL:/ {\
			start = index($$0, ": ") + 2;\
			end = index($$0, " (");\
			failed_test = substr($$0, start, end - start); \
			for (i = 1; i <= test_count; i++) {\
				if (tests[i] == failed_test) {\
					if (errors[i]) print errors[i];\
					if (messages[i]) print messages[i];\
					delete tests[i];\
					delete errors[i];\
					delete messages[i];\
				}\
			}\
		}' | \
	awk ' \
		{if (index($$0, "/") > 0) { \
				split($$0, parts, "/"); \
				split(parts[1], colin_sub, ": "); \
				print colin_sub[1] ": " parts[length(parts)] \
		} else { \
				print $$0 \
		}}' | \
	sed -e 's/--- PASS/\x1b[32m--- PASS\x1b[0m/' \
			-e 's/--- FAIL/\x1b[31m--- FAIL\x1b[0m/' \
			-e 's/--- SKIP/\x1b[33m--- SKIP\x1b[0m/' \
			-e 's/.*warning.*/\x1b[0;33m&\x1b[1;33m/' \
			-e '/actual  :/s/.*/\x1b[38;5;9m&\x1b[0m/' \
			-e '/expected:/s/.*/\x1b[38;5;10m&\x1b[0m/' \
			-e 's/\(Expected: \)\([^,]*\)\(, Got\)/\x1b[38;5;10m\1\2\x1b[0m\3/' \
			-e 's/\(Got: \)\(.*\)$$/\x1b[38;5;9m\1\2\x1b[0m/'

%:
	@:

.PHONY: migrate %
# Handle migrations with optional file and direction arguments
migrate:
	@direction="up"; \
	file=""; \
	if [ -n "$(word 2, $(MAKECMDGOALS))" ]; then \
		if [ -n "$(word 3, $(MAKECMDGOALS))" ]; then \
			if [ "$(word 3, $(MAKECMDGOALS))" = "up" ] || [ "$(word 3, $(MAKECMDGOALS))" = "down" ]; then \
				file="$(word 2, $(MAKECMDGOALS))"; \
				direction="$(word 3, $(MAKECMDGOALS))"; \
			else \
				file="$(word 2, $(MAKECMDGOALS))"; \
				direction="up"; \
			fi; \
		else \
			if [ "$(word 2, $(MAKECMDGOALS))" = "up" ] || [ "$(word 2, $(MAKECMDGOALS))" = "down" ]; then \
				direction="$(word 2, $(MAKECMDGOALS))"; \
			else \
				file="$(word 2, $(MAKECMDGOALS))"; \
			fi; \
		fi; \
	fi; \
	if [ -n "$$file" ]; then \
		echo "Running $$file migration $$direction..."; \
		go run cmd/api/main.go migrate "$$direction" "$$file"; \
	else \
		echo "Running all migrations $$direction..."; \
		go run cmd/api/main.go migrate "$$direction" ""; \
	fi

# Handle arguments passed to migrate
%:
	@:

# Seed the database
seed:
	@go run cmd/api/main.go seed

# Truncate the database tables
truncate:
	@go run cmd/api/main.go truncate

# Refresh truncates and seeds the database
refresh:
	@go run cmd/api/main.go refresh

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi


.PHONY: all build run test clean
