# SwagGPT

**SwagGPT** is an experimental CLI tool designed to automatically generate and add Swagger comments to Gin and Echo handler functions in Go projects. The tool uses OpenAI's GPT-4 model to analyze the handler functions and create comprehensive Swagger documentation, improving API documentation quality and consistency.

## Features

- **Automated Swagger Comment Creation**: Leverages OpenAI API to automatically produce Swagger comments from the content of handler functions.
- **Compatibility with Gin and Echo Frameworks**: Primarily designed for Gin, but can be extended to support other web frameworks such as Echo.
- **Dry Run Feature**: Allows users to preview the generated comments without altering the actual files.
- **Cost Prediction**: Provides an estimated cost prior to execution.
- **Concurrent Execution**: Enables simultaneous processing of multiple handler functions, improving the efficiency and speed of the documentation generation process.

## Installation

To install the CLI tool, use the following command:

```sh
go install github.com/insectkorea/swagGPT/cmd/swaggpt@latest
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
swaggpt add-comments --dir /path/to/your/code --model gpt-4o 
```

If you have separate file that defines routes, you can add 

```sh
swaggpt add-comments --dir /path/to/your/code --model gpt-4o --route-file /path/to/your/route-file
```

Please make sure your files are under source version control, as swagGTP will overwrite contents.

You can use the `--dry-run` flag to preview the changes without writing them to files:

```sh
swaggpt add-comments --dir /path/to/your/code --model gpt-4o --dry-run
```

Note that while dry-run does not write to your files, but it does make API requests to Open AI.

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

## Caution

This may cause unexpected amount of billing as the program needs to pass the whole handler code to correctly understand the context.