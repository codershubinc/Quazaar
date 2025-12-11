using System.Text;
using Windows.Media.Control;

namespace QuazaarMedia
{
    public static class MediaSessionHelper
    {
        public static async Task OutputInfo(GlobalSystemMediaTransportControlsSession session)
        {
            try
            {
                // Get synchronous data first to avoid COM threading issues after await
                var timeline = session.GetTimelineProperties();
                var info = session.GetPlaybackInfo();

                // Then get async data
                var props = await session.TryGetMediaPropertiesAsync();
                // player name eg spotify
                var appName = session.SourceAppUserModelId;

                double positionMs = timeline.Position.TotalMilliseconds;
                double durationMs = timeline.EndTime.TotalMilliseconds;
                string status = info.PlaybackStatus.ToString();

                // Extrapolate position if playing because SMTC doesn't update constantly
                if (info.PlaybackStatus == GlobalSystemMediaTransportControlsSessionPlaybackStatus.Playing)
                {
                    double rate = info.PlaybackRate ?? 1.0;
                    var timeSinceUpdate = DateTimeOffset.Now.Subtract(timeline.LastUpdatedTime);
                    var adjustedTime = TimeSpan.FromTicks((long)(timeSinceUpdate.Ticks * rate));
                    var currentPos = timeline.Position.Add(adjustedTime);

                    // Clamp to duration
                    if (currentPos > timeline.EndTime) currentPos = timeline.EndTime;
                    if (currentPos < timeline.StartTime) currentPos = timeline.StartTime;

                    positionMs = currentPos.TotalMilliseconds;
                }

                var artworkPath = await Artwork.SaveArtworkAsync(session);

                var sb = new StringBuilder();
                sb.Append("{");
                sb.Append($"\"Title\": \"{Utils.Escape(props.Title)}\",");
                sb.Append($"\"Artist\": \"{Utils.Escape(props.Artist)}\",");
                sb.Append($"\"Album\": \"{Utils.Escape(props.AlbumTitle)}\",");
                sb.Append($"\"Status\": \"{status}\",");
                sb.Append($"\"Position\": {positionMs},");
                sb.Append($"\"Duration\": {durationMs}");
                sb.Append($",\"App\": \"{Utils.Escape(appName)}\"").Replace(".exe", "");
                sb.Append($",\"ArtworkUri\": \"{Utils.Escape(artworkPath)}\"");
                sb.Append("}");

                Console.WriteLine(sb.ToString());
            }
            catch (Exception ex)
            {
                Console.WriteLine($"{{\"error\": \"{Utils.Escape(ex.GetType().Name + ": " + ex.Message)}\"}}");
            }
        }
    }
}
