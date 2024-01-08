# MistCrawler

MistCrawler is a Go script for checking the validity of URLs from a list and saving valid URLs to an output file.

## Features

- Concurrently checks URLs for validity.
- Supports specifying the number of threads for checking URLs.
- Provides a silent mode to suppress stdout output.
- Includes a help option (`-h`) to display usage information.

## Usage

1. **Clone the Repository**

   ```bash
   git clone [https://github.com/yourusername/MistCrawler.git](https://github.com/MistSpark/MistCrawler.git)

2. **Build**
   ```bash
   cd MistCrawler
   go build MistCrawler.go

4. **Run**
   ```bash
   ./MistCrawler -u urls.txt -o valid_urls.txt -t 4 -silent
   ```
   This command checks the URLs in urls.txt, saves valid URLs to valid_urls.txt, uses 4 threads for concurrent checking, and runs in silent mode.

## License
This project is licensed under the MIT License.

## Acknowledgments
[Go Programming Language](https://golang.org/)


## Contributing
Contributions are welcome! Feel free to open an issue or submit a pull request.

## Author
MistSpark ([@MistSpark](https://github.com/MistSpark))
