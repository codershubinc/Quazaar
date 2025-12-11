
using Windows.Media.Control;

namespace QuazaarMedia
{
    internal class Artwork
    {
        public static async Task<string> SaveArtworkAsync(GlobalSystemMediaTransportControlsSessionMediaProperties mediaProperties)
        {
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
                return artWorkPath;
            }

            if (mediaProperties.Thumbnail != null)
            {
                var thumbnailStreamRef = mediaProperties.Thumbnail;
                using var thumbnailStream = await thumbnailStreamRef.OpenReadAsync();
                using var fileStream = new FileStream(artWorkPath, FileMode.Create);
                await thumbnailStream.AsStreamForRead().CopyToAsync(fileStream);

                return artWorkPath;
            }

            return "noArtwork"; // If no thumbnail was available
        }
    }
}
