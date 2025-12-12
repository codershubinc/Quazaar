using Windows.Media.Control;
using QuazaarMedia.Utilities;

namespace QuazaarMedia.Services
{
    public class MediaStatusService
    {
        public async Task GetAndPrintStatus(GlobalSystemMediaTransportControlsSession session)
        {
            if (session == null)
            {
                Console.WriteLine("{\"status\": \"idle\"}");
                return;
            }

            try
            {
                var props = await session.TryGetMediaPropertiesAsync();
                var timeline = session.GetTimelineProperties();
                var playbackInfo = session.GetPlaybackInfo();

                if (props == null || playbackInfo == null)
                {
                    Console.WriteLine("{\"status\": \"idle\"}");
                    return;
                }

                string artworkUri = await ArtworkUriHelper.ExtractArtwork(props.Thumbnail);
                double currentPosition = CalculateCurrentPosition(timeline, playbackInfo);

                PrintStatusJson(props, session, playbackInfo, currentPosition, timeline, artworkUri);
            }
            catch (Exception ex)
            {
                Console.WriteLine($"{{\"status\": \"error\", \"message\": \"{JsonHelper.Escape(ex.Message)}\"}}");
            }
        }
        private double CalculateCurrentPosition(
            GlobalSystemMediaTransportControlsSessionTimelineProperties timeline,
            GlobalSystemMediaTransportControlsSessionPlaybackInfo playbackInfo)
        {
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
                if (currentPosition > timeline.EndTime.TotalMilliseconds)
                    currentPosition = timeline.EndTime.TotalMilliseconds;
            }
            return currentPosition;
        }

        private void PrintStatusJson(
            GlobalSystemMediaTransportControlsSessionMediaProperties props,
            GlobalSystemMediaTransportControlsSession session,
            GlobalSystemMediaTransportControlsSessionPlaybackInfo playbackInfo,
            double currentPosition,
            GlobalSystemMediaTransportControlsSessionTimelineProperties timeline,
            string artworkUri)
        {
            var json = $@"{{
            ""status"": ""playing"",
            ""Title"": ""{JsonHelper.Escape(props.Title)}"",
            ""Artist"": ""{JsonHelper.Escape(props.Artist)}"",
            ""Album"": ""{JsonHelper.Escape(props.AlbumTitle)}"",
            ""App"": ""{JsonHelper.Escape(session.SourceAppUserModelId)}"",
            ""Status"": ""{(playbackInfo.PlaybackStatus == GlobalSystemMediaTransportControlsSessionPlaybackStatus.Playing ? "Playing" : "Paused")}"",
            ""Position"": {currentPosition},
            ""Duration"": {(timeline != null ? timeline.EndTime.TotalMilliseconds : 0)},
            ""ArtworkUri"": ""{JsonHelper.Escape(artworkUri)}""
            }}";

            Console.WriteLine(json.Replace(Environment.NewLine, ""));
        }
    }
}
