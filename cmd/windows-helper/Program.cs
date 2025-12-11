using System;
using System.Text;
using System.Threading.Tasks;
using Windows.Media.Control;

// Quazaar Media Helper - Native AOT Optimized
// Reads commands from Stdin, executes Media API calls, writes JSON to Stdout.

class Program
{
    static async Task Main(string[] args)
    {
        // 1. Setup Encoding for Pipe Communication (Critical for Go IPC)
        Console.InputEncoding = Encoding.UTF8;
        Console.OutputEncoding = Encoding.UTF8;

        try 
        {
            // Connect to Windows Media API
            var manager = await GlobalSystemMediaTransportControlsSessionManager.RequestAsync();
            
            // 2. Handshake: Tell Go we are ready
            Console.WriteLine("{\"status\": \"ready\"}");

            // 3. Command Loop
            string? line;
            while ((line = Console.ReadLine()) != null)
            {
                if (string.IsNullOrWhiteSpace(line)) continue;

                try 
                {
                    // Basic parsing
                    var action = ParseAction(line);
                    var session = manager.GetCurrentSession();

                    if (session == null) {
                        Console.WriteLine("{\"error\": \"no_session\"}");
                        continue;
                    }

                    switch (action)
                    {
                        case "info": 
                            await OutputInfo(session); 
                            break;
                        case "play": 
                            await session.TryPlayAsync(); 
                            Ok(); 
                            break;
                        case "pause": 
                            await session.TryPauseAsync(); 
                            Ok(); 
                            break;
                        case "play_pause": 
                        case "toggle":
                            await session.TryTogglePlayPauseAsync(); 
                            Ok(); 
                            break;
                        case "next": 
                            await session.TrySkipNextAsync(); 
                            Ok(); 
                            break;
                        case "prev": 
                            await session.TrySkipPreviousAsync(); 
                            Ok(); 
                            break;
                        case "seek_forward":
                            // Seek forward 10 seconds
                            var posFwd = session.GetTimelineProperties().Position.Add(TimeSpan.FromSeconds(10));
                            await session.TryChangePlaybackPositionAsync(posFwd.Ticks);
                            Ok();
                            break;
                        case "seek_backward":
                            // Seek backward 10 seconds
                            var posBack = session.GetTimelineProperties().Position.Subtract(TimeSpan.FromSeconds(10));
                            // Ensure we don't seek before 0
                            if (posBack.Ticks < 0) posBack = TimeSpan.Zero;
                            await session.TryChangePlaybackPositionAsync(posBack.Ticks);
                            Ok();
                            break;
                        case "seek":
                            // Parse position from JSON line manually or simply
                            // This is a simplified parser for "Args": {"Position": 12345}
                            long seekPos = ParseLongArg(line, "Position");
                            if (seekPos >= 0) 
                            {
                                await session.TryChangePlaybackPositionAsync(TimeSpan.FromMilliseconds(seekPos).Ticks);
                                Ok();
                            }
                            else 
                            {
                                Console.WriteLine("{\"error\": \"invalid_seek_position\"}");
                            }
                            break;
                        case "volume":
                             // Windows SMTC usually doesn't allow setting system volume directly for other apps.
                             // It relies on the app to handle volume.
                             // We can try to use key simulation for system volume as a fallback if needed,
                             // but TryChangePlaybackPositionAsync is for the session.
                             // Currently, SMTC API for volume (TryChangePlaybackRateAsync) is for speed, not audio level.
                             // Real system volume control requires different APIs (ISimpleAudioVolume), which is complex in AOT.
                             // For now, we return "not_supported" or implement a hack.
                             Console.WriteLine("{\"error\": \"volume_control_not_supported_via_smtc\"}");
                             break;
                        default:
                            Console.WriteLine("{\"error\": \"unknown_command\"}");
                            break;
                    }
                }
                catch (Exception ex)
                {
                    Console.WriteLine($"{{\"error\": \"{Escape(ex.Message)}\"}}");
                }
            }
        }
        catch (Exception e)
        {
            Console.WriteLine($"{{\"fatal_error\": \"{Escape(e.Message)}\"}}");
        }
    }

    static void Ok() => Console.WriteLine("{\"status\": \"ok\"}");

    static string ParseAction(string json)
    {
        if (json.Contains("play_pause") || json.Contains("toggle")) return "play_pause";
        if (json.Contains("play")) return "play";
        if (json.Contains("pause")) return "pause";
        if (json.Contains("next")) return "next";
        if (json.Contains("prev")) return "prev";
        if (json.Contains("seek_forward")) return "seek_forward";
        if (json.Contains("seek_backward")) return "seek_backward";
        if (json.Contains("seek")) return "seek"; // Check this after forward/back
        if (json.Contains("volume")) return "volume";
        if (json.Contains("info")) return "info";
        return "unknown";
    }

    // Helper to extract numeric args like "Position": 12345
    static long ParseLongArg(string json, string key)
    {
        try {
            string search = $"\"{key}\":";
            int idx = json.IndexOf(search);
            if (idx == -1) return -1;
            
            int start = idx + search.Length;
            // Find end of number (comma or closing brace)
            int endComma = json.IndexOf(',', start);
            int endBrace = json.IndexOf('}', start);
            
            int end = -1;
            if (endComma == -1) end = endBrace;
            else if (endBrace == -1) end = endComma;
            else end = Math.Min(endComma, endBrace);
            
            if (end == -1) return -1;
            
            string numStr = json.Substring(start, end - start).Trim();
            return long.Parse(numStr);
        } catch { return -1; }
    }

    static async Task OutputInfo(GlobalSystemMediaTransportControlsSession session)
    {
        var props = await session.TryGetMediaPropertiesAsync();
        var timeline = session.GetTimelineProperties();
        var info = session.GetPlaybackInfo();
        
        double positionMs = timeline.Position.TotalMilliseconds;
        double durationMs = timeline.EndTime.TotalMilliseconds;
        string status = info.PlaybackStatus.ToString();

        var sb = new StringBuilder();
        sb.Append("{");
        sb.Append($"\"Title\": \"{Escape(props.Title)}\",");
        sb.Append($"\"Artist\": \"{Escape(props.Artist)}\",");
        sb.Append($"\"Album\": \"{Escape(props.AlbumTitle)}\",");
        sb.Append($"\"Status\": \"{status}\",");
        sb.Append($"\"Position\": {positionMs},");
        sb.Append($"\"Duration\": {durationMs}");
        sb.Append("}");
        
        Console.WriteLine(sb.ToString());
    }

    static string Escape(string? s) 
    {
        if (string.IsNullOrEmpty(s)) return "";
        return s.Replace("\\", "\\\\")
                .Replace("\"", "\\\"")
                .Replace("\n", " ")
                .Replace("\r", "");
    }
}