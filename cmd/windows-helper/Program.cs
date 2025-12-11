using System;
using System.IO;
using System.Text;
using System.Text.Json;
using System.Text.Json.Serialization;
using System.Threading.Tasks;
using Windows.Media.Control;
using Windows.Storage.Streams;
using System.Diagnostics;

namespace QuazaarMedia
{
    class Program
    {
        private static GlobalSystemMediaTransportControlsSessionManager _manager;
        private static int _commandCount = 0;

        static async Task Main(string[] args)
        {
            // Ensure console uses UTF8
            Console.OutputEncoding = Encoding.UTF8;
            Console.InputEncoding = Encoding.UTF8;

            Log("Starting Sidecar...");

            // Handshake for Go app
            Console.WriteLine("{\"status\": \"ready\"}");

            // Initialize manager once
            try
            {
                _manager = await GlobalSystemMediaTransportControlsSessionManager.RequestAsync();
                Log("Manager initialized");
            }
            catch (Exception ex)
            {
                Log($"Manager init failed: {ex}");
            }

            // Main command loop - Request/Response model
            await ReadCommands();
        }

        static void Log(string message)
        {
            try
            {
                File.AppendAllText("sidecar.log", $"{DateTime.Now}: {message}{Environment.NewLine}");
            }
            catch { }
        }

        static async Task ReadCommands()
        {
            while (true)
            {
                string line = null;
                try
                {
                    line = await Console.In.ReadLineAsync();
                }
                catch (Exception ex)
                {
                    Log($"ReadLine error: {ex}");
                    break;
                }

                if (line == null)
                {
                    Log("End of stream (stdin closed)");
                    break;
                }

                Log($"Received: {line}");

                try
                {
                    // Use Source Generator Context for deserialization
                    var cmd = JsonSerializer.Deserialize(line, AppJsonContext.Default.CommandObj);
                    if (cmd != null)
                    {
                        await HandleCommand(cmd);
                    }
                }
                catch (JsonException ex)
                {
                    Log($"JSON error: {ex}");
                    Console.WriteLine("{\"status\": \"error\", \"message\": \"json parse error\"}");
                }
                catch (Exception ex)
                {
                    Log($"Command error: {ex}");
                    Console.WriteLine($"{{\"status\": \"error\", \"message\": \"{Escape(ex.Message)}\"}}");
                }

                // Periodic GC to keep memory low
                _commandCount++;
                if (_commandCount > 100)
                {
                    _commandCount = 0;
                    GC.Collect();
                }
            }
        }

        static async Task HandleCommand(CommandObj cmd)
        {
            // Ensure manager exists
            if (_manager == null)
            {
                try
                {
                    _manager = await GlobalSystemMediaTransportControlsSessionManager.RequestAsync();
                }
                catch { }
            }

            if (_manager == null)
            {
                // If we still can't get manager, return error or idle
                if (cmd.Action == "info") Console.WriteLine("{\"status\": \"idle\"}");
                else Console.WriteLine("{\"status\": \"error\", \"message\": \"no manager\"}");
                return;
            }

            var session = _manager.GetCurrentSession();

            if (cmd.Action == "info")
            {
                await PrintStatus(session);
                return;
            }

            // Control commands
            if (session != null)
            {
                try
                {
                    switch (cmd.Action)
                    {
                        case "play_pause":
                            await session.TryTogglePlayPauseAsync();
                            break;
                        case "play":
                            await session.TryPlayAsync();
                            break;
                        case "pause":
                            await session.TryPauseAsync();
                            break;
                        case "next":
                            await session.TrySkipNextAsync();
                            break;
                        case "prev":
                            await session.TrySkipPreviousAsync();
                            break;
                        case "seek_forward":
                            var timeline = session.GetTimelineProperties();
                            if (timeline != null)
                            {
                                var newPos = timeline.Position.Add(TimeSpan.FromSeconds(10));
                                if (newPos > timeline.EndTime) newPos = timeline.EndTime;
                                await session.TryChangePlaybackPositionAsync(newPos.Ticks);
                            }
                            break;
                        case "seek_backward":
                            var timeline2 = session.GetTimelineProperties();
                            if (timeline2 != null)
                            {
                                var newPos = timeline2.Position.Subtract(TimeSpan.FromSeconds(10));
                                if (newPos < TimeSpan.Zero) newPos = TimeSpan.Zero;
                                await session.TryChangePlaybackPositionAsync(newPos.Ticks);
                            }
                            break;
                        case "seek":
                            // Expects "Position" in seconds (int or float)
                            // But CommandObj needs to support it.
                            // For now, basic support.
                            break;
                    }
                    Console.WriteLine("{\"status\": \"ok\"}");
                }
                catch (Exception ex)
                {
                    Console.WriteLine($"{{\"status\": \"error\", \"message\": \"{Escape(ex.Message)}\"}}");
                }
            }
            else
            {
                Console.WriteLine("{\"status\": \"error\", \"message\": \"no active session\"}");
            }
        }

