namespace QuazaarMedia
{
    public static class Utils
    {
        public static void Ok() => Console.WriteLine("{\"status\": \"ok\"}");

        public static string ParseAction(string json)
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
            if (json.Contains("artwork_base64")) return "artwork_base64"; // New command
            if (json.Contains("artwork")) return "artwork";
            if (json.Contains("info")) return "info";
            return "unknown";
        }

        public static long ParseLongArg(string json, string key)
        {
            try
            {
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
            }
            catch { return -1; }
        }

        public static string Escape(string? s)
        {
            if (string.IsNullOrEmpty(s)) return "";
            return s.Replace("\\", "\\\\")
                    .Replace("\"", "\\\"")
                    .Replace("\n", " ")
                    .Replace("\r", "");
        }
    }
}
