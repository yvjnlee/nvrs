# Define the virtual environment directory
VENV_DIR := venv

# Default target
.PHONY: all
all: install

# Create a virtual environment if it doesn't exist
$(VENV_DIR)/bin/activate: 
	python3 -m venv $(VENV_DIR)

# Install dependencies
.PHONY: install
install: $(VENV_DIR)/bin/activate
	$(VENV_DIR)/bin/pip install -r requirements.txt

# Run tests in the tests/ directory
.PHONY: test
test: install
	$(VENV_DIR)/bin/python -m unittest discover -s tests

# Run the gRPC client
.PHONY: run-grpc-client
run-grpc-client:
	$(VENV_DIR)/bin/python grpc_client/grpc_client.py

# Run the main agent application
.PHONY: run-agent
run-agent:
	$(VENV_DIR)/bin/python agents/agent.py

# Clean up environment
.PHONY: clean
clean:
	rm -rf $(VENV_DIR)
	find . -name "__pycache__" -type d -exec rm -r {} +

# Update requirements.txt
.PHONY: freeze
freeze:
	$(VENV_DIR)/bin/pip freeze > requirements.txt
