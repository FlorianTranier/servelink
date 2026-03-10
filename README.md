# ServeLink

A lightweight HTTP file server for exposing directory contents over a REST API. Built with Go and packaged as a Docker image for easy deployment.

## Features

- **Directory Browsing**: List directory contents with configurable introspection depth
- **File Serving**: Serve individual files over HTTP
- **API Key Authentication**: Secure access with environment-based API keys
- **Security Hardening**: Built-in helmet middleware for security headers
- **Minimal Footprint**: Uses scratch-based Docker image for minimal size

## Docker Hub

The image is available on Docker Hub: `blanymon/servelink`

## Quick Start

### Using Docker Run

```bash
docker run -d \
  -p 8080:8080 \
  -v /path/to/your/files:/app/mnt \
  -e SUPER_SECRET_KEY=your-secret-key-here \
  blanymon/servelink
```

### Using Docker Compose

```yaml
version: '3.8'

services:
  servelink:
    image: blanymon/servelink
    ports:
      - "8080:8080"
    volumes:
      - /path/to/your/files:/app/mnt
    environment:
      - SUPER_SECRET_KEY=your-secret-key-here
```

## Configuration

### Environment Variables

- **`SUPER_SECRET_KEY`** (required): API key for authentication. All requests must include this key in the `x-api-key` header.

### Volumes

- **`/app/mnt`**: Mount your files/directories to this path in the container.

## API Endpoints

### List Directory Contents

```http
GET /dir?introspectionLevel=<level>
```

**Query Parameters:**
- `introspectionLevel` (optional, default: `0`): How many levels deep to recursively list subdirectories
  - `0`: Only list the root directory
  - `1`: List root and immediate subdirectories
  - `n`: List up to n levels deep

**Headers:**
- `x-api-key`: Your secret API key

**Example:**
```bash
curl -H "x-api-key: your-secret-key-here" \
  "http://localhost:8080/dir?introspectionLevel=2"
```

**Response:**
```json
{
  "Name": ".",
  "Path": ".",
  "IsDir": true,
  "Metadata": null,
  "Children": {
    "file.txt": {
      "Name": "file.txt",
      "Path": "mnt/file.txt",
      "IsDir": false,
      "Metadata": {...},
      "Children": {}
    },
    "subfolder": {
      "Name": "subfolder",
      "Path": "mnt/subfolder",
      "IsDir": true,
      "Metadata": {...},
      "Children": {...}
    }
  }
}
```

### Download/Serve File

```http
GET /file/:path
```

**Path Parameters:**
- `path`: Relative path to the file (from the mounted `/app/mnt` directory)

**Headers:**
- `x-api-key`: Your secret API key

**Example:**
```bash
curl -H "x-api-key: your-secret-key-here" \
  "http://localhost:8080/file/documents/example.pdf" \
  -o example.pdf
```

## Security Considerations

1. **Always set a strong API key** via the `SUPER_SECRET_KEY` environment variable
2. All requests require authentication via the `x-api-key` header
3. Do not expose the service to the internet without additional security measures (reverse proxy, firewall, etc.)
4. The service serves files from the mounted volume - ensure only intended files are accessible

## Building from Source

```bash
# Clone the repository
git clone https://github.com/floriantranier/servelink.git
cd servelink

# Build the Docker image
docker build -t servelink .

# Run the container
docker run -d \
  -p 8080:8080 \
  -v /path/to/files:/app/mnt \
  -e SUPER_SECRET_KEY=your-secret-key \
  servelink
```

## License

See the repository for license information.

## Author

Florian Tranier
