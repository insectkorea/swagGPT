# SwagGPT

**SwagGPT** is an experimental CLI tool designed to automatically generate and add Swagger comments to Gin and Echo handler functions in Go projects. The tool uses OpenAI's GPT-4 model to analyze the handler functions and create comprehensive Swagger documentation, improving API documentation quality and consistency.

## Features

- **Automatic Swagger Comment Generation**: Uses OpenAI API to generate Swagger comments based on handler function content.
- **Supports Gin and Echo Frameworks**: Initially focused on Gin but extendable to other web frameworks like Echo.
- **Dry Run Mode**: Preview the generated comments without modifying the actual files.
- **Undo Changes**: Restore original files from backups if needed.

## Installation

To install the CLI tool, use the following command:

```sh
go get github.com/insectkorea/swagGPT
```

## Usage

### Set Up Environment

Make sure to set your OpenAI API key in the environment:

```sh
export OPENAI_API_KEY=your_openai_api_key
```

### Run the CLI Tool

To add Swagger comments to handler functions in a specified directory:

```sh
swagGPT add-comments --dir /path/to/your/code --config /path/to/config.yaml
```

You can use the `--dry-run` flag to preview the changes without writing them to files:

```sh
swagGPT add-comments --dir /path/to/your/code --config /path/to/config.yaml --dry-run
```

### Restore Original Files

To restore original files from backups:

```sh
swagGPT undo --dir /path/to/your/code
```

### Configuration

Customize the generated Swagger comments using the `configs/config.yaml` file.

```yaml
swagger:
  summary_template: "Summary for {function_name}"
  description_template: "Description for {function_name}"
  tags: ["example"]
  accept: "json"
  produce: "json"
  success_response: "200 {string} string \"OK\""
  router_template: "/example/{function_name} [get]"
```

### Example

Assuming your project structure is as follows:

```
/your/project/
├── handlers/
│   └── CreateByMember.go
├── configs/
│   └── config.yaml
```

Navigate to your project directory and run:

```sh
cd /your/project
swagGPT add-comments --dir ./handlers --config ./configs/config.yaml --dry-run
```

## Running Tests

To run the tests, use the following command:

```sh
go test ./...
```

### Contributing Guidelines

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -am 'Add new feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Create a new Pull Request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Disclaimer

This is an **experimental project**. The accuracy and reliability of the generated Swagger comments depend on the quality of the underlying AI model and the specific implementation details. Use this tool at your own risk, and always review the generated comments for correctness before deploying them in a production environment.