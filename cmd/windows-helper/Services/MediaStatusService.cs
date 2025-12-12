using System;
using System.IO;
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
                LogError("Session is null");
                Console.WriteLine("{\"status\": \"idle\"}");
                return;
            }

            try
            {
                GlobalSystemMediaTransportControlsSessionMediaProperties props = null;
                GlobalSystemMediaTransportControlsSessionTimelineProperties timeline = null;
                GlobalSystemMediaTransportControlsSessionPlaybackInfo playbackInfo = null;

                // Try to get props with timeout
                try
                {
                    props = await session.TryGetMediaPropertiesAsync();
                }
                catch (Exception ex)
                {
                    LogError($"TryGetMediaPropertiesAsync failed: {FormatExceptionDetails(ex)}");
                }

                // Try to get timeline
                try
                {
                    timeline = session.GetTimelineProperties();
                }
                catch (Exception ex)
                {
                    LogError($"GetTimelineProperties failed: {FormatExceptionDetails(ex)}");
                }

                // Try to get playback info
                try
                {
                    playbackInfo = session.GetPlaybackInfo();
                }
                catch (Exception ex)
                {
                    LogError($"GetPlaybackInfo failed: {FormatExceptionDetails(ex)}");
                }

                if (props == null || playbackInfo == null)
                {
                    LogError("Media properties or playback info is null");
                    Console.WriteLine("{\"status\": \"idle\"}");
                    return;
                }

                string artworkUri = await ArtworkUriHelper.ExtractArtwork(props.Thumbnail);
                double currentPosition = CalculateCurrentPosition(timeline, playbackInfo);

                PrintStatusJson(props, session, playbackInfo, currentPosition, timeline, artworkUri);
            }
            catch (Exception ex)
            {
                var details = FormatExceptionDetails(ex);
                LogError($"GetAndPrintStatus error: {details}");
                var msg = details;
                Console.WriteLine($"{{\"status\": \"error\", \"message\": \"{JsonHelper.Escape(msg)}\"}}");
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

        private static string FormatExceptionDetails(Exception ex)
        {
            var details = ex.GetType().Name;

            // For COMException, include HResult
            if (ex is System.Runtime.InteropServices.COMException comEx)
            {
                details += $" [HRESULT: 0x{comEx.HResult:X8}]";
            }

            if (!string.IsNullOrEmpty(ex.Message))
            {
                details += $" - {ex.Message}";
            }

            if (ex.InnerException != null)
            {
                details += $" | Inner: {ex.InnerException.GetType().Name}";
            }

            return details;
        }

        private static void LogError(string message)
        {
            try
            {
                string logPath = Path.Combine(Path.GetTempPath(), "sidecar.log");
                string timestamp = DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss.fff");
                File.AppendAllText(logPath, $"[{timestamp}] MediaStatusService: {message}\n");
            }
            catch { /* Ignore logging errors */ }
        }
    }
}
