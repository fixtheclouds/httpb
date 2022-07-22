# HTTP Benchmarks

Implementing primitive HTTP benchmarking scripts in various languages just for fun.

## Basic requirements
- accept URL as the only required argument;
- fetch that URL in an endless loop;
- write time spent for each requests (if successful) and HTTP status (or an error) to STDOUT;
- list results after catching SIGINT:
  - minimum, average and maximum time for requests;
  - total requests count;
  - number of successful and failed requests.
