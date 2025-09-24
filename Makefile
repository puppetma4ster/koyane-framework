# Name of the compiled binary
BINARY = koyane-framework

# Config, database, and settings files
CONFIG_FILE = ./config.yaml
DB_FILE = ./wordLists.db
SETTINGS_FILE = ./settings.yaml

# Build directory (current directory)
BUILD_DIR = .

### Standard installation paths
# where the binary will be installed
BIN_PATH := /usr/local/bin
# system-wide config
CONFIG_PATH := /etc/koyane-framework
# system-wide database
DB_PATH := /var/lib/koyane-framework
# user-specific settings
SETTINGS_PATH := ~/.config/koyane-framework
# user saves
SAVE_PATH = ~/.koyane_framework_saves

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
	@echo "Uninstalling..."
	@echo ""
	@echo "Uninstall binary..."
	rm -f "$(BIN_PATH)/$(BINARY)"
	@echo "Done!"
	@echo "Uninstall configurations..."
	rm -rf $(CONFIG_PATH)
	@echo "Done!"
	@echo "Uninstall databases..."
	rm -rf $(DB_PATH)
	@echo "Done!"
	@echo "Uninstall Setting Files..."
	rm -rf $(SETTINGS_PATH)
	@echo "Done!"
	@echo "Uninstall Save Files..."
	rm -rf $(SAVE_PATH)
	@echo "Done!"
