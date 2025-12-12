using System;
using System.Threading.Tasks;
using Windows.Media.Control;
using QuazaarMedia.Utilities;

namespace QuazaarMedia.Services
{
    public class PlaybackControlService
    {
        public async Task ExecutePlaybackCommand(GlobalSystemMediaTransportControlsSession session, string action)
        {
            if (session == null)
            {
                Console.WriteLine("{\"status\": \"error\", \"message\": \"no active session\"}");
                return;
            }

            try
            {
                switch (action)
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
                        await SeekForward(session);
                        break;
                    case "seek_backward":
                        await SeekBackward(session);
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
                Console.WriteLine($"{{\"status\": \"error\", \"message\": \"{JsonHelper.Escape(ex.Message)}\"}}");
            }
        }

        private async Task SeekForward(GlobalSystemMediaTransportControlsSession session)
        {
            var timeline = session.GetTimelineProperties();
            if (timeline != null)
            {
                var newPos = timeline.Position.Add(TimeSpan.FromSeconds(10));
                if (newPos > timeline.EndTime) newPos = timeline.EndTime;
                await session.TryChangePlaybackPositionAsync(newPos.Ticks);
            }
        }

        private async Task SeekBackward(GlobalSystemMediaTransportControlsSession session)
        {
            var timeline = session.GetTimelineProperties();
            if (timeline != null)
            {
                var newPos = timeline.Position.Subtract(TimeSpan.FromSeconds(10));
                if (newPos < TimeSpan.Zero) newPos = TimeSpan.Zero;
                await session.TryChangePlaybackPositionAsync(newPos.Ticks);
            }
        }
    }
}