        static async Task PrintStatus(GlobalSystemMediaTransportControlsSession session)
        {
            Log("PrintStatus: Start");
            if (session == null)
            {
                Console.WriteLine("{\"status\": \"idle\"}");
                Log("PrintStatus: Session null");
                return;
            }

            try
            {
                var props = await session.TryGetMediaPropertiesAsync();
                Log("PrintStatus: Got props");
                var timeline = session.GetTimelineProperties();
                var playbackInfo = session.GetPlaybackInfo();

                if (props == null || playbackInfo == null)
                {
                    Console.WriteLine("{\"status\": \"idle\"}");
                    Log("PrintStatus: Props or PlaybackInfo null");
                    return;
                }

                string artworkUri = "";
                if (props.Thumbnail != null)
                {
                    try
                    {
                        Log("PrintStatus: Saving artwork");
                        var filePath = Path.Combine(Path.GetTempPath(), "quazaar_art.jpg");
                        using (var stream = await props.Thumbnail.OpenReadAsync())
                        using (var readStream = stream.AsStreamForRead())
                        using (var fileStream = File.Create(filePath))
                        {
                            await readStream.CopyToAsync(fileStream);
                        }
                        artworkUri = filePath;
                        Log($"PrintStatus: Artwork saved to {artworkUri}");
                    }
                    catch (Exception ex) { Log($"Artwork error: {ex}"); }
                }

                // Calculate current position
                double currentPosition = 0;
                if (timeline != null)
                {
                    currentPosition = timeline.Position.TotalMilliseconds;
                    if (playbackInfo.PlaybackStatus == GlobalSystemMediaTransportControlsSessionPlaybackStatus.Playing)
                    {
                        var timeDiff = DateTimeOffset.UtcNow.Subtract(timeline.LastUpdatedTime).TotalMilliseconds;
                        currentPosition += timeDiff;
                    }
                    // Clamp to duration
                    if (currentPosition > timeline.EndTime.TotalMilliseconds) currentPosition = timeline.EndTime.TotalMilliseconds;
                }

                // Manual JSON construction
                var json = $@"{{
""status"": ""playing"",
""Title"": ""{Escape(props.Title)}"",
""Artist"": ""{Escape(props.Artist)}"",
""Album"": ""{Escape(props.AlbumTitle)}"",
""App"": ""{Escape(session.SourceAppUserModelId)}"",
""Status"": ""{(playbackInfo.PlaybackStatus == GlobalSystemMediaTransportControlsSessionPlaybackStatus.Playing ? "Playing" : "Paused")}"",
""Position"": {currentPosition},
""Duration"": {(timeline != null ? timeline.EndTime.TotalMilliseconds : 0)},
""ArtworkUri"": ""{Escape(artworkUri)}""
}}";

                Console.WriteLine(json.Replace(Environment.NewLine, ""));
                Log("PrintStatus: JSON sent");
            }
            catch (Exception ex)
            {
                Log($"PrintStatus Error: {ex}");
                Console.WriteLine($"{{\"status\": \"error\", \"message\": \"{Escape(ex.Message)}\"}}");
            }
        }

        static string Escape(string s)
        {
            if (string.IsNullOrEmpty(s)) return "";
            return s.Replace("\\", "\\\\").Replace("\"", "\\\"").Replace("\n", " ").Replace("\r", "");
        }
    }

    internal class CommandObj
    {
        public string Action { get; set; }
        public long Position { get; set; } // For seek
        public int Level { get; set; } // For volume
    }

    [JsonSerializable(typeof(CommandObj))]
    internal partial class AppJsonContext : JsonSerializerContext
    {
    }
}
