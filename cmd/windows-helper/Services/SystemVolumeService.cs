using System;
using System.Runtime.InteropServices; 

namespace QuazaarMedia.Services
{
    public class VolumeControlService
    {
        public void SetVolume(int level)
        {
            try
            {
                var volume = GetVolumeObject();
                float scalar = Math.Clamp(level, 0, 100) / 100.0f;
                volume.SetMasterVolumeLevelScalar(scalar, Guid.Empty);
                Console.WriteLine("{\"status\": \"ok\"}");
            }
            catch (Exception ex) { Console.WriteLine($"{{\"status\": \"error\", \"message\": \"SetVolume Error: {ex.Message}\"}}"); }
        }

        public void ToggleMute()
        {
            try
            {
                var volume = GetVolumeObject();
                volume.GetMute(out bool isMuted);
                volume.SetMute(!isMuted, Guid.Empty);
                Console.WriteLine("{\"status\": \"ok\"}");
            }
            catch (Exception ex) { Console.WriteLine($"{{\"status\": \"error\", \"message\": \"ToggleMute Error: {ex.Message}\"}}"); }
        }

        public void VolumeUp()
        {
            try
            {
                var volume = GetVolumeObject();
                volume.GetMasterVolumeLevelScalar(out float currentLevel);
                float newLevel = Math.Clamp(currentLevel + 0.05f, 0.0f, 1.0f);
                volume.SetMasterVolumeLevelScalar(newLevel, Guid.Empty);
                Console.WriteLine("{\"status\": \"ok\"}");
            }
            catch (Exception ex) { Console.WriteLine($"{{\"status\": \"error\", \"message\": \"VolumeUp Error: {ex.Message}\"}}"); }
        }

        public void VolumeDown()
        {
            try
            {
                var volume = GetVolumeObject();
                volume.GetMasterVolumeLevelScalar(out float currentLevel);
                float newLevel = Math.Clamp(currentLevel - 0.05f, 0.0f, 1.0f);
                volume.SetMasterVolumeLevelScalar(newLevel, Guid.Empty);
                Console.WriteLine("{\"status\": \"ok\"}");
            }
            catch (Exception ex) { Console.WriteLine($"{{\"status\": \"error\", \"message\": \"VolumeDown Error: {ex.Message}\"}}"); }
        }

        private IAudioEndpointVolume GetVolumeObject()
        {
            IMMDeviceEnumerator deviceEnumerator = (IMMDeviceEnumerator)new MMDeviceEnumerator();
            IMMDevice speakers;
            deviceEnumerator.GetDefaultAudioEndpoint(EDataFlow.eRender, ERole.eMultimedia, out speakers);
            object o;
            speakers.Activate(typeof(IAudioEndpointVolume).GUID, 0, IntPtr.Zero, out o);
            return (IAudioEndpointVolume)o;
        }

        [ComImport, Guid("BCDE0395-E52F-467C-8E3D-C4579291692E")] internal class MMDeviceEnumerator { }
        [ComImport, Guid("A95664D2-9614-4F35-A746-DE8DB63617E6"), InterfaceType(ComInterfaceType.InterfaceIsIUnknown)]
        internal interface IMMDeviceEnumerator { int EnumAudioEndpoints(EDataFlow dataFlow, uint dwStateMask, out IntPtr ppDevices); int GetDefaultAudioEndpoint(EDataFlow dataFlow, ERole role, out IMMDevice ppEndpoint); int GetDevice(string pwstrId, out IMMDevice ppDevice); int RegisterEndpointNotificationCallback(IntPtr pClient); int UnregisterEndpointNotificationCallback(IntPtr pClient); }
        [ComImport, Guid("D666063F-1587-4E43-81F1-B948E807363F"), InterfaceType(ComInterfaceType.InterfaceIsIUnknown)]
        internal interface IMMDevice { int Activate(ref Guid iid, int dwClsCtx, IntPtr pActivationParams, [MarshalAs(UnmanagedType.IUnknown)] out object ppInterface); int OpenPropertyStore(int stgmAccess, out IntPtr ppProperties); int GetId(out IntPtr ppstrId); int GetState(out int pdwState); }
        [ComImport, Guid("5CDF2C82-841E-4546-9722-0CF74078229A"), InterfaceType(ComInterfaceType.InterfaceIsIUnknown)]
        internal interface IAudioEndpointVolume { int RegisterControlChangeNotify(IntPtr pNotify); int UnregisterControlChangeNotify(IntPtr pNotify); int GetChannelCount(out int pnChannelCount); int SetMasterVolumeLevel(float fLevelDB, Guid pguidEventContext); int SetMasterVolumeLevelScalar(float fLevel, Guid pguidEventContext); int GetMasterVolumeLevel(out float pfLevelDB); int GetMasterVolumeLevelScalar(out float pfLevel); int SetChannelVolumeLevel(uint nChannel, float fLevelDB, Guid pguidEventContext); int SetChannelVolumeLevelScalar(uint nChannel, float fLevel, Guid pguidEventContext); int GetChannelVolumeLevel(uint nChannel, out float pfLevelDB); int GetChannelVolumeLevelScalar(uint nChannel, out float pfLevel); int SetMute([MarshalAs(UnmanagedType.Bool)] bool bMute, Guid pguidEventContext); int GetMute(out bool pbMute); int GetVolumeStepInfo(out uint pnStep, out uint pnStepCount); int VolumeStepUp(Guid pguidEventContext); int VolumeStepDown(Guid pguidEventContext); int QueryHardwareSupport(out uint pdwHardwareSupportMask); int GetVolumeRange(out float pflVolumeMindB, out float pflVolumeMaxdB, out float pflVolumeIncrementdB); }
        internal enum EDataFlow { eRender, eCapture, eAll, EDataFlow_enum_count }
        internal enum ERole { eConsole, eMultimedia, eCommunications, ERole_enum_count }
    }
}