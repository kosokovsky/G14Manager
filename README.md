# ROGManager: A Replacement For Armoury Crate (mostly)

![Test and Build](https://github.com/zllovesuki/ROGManager/workflows/Test%20and%20Build/badge.svg) ![Build Release](https://github.com/zllovesuki/ROGManager/workflows/Build%20Release/badge.svg)

## Disclaimer

Your warranty is now void. Proceed at your own risk.

## Help Wanted

After some reverse engineering, it looks like `atkwmiacpi64.sys` hooks itself into the WMI as the provider (and talks to `ACPI\PNP0C14\ATK`). However, even if ATKAMIACPIIO the driver is running, without Asus Optimization running, WMI has no events populated.

Not sure why this is the case, but I will investigate into this further and aim to remove Asus Optimization as part of the requirement.

## Requirements

ROGManager requires "Asus Optimization" to be running as a Service, since "Asus Optimization" loads the `atkwmiacpi64.sys` driver and interpret ACPI events as key presses, and exposes a `\\.\ATKACPI` device to be used. You do not need any other softwares from Asus running to use ROGManager; you can safely uninstall them from your system. However, some softwares are installed as Windows driver, and you should disable them in Services:

![Running Services](images/services.png)

![Running Processes](images/processes.png)

The OSD functionality is provided by `AsusOSD.exe`, which should also be under "Asus Optimization." 

```
PS C:\Windows\System32\DriverStore\FileRepository\asussci2.inf_amd64_b12b0d488bd75133\ASUSOptimization> dir


    Directory: C:\Windows\System32\DriverStore\FileRepository\asussci2.inf_amd64_b12b0d488bd75133\ASUSOptimization


Mode                 LastWriteTime         Length Name
----                 -------------         ------ ----
------         7/28/2020     02:41           3684 ASUS Optimization 36D18D69AFC3.xml
------         7/28/2020     02:52         218024 AsusHotkeyExec.exe
------         7/28/2020     02:52         273832 AsusOptimization.exe
------         7/28/2020     02:53         262056 AsusOptimizationStartupTask.exe
------         7/28/2020     02:53         117160 AsusOSD.exe
------         7/28/2020     02:53         844200 AsusSplendid.exe
------         7/28/2020     02:53         177576 AsusWiFiRangeboost.exe
------         7/28/2020     02:53         184744 AsusWiFiSmartConnect.exe
------         7/28/2020     02:53          44680 atkwmiacpi64.sys
------         7/28/2020     02:53         236952 CCTAdjust.dll
------         7/28/2020     02:53         204184 VideoEnhance_v406_20180511_x64.dll
```

Recommend running ROGManager.exe on startup in Task Scheduler.

## Remapping the ROG Key

Use case: You can compile your `.ahk` to `.exe` and run your macros.

By default, it will launch Task Manager when you press the ROG Key once.

To specify which program to launch when pressed multiple times, pass your path to the desired program as argument to `-rog` multiple times. For example:

```
.\ROGManager.exe -rog "Taskmgr.exe" -rog "start Spotify.exe"
```

This will launch Task Manager when you press the ROG key once, and Spotify when you press twice.

## Changing the Fan Curve

For the initial release, you have to change fan curve in `system\thermal\default.go`. In a future release ROGManager will allow you to specify the fan curve without rebuilding the binary. However, the default fan curve should be sufficient for most users.

Use the `Fn + F5` key combo to cycle through all the profiles. Fanless -> Quiet -> Slient -> Performance.

The key combo has a time delay. If you press the combo X times, it will apply the the next X profile. For example, if you are currently on "Fanless" profile, pressing `Fn + F5` twice will apply the "Slient" profile.

## How to Build

1. Install golang 1.14+ if you don't have it already
2. Install mingw x86_64 for `gcc.exe`
2. Install `rsrc`: `go get github.com/akavel/rsrc`
3. Generate `syso` file: `\path\to\rsrc.exe -arch amd64 -manifest ROGManager.exe.manifest -ico go.ico -o ROGManager.exe.syso`
4. Build the binary: `.\build.ps1`

## Developing

Use `.\run.ps1` as it does not compile using CGo.