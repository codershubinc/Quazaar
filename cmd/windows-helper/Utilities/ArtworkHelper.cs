using Windows.Storage.Streams;

namespace QuazaarMedia.Utilities
{
    public static class ArtworkUriHelper
    {
        public static async Task<string> ExtractArtwork(IRandomAccessStreamReference thumbnail)
        {
            string artworkUri = "";
            if (thumbnail != null)
            {
                try
                {
                    var filePath = Path.Combine(Path.GetTempPath(), "quazaar_art.jpg");
                    using (var stream = await thumbnail.OpenReadAsync())
                    using (var readStream = stream.AsStreamForRead())
                    using (var fileStream = File.Create(filePath))
                    {
                        await readStream.CopyToAsync(fileStream);
                    }
                    artworkUri = filePath;
                }
                catch (Exception) { }
            }
            return artworkUri;
        }

    }
}