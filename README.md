# Crashlooper

A simple container that crash 💥 after a set amount of time ⏰

## Usage

```bash
$ docker build -t crashlooper .
$ docker run --rm -it crashlooper --help
Usage:
  crashlooper [flags]

Flags:
      --crash-after duration                 Server will crash itself after specified period (default=0 means never)
  -h, --help                                 help for crashlooper
      --log-level string                     Server log level (default "info")
      --memory-increment string              crashlooper memory usage increment
      --memory-increment-interval duration   crashlooper memory usage increment interval (default 1s)
      --memory-target string                 crashlooper memory usage target
      --port string                          Server bind port (default "3000")
```

## Example

```bash
docker run --rm -it crashlooper --crash-after 10s
```
