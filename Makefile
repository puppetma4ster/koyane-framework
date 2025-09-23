# Name of the compiled binary
BINARY = koyane-framework

# Config, database, and settings files
CONFIG_FILE = ./config.yaml
DB_FILE = ./wordLists.db
SETTINGS_FILE = ./settings.yaml

# Build directory (current directory)
BUILD_DIR = .

# Standard installation paths
BIN_PATH = /usr/local/bin          # where the binary will be installed
CONFIG_PATH = /etc/koyane-framework  # system-wide config
DB_PATH = /var/lib/koyane-framework  # system-wide database
SETTINGS_PATH = ~/.config/koyane-framework  # user-specific settings

# -------------------------
# Build target: compiles the program
# -------------------------
build:
	@echo "Start to compile source code..."
	go build -o "$(BINARY)"
	@echo "Compilation done!"

# -------------------------
# Install target: creates necessary directories and copies files
# -------------------------
install:
	@echo "Creating directories..."
	mkdir -p $(CONFIG_PATH) $(DB_PATH) $(SETTINGS_PATH)
	@echo "Copying files..."
	cp $(CONFIG_FILE) $(CONFIG_PATH)
	cp $(DB_FILE) $(DB_PATH)
	cp $(BINARY) $(BIN_PATH)
	cp $(SETTINGS_FILE) $(SETTINGS_PATH)
	@echo "Installation done!"

# -------------------------
# Clean target: removes the compiled binary
# -------------------------
clean:
	@echo "Cleaning..."
	rm -f $(BINARY)
	@echo "Cleaning done!"

# -------------------------
# Uninstall target: removes all installed files and directories
# -------------------------
uninstall:
	rm -f $(BIN_PATH)
	rm -rf $(CONFIG_PATH)
	rm -rf $(DB_PATH)
	rm -rf $(SETTINGS_PATH)
