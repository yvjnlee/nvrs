# Define the virtual environment directory
VENV_DIR := venv

# Default target
.PHONY: all
all: install

# Create a virtual environment
.PHONY: venv
venv:
	python3 -m venv $(VENV_DIR)

# Install dependencies
.PHONY: install
install: venv
	$(VENV_DIR)/bin/pip install -r requirements.txt

# Run tests (assuming tests are in a tests/ directory)
.PHONY: test
test: install
	$(VENV_DIR)/bin/python -m unittest discover -s tests

# Run the main application
.PHONY: run
run: install
	$(VENV_DIR)/bin/python main.py

# Clean up environment
.PHONY: clean
clean:
	rm -rf $(VENV_DIR)
	find . -name "__pycache__" -type d -exec rm -r {} +

# Update requirements.txt
.PHONY: freeze
freeze:
	$(VENV_DIR)/bin/pip freeze > requirements.txt
