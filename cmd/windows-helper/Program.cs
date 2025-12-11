using System.Text;
using Windows.Media.Control;
using QuazaarMedia;

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
                    var action = Utils.ParseAction(line);
                    var session = manager.GetCurrentSession();

                    if (session == null)
                    {
                        Console.WriteLine("{\"error\": \"no_session\"}");
                        continue;
                    }

                    switch (action)
                    {
                        case "info":
                            await MediaSessionHelper.OutputInfo(session);
                            break;
                        case "artwork":
                            await Artwork.SaveArtworkAsync(session);
                            break;
                        case "play":
                            await session.TryPlayAsync();
                            Utils.Ok();
                            break;
                        case "pause":
                            await session.TryPauseAsync();
                            Utils.Ok();
                            break;
                        case "play_pause":
                        case "toggle":
                            await session.TryTogglePlayPauseAsync();
                            Utils.Ok();
                            break;
                        case "next":
                            await session.TrySkipNextAsync();
                            Utils.Ok();
                            break;
                        case "prev":
                            await session.TrySkipPreviousAsync();
                            Utils.Ok();
                            break;
                        case "seek_forward":
                            // Seek forward 10 seconds
                            var posFwd = session.GetTimelineProperties().Position.Add(TimeSpan.FromSeconds(10));
                            await session.TryChangePlaybackPositionAsync(posFwd.Ticks);
                            Utils.Ok();
                            break;
                        case "seek_backward":
                            // Seek backward 10 seconds
                            var posBack = session.GetTimelineProperties().Position.Subtract(TimeSpan.FromSeconds(10));
                            // Ensure we don't seek before 0
                            if (posBack.Ticks < 0) posBack = TimeSpan.Zero;
                            await session.TryChangePlaybackPositionAsync(posBack.Ticks);
                            Utils.Ok();
                            break;
                        case "seek":
                            long seekPos = Utils.ParseLongArg(line, "Position");
                            if (seekPos >= 0)
                            {
                                await session.TryChangePlaybackPositionAsync(TimeSpan.FromMilliseconds(seekPos).Ticks);
                                Utils.Ok();
                            }
                            else
                            {
                                Console.WriteLine("{\"error\": \"invalid_seek_position\"}");
                            }
                            break;
                        case "volume":
                            Console.WriteLine("{\"error\": \"volume_control_not_supported_via_smtc\"}");
                            break;
                        default:
                            Console.WriteLine("{\"error\": \"unknown_command\"}");
                            break;
                    }
                }
                catch (Exception ex)
                {
                    Console.WriteLine($"{{\"error\": \"{Utils.Escape(ex.Message)}\"}}");
                }
            }
        }
        catch (Exception e)
        {
            Console.WriteLine($"{{\"fatal_error\": \"{Utils.Escape(e.Message)}\"}}");
        }
    }
}
