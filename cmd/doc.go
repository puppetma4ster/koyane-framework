// Package cmd provides all command-line interface (CLI) commands for the Koyane Framework.
//
// The commands are implemented using the Cobra library and cover functionality such as
// analyzing wordlists, generating new lists, editing existing ones, searching, showing,
// and retrieving data. Each command handles argument parsing, flag validation, and calls
// the underlying core logic from other packages.
//
// Available commands include:
//
//   - analyze   : Analyze a wordlist or dataset
//   - generate  : Generate a new wordlist
//   - edit      : Edit an existing wordlist
//   - search    : Search for specific entries or wordlists
//   - show      : Show detailed information about a wordlist
//   - get       : Download a wordlist or related resource
//
// Example usage:
//
//	$ koyane search --name "rockyou"
//	$ koyane show 1
//	$ koyane generate -M ?d?d?L!aedc
//
// This package separates CLI handling from business logic to make testing and maintenance easier.
package cmd
