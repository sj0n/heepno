# Copilot Instructions for Heepno

## Project Overview
Heepno is a CLI tool for transcribing audio files using multiple providers. The architecture separates CLI command definitions from the transcription logic, with each provider implementing a common interface for consistency.

## Architecture & Key Components
- **cmd/main.go**: Entry point for the CLI.
- **pkg/root.go**: Defines the root Cobra command.
- **pkg/assemblyai.go, pkg/deepgram.go, pkg/openai.go**: Define the provider-specific subcommands (`aai`, `dg`, `openai`).
- **pkg/config/config.go**: Manages application configuration, including API keys.
- **pkg/interfaces/provider.go**: Defines the core `Provider` interface that all transcription services must implement.
- **pkg/interfaces/assemblyai.go, pkg/interfaces/deepgram.go, pkg/interfaces/openai.go**: The concrete implementations of the `Provider` interface for each service.
- **pkg/shared/**: Shared utilities for console output, file writing, and formatting.

### Data Flow
1.  User invokes a CLI command (e.g., `heepno aai <file>`).
2.  Configuration is loaded from `pkg/config/config.go`.
3.  The command is routed via `cmd/main.go` and `pkg/root.go` to the appropriate subcommand in `pkg/{provider}.go`.
4.  The subcommand instantiates the provider implementation from `pkg/interfaces/{provider}.go`.
5.  The implementation's `Transcribe` method is called, which handles API communication.
6.  Shared utilities in `pkg/shared/` manage console output and file writing.

## Developer Workflows
- **Build**: `go build -o heepno ./cmd`
- **Run**: `./heepno [command] [args]`
- **Test**: No explicit test files found; add tests under `pkg/` as needed.
- **Debug**: Use standard Go debugging tools (e.g., `dlv debug ./cmd`).

## Project Conventions
- **Interface-based Providers**: All transcription providers must implement the `Provider` interface from `pkg/interfaces/provider.go`.
- **Separation of Concerns**: CLI command logic is in `pkg/{provider}.go`, while the actual API interaction logic is in `pkg/interfaces/{provider}.go`.
- **Centralized Configuration**: All configuration is managed through the `pkg/config` package.
- **CLI Commands**: Defined in `pkg/{provider}.go` files and added to the root command in `pkg/root.go`.

## Integration & Dependencies
- Relies on external APIs: Deepgram, OpenAI, AssemblyAI.
- API keys are handled via environment variables, loaded by the `pkg/config` package.

## Key Files
- `cmd/main.go`: CLI entry point
- `pkg/root.go`: Root CLI command
- `pkg/config/config.go`: Application configuration
- `pkg/interfaces/provider.go`: Core provider interface
- `pkg/assemblyai.go`, `pkg/deepgram.go`, `pkg/openai.go`: Provider-specific CLI command definitions
- `pkg/interfaces/assemblyai.go`, `pkg/interfaces/deepgram.go`, `pkg/interfaces/openai.go`: Provider implementations
- `pkg/shared/`: