using System;
using System.Management;
using System.Runtime.Versioning;

namespace QuazaarMedia.Services
{
    [SupportedOSPlatform("windows")]
    public class BrightnessControlService
    {
        public void SetBrightness(int level)
        {
            try
            {
                level = Math.Clamp(level, 0, 100);
                using var searcher = new ManagementObjectSearcher("root\\WMI", "SELECT * FROM WmiMonitorBrightnessMethods");
                using var collection = searcher.Get();
                foreach (ManagementObject mObj in collection)
                {
                    mObj.InvokeMethod("WmiSetBrightness", new object[] { 1, level });
                }
                Console.WriteLine("{\"status\": \"ok\"}");
            }
            catch (Exception ex)
            {
                Console.WriteLine($"{{\"status\": \"error\", \"message\": \"SetBrightness Error: {ex.Message}\"}}");
            }
        }

        public void BrightnessUp()
        {
            try
            {
                int current = GetCurrentBrightnessLevel();
                SetBrightness(current + 5);
            }
            catch (Exception ex)
            {
                Console.WriteLine($"{{\"status\": \"error\", \"message\": \"BrightnessUp Error: {ex.Message}\"}}");
            }
        }

        public void BrightnessDown()
        {
            try
            {
                int current = GetCurrentBrightnessLevel();
                SetBrightness(current - 5);
            }
            catch (Exception ex)
            {
                Console.WriteLine($"{{\"status\": \"error\", \"message\": \"BrightnessDown Error: {ex.Message}\"}}");
            }
        }

        public void GetAndPrintBrightness()
        {
            try
            {
                int current = GetCurrentBrightnessLevel();
                Console.WriteLine($"{{\"status\": \"ok\", \"brightness\": {current}}}");
            }
            catch (Exception ex)
            {
                Console.WriteLine($"{{\"status\": \"error\", \"message\": \"GetBrightness Error: {ex.Message}\"}}");
            }
        }

        private int GetCurrentBrightnessLevel()
        {
            using var searcher = new ManagementObjectSearcher("root\\WMI", "SELECT CurrentBrightness FROM WmiMonitorBrightness");
            using var collection = searcher.Get();
            foreach (ManagementObject mObj in collection)
            {
                return Convert.ToInt32(mObj["CurrentBrightness"]);
            }
            throw new Exception("No brightness control found");
        }
    }
}
