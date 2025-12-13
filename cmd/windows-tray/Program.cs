using System;
using System.Diagnostics;
using System.Drawing;
using System.IO;
using System.Windows.Forms;

namespace QuazaarTray
{
    static class Program
    {
        private static NotifyIcon? _trayIcon;
        private static Process? _quazaarProcess;

        [STAThread]
        static void Main()
        {
            try
            {
                Application.SetHighDpiMode(HighDpiMode.SystemAware);
                Application.EnableVisualStyles();
                Application.SetCompatibleTextRenderingDefault(false);

                // Start Quazaar
                StartQuazaar();

                // Setup Tray Icon
                _trayIcon = new NotifyIcon
                {
                    Text = "Quazaar",
                    Visible = true
                };

                // Load Icon
                string iconPath = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "icon.ico");
                if (File.Exists(iconPath))
                {
                    try
                    {
                        _trayIcon.Icon = new Icon(iconPath);
                    }
                    catch
                    {
                        try
                        {
                            // Fallback: Try loading as Image (PNG/JPG) and convert
                            using (var bmp = new Bitmap(iconPath))
                            {
                                _trayIcon.Icon = Icon.FromHandle(bmp.GetHicon());
                            }
                        }
                        catch (Exception ex)
                        {
                            Log($"Failed to load icon: {ex.Message}. Using default.");
                            _trayIcon.Icon = SystemIcons.Application;
                        }
                    }
                }
                else
                {
                    _trayIcon.Icon = SystemIcons.Application;
                }

                // Context Menu
                var contextMenu = new ContextMenuStrip();
                contextMenu.Items.Add("Open Quazaar", null, (s, e) => OpenBrowser());
                contextMenu.Items.Add("-");
                contextMenu.Items.Add("Exit", null, (s, e) => ExitApp());
                _trayIcon.ContextMenuStrip = contextMenu;

                // Double click to open
                _trayIcon.DoubleClick += (s, e) => OpenBrowser();

                // Open browser on start
                OpenBrowser();

                Application.Run();
            }
            catch (Exception ex)
            {
                MessageBox.Show($"Critical Error: {ex.Message}\n{ex.StackTrace}", "Quazaar Tray Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
            }
        }

        private static void StartQuazaar()
        {
            try
            {
                Log("Starting Quazaar process...");
                var exePath = Path.Combine(AppDomain.CurrentDomain.BaseDirectory, "quazaar.exe");
                if (!File.Exists(exePath))
                {
                    Log($"Error: quazaar.exe not found at {exePath}");
                    MessageBox.Show($"quazaar.exe not found at {exePath}", "Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
                    return;
                }

                var startInfo = new ProcessStartInfo
                {
                    FileName = exePath,
                    UseShellExecute = false,
                    CreateNoWindow = true,
                    WorkingDirectory = AppDomain.CurrentDomain.BaseDirectory
                };
                _quazaarProcess = Process.Start(startInfo);
                Log("Quazaar process started.");
            }
            catch (Exception ex)
            {
                Log($"Failed to start Quazaar: {ex.Message}");
                MessageBox.Show($"Failed to start Quazaar: {ex.Message}", "Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
            }
        }

        private static void OpenBrowser()
        {
            try
            {
                Log("Opening browser...");
                Process.Start(new ProcessStartInfo("http://localhost:8765") { UseShellExecute = true });
            }
            catch (Exception ex)
            {
                Log($"Failed to open browser: {ex.Message}");
            }
        }

        private static void ExitApp()
        {
            Log("Exiting application...");
            _trayIcon!.Visible = false;
            if (_quazaarProcess != null && !_quazaarProcess.HasExited)
            {
                try
                {
                    _quazaarProcess.Kill();
                    Log("Quazaar process killed.");
                }
                catch (Exception ex)
                {
                    Log($"Failed to kill process: {ex.Message}");
                }
            }
            Application.Exit();
        }

        private static void Log(string message)
        {
            Console.WriteLine($"[{DateTime.Now}] {message}");
        }
    }
}
