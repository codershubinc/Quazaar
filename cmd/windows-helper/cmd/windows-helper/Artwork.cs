
using Windows.Media.Control;

namespace QuazaarMedia
{
    internal class Artwork
    {
        public static async Task<string> SaveArtworkAsync(GlobalSystemMediaTransportControlsSession currentSession)
        {
            var mediaProperties = await currentSession.TryGetMediaPropertiesAsync();
            if (mediaProperties == null) return "noArtwork";

            var generateFileName = $"{mediaProperties.Title}.jpg";
            var invalidChars = Path.GetInvalidFileNameChars();
            var sanitizedFileName = string.Join("_", generateFileName.Split(invalidChars, StringSplitOptions.RemoveEmptyEntries)).TrimEnd('.');

            var tempDir = Path.Combine(Directory.GetCurrentDirectory(), "temp");
            // Ensure the directory exists
            Directory.CreateDirectory(tempDir);

            var artWorkPath = Path.Combine(tempDir, sanitizedFileName);

            if (File.Exists(artWorkPath))
            {
                // Console.WriteLine("alredy stored  ::skiping storing ....");
                // Console.WriteLine($"{{\"artwork_path\": \"{Utils.Escape(artWorkPath)}\"}}");
                return artWorkPath;
            }

            if (mediaProperties.Thumbnail != null)
            {
                var thumbnailStreamRef = mediaProperties.Thumbnail;
                using var thumbnailStream = await thumbnailStreamRef.OpenReadAsync();
                using var fileStream = new FileStream(artWorkPath, FileMode.Create);
                await thumbnailStream.AsStreamForRead().CopyToAsync(fileStream);

                // Console.WriteLine($"{{\"artwork_path\": \"{Utils.Escape(artWorkPath)}\"}}");
                return artWorkPath;
            }

            return "noArtwork"; // If no thumbnail was available
        }
    }
}
