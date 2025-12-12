using System;
using System.Text;
using System.Text.Json;
using System.Threading.Tasks;

namespace QuazaarMedia
{
    class Program
    {
        private static int _commandCount = 0;
        private static MediaController? _controller;

        static async Task Main(string[] args)
        {
            // Ensure console uses UTF8
            Console.OutputEncoding = Encoding.UTF8;
            Console.InputEncoding = Encoding.UTF8;

            // Handshake for Go app
            Console.WriteLine("{\"status\": \"ready\"}");

            // Initialize controller
            _controller = new MediaController();
            await _controller.InitializeAsync();

            // Main command loop - Request/Response model
            await ReadCommands();
        }

        static async Task ReadCommands()
        {
            while (true)
            {
                string? line = null;
                try
                {
                    line = await Console.In.ReadLineAsync();
                }
                catch (Exception)
                {
                    break;
                }

                if (line == null)
                {
                    break;
                }

                try
                {
                    // Use Source Generator Context for deserialization
                    var cmd = JsonSerializer.Deserialize(line, AppJsonContext.Default.CommandObj);
                    if (cmd != null && _controller != null)
                    {
                        await _controller.HandleCommand(cmd);
                    }
                }
                catch (JsonException)
                {
                    Console.WriteLine("{\"status\": \"error\", \"message\": \"json parse error\"}");
                }
                catch (Exception ex)
                {
                    var msg = ex.Message.Replace("\\", "\\\\").Replace("\"", "\\\"").Replace("\n", " ").Replace("\r", "");
                    Console.WriteLine($"{{\"status\": \"error\", \"message\": \"{msg}\"}}");
                }

                // Periodic GC to keep memory low
                _commandCount++;
                if (_commandCount > 10)
                {
                    _commandCount = 0;
                    GC.Collect();
                }
            }
        }
    }
}
