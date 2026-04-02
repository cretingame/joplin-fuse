# Joplin Fuse

Joplin Fuse is a Go-based tool that mounts your Joplin notes into a filesystem using FUSE (Filesystem in Userspace). This allows you to browse, read, and interact with your Joplin notes as if they were regular files on your system.

## Features

- Mount Joplin notes as a read-only filesystem
- Read and edit notes directly from your terminal or file explorer
- Simple command-line interface

## Roadmap

- [x] Note Editing
- [ ] File and folder modification dates
- [ ] Write Conflict Management
    - [ ] Detect Conflict
    - [ ] Copy the clonflicted note
- [ ] Joplin Link Conversion: `[MyNote](:/3d9be863f25945b88734bfc3012f6b0b)` -> `[MyNote](./MyNotebook/MyNote)`
- [ ] Removing superflous log messages

## Installation

### With Go

To install Joplin Fuse, use the following command:

```bash
go install github.com/cretingame/joplin-fuse@latest
```

To uninstall:
```bash
rm ~/go/bin/joplin-fuse
```

### Ubuntu/Debian Package

Download the last package with the *.deb extention from the release page: https://github.com/cretingame/joplin-fuse/releases .

Execute the command below with root privilege:
```bash
dpkg -i build/joplin-fuse_0.0.X_amd64.deb
```

To uninstall:
```bash
dpkg -r joplin-fuse
```

## Usage

```bash
joplin-fuse [MOUNTING POINT]
```

For example:

```bash
mkdir ~/JoplinMount
joplin-fuse ~/JoplinMount
```

This will mount your Joplin notes at the specified mount point.

## Requirements

- FUSE installed and configured on your system
- Joplin desktop app or server with API access enabled

## API Configuration

Joplin Fuse connects to the Joplin API. Make sure the API is enabled in your Joplin settings.

Default configuration assumes the API is available at http://127.0.0.1:41184. If you're using a different host or port, you may need to set environment variables or provide configuration options (update this section as appropriate to your implementation).

You can check the API status by visiting:

http://127.0.0.1:41184/ping

If the response is `"JoplinClipperServer"`, the API is running.

## Building

Clone the repository and install dependencies:

```bash
git clone https://github.com/yourusername/joplin-fuse.git
cd joplin-fuse
go build
```

You can then run the binary:

```bash
./joplin-fuse [MOUNTING POINT]
```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

GNU Affero General Public License v3.0 or later
