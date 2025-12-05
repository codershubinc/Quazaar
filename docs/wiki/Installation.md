# Installation & Setup

## Prerequisites

- **Go**: Version 1.21 or higher.
- **Git**: For cloning the repository.
- **Music Player Support**:

  - **Linux**: `playerctl` is recommended for MPRIS integration.

    ```bash
    # Arch Linux
    sudo pacman -S playerctl

    # Ubuntu/Debian
    sudo apt install playerctl
    ```

  - **Windows**: No additional dependencies required.

## Installation Steps

1. **Clone the Repository**

   ```bash
   git clone https://github.com/codershubinc/Quazaar.git
   cd Quazaar
   ```

2. **Download Dependencies**

   ```bash
   go mod download
   ```

3. **Build the Server**

   ```bash
   go build -o quazaar ./cmd/server
   ```

4. **Run the Server**

   ```bash
   ./quazaar
   ```

The server will start on `http://0.0.0.0:8765`.

## Configuration

Quazaar uses a SQLite database stored at `~/.quazaar/quazaar.db`. The database is automatically initialized on the first run.

### Environment Variables

You can configure the server using environment variables (create a `.env` file in the root directory):

```env
PORT=8765
SPOTIFY_CLIENT_ID=your_client_id
SPOTIFY_CLIENT_SECRET=your_client_secret
SPOTIFY_REDIRECT_URI=http://localhost:8765/callback
```

## Running Tests

To ensure everything is set up correctly, you can run the test suite:

```bash
go test ./...
```
