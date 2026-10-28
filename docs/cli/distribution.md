# Distribution

Signing, installers, bundled libraries, file associations, CI and the patched Gio the CLI builds with.

## Signing

`giode build -release` builds for distribution and signs everything it can with the
certificates in `giode.toml`; without `-release`, Android uses the debug key and the
rest are unsigned.

| Platform | What signs it | Set up |
| --- | --- | --- |
| Android | `android.keystore` | `GIODE_KEYSTORE_PASS=... giode keystore` creates one. Keep it and its password safe: Google Play only accepts updates signed with it (with Play App Signing, it's your upload key). |
| Windows | `windows.signing`: an Authenticode `.pfx`, or a `command` | Signs the executable and the setup program, with a timestamp so signatures outlive the certificate. Uses `signtool` on Windows and `osslsigncode` elsewhere. For keys in a hardware token or a cloud service (Azure Trusted Signing, a KMS), give `command` with `{file}` where the file goes. `giode cert` makes a self-signed certificate to test with; for public releases, use one from a certificate authority so SmartScreen doesn't warn. |
| macOS | `macos.identity` (or `macos.profile`) | A "Developer ID Application" identity in your keychain. With `macos.notary_id` and the app-specific password in `$GIODE_NOTARY_PASS`, the app is also notarized, so Gatekeeper opens it without warnings. |
| iOS | `ios.profile` | Device builds always need a provisioning profile; the matching identity must be in your keychain. |

`giode sign <file>...` signs other Windows files the same way.

## Windows installer

`giode build windows-installer` makes `dist/windows/<App>-<version>-setup.exe` with NSIS:
a wizard (welcome, optional license, folder, progress, finish with a "Run" box), a Start
menu shortcut (and a desktop one with `desktop_shortcut`), an entry in Settings > Apps
with publisher, version and size, and an uninstaller. By default it installs for the
current user into `%LOCALAPPDATA%\Programs`, without an administrator prompt;
`per_machine = true` installs for everyone into Program Files. It runs silently with
`/S`. The uninstaller isn't signed yet.

## Native libraries

A cgo app linking shared libraries (libmpv and the dozens it brings) needs them in its
packages. With `[native] bundle = true`, giode reads the built binary's imports
(PE, ELF or Mach-O, recursively, so it works when cross-building), leaves out the
system's (Windows' DLLs, macOS' `/usr/lib` and `/System`, a Linux desktop's C library,
graphics, windowing and sound stacks) and bundles the rest:

| Target | Where they go |
| --- | --- |
| `windows`, `windows-installer`, `msix`, `winget` | next to the `.exe`, installed (and uninstalled) with it |
| `macos`, `dmg` | `Contents/Frameworks`, renamed `@rpath/<name>` with every reference rewritten and the executable given the run path (as `dylibbundler` does), signed with the app |
| `linux`, `appimage`, `flatpak`, `deb`, `rpm` | `lib/` next to the binary, or the package's `lib/<command>`; the binary's run path finds them |

```toml
[native]
bundle = true
exclude = ["libvulkan*"]              # never bundled

[native.windows]
search = ["C:/msys64/mingw64/bin"]    # where the DLLs are (MSYS2, a MinGW sysroot)

[native.macos]
search = ["/opt/homebrew/lib"]

[native.linux]
system = true                         # .deb and .rpm depend on the system's packages instead
deb_depends = ["libmpv2"]             # the .rpm requires the libraries' sonames
libraries = ["/usr/lib/x86_64-linux-gnu/libva-drm.so.2"]  # also bundled: plugins loaded at run time
```

Libraries a binary imports that aren't the system's and aren't found are listed in a
warning. `giode doctor` checks the tools: `install_name_tool` and `codesign` on macOS,
`windows.cc` and pkg-config for Windows cgo builds.

## File associations

`[[file_types]]` puts the app in "Open with" for those files on every desktop:

```toml
[[file_types]]
name = "MKV video"
extensions = ["mkv", "mk3d"]
mime_types = ["video/x-matroska"]

## CI for cgo apps

`giode ci`'s workflow installs what Gio needs; a cgo app's C libraries come from
`[ci] setup`, commands run on each runner before building. `windows_msys2` builds the
Windows packages on a Windows runner in MSYS2's MINGW64 environment, where libraries like
libmpv are a `pacman` away, instead of cross-building them on Linux:

```toml
[ci]
windows_msys2 = true

[ci.setup]
macos = "brew install mpv"
linux = "sudo apt-get install -y libmpv-dev"
windows = "pacman -S --noconfirm mingw-w64-x86_64-mpv"
```

## Android libraries

Android builds take libraries from Maven repositories, as Gradle would:
`android.dependencies` lists `group:artifact:version` coordinates, which
giode resolves with their dependencies from Google's Maven repository and
Maven Central (cached in the user cache directory), and
`android.libraries` adds local `.aar` and `.jar` files. Their classes are
dexed with the app's (multidex when needed), their resources linked with
R classes for their packages, their manifests merged into the app's, and
their native libraries (for the app's ABIs) and assets packaged.

A Go package can ask for libraries itself: an `android_dependencies.txt`
next to its Go files lists coordinates, and an `android_manifest.xml`
adds permissions and `<application>` components, for every app importing
it. That's how giode's `push` (Firebase Cloud Messaging), `store` (Play
Billing) and `review` (Play In-App Review) packages bring what they need;
`push` also needs `android.google_services`, the `google-services.json`
file from the Firebase console, which giode turns into the resources
Firebase reads, as Google's Gradle plugin does.

## Patched Gio

`giode build` and `giode run` build against a copy of your Gio version with giode's
patches, adding the `giode_gio` build tag:

- **Underlay hook**, on every platform: what giode's `Config.Underlay` and
  `Config.Renderer` use to draw video or 3D under the UI. The window calls the hook with
  the frame's render target before drawing the UI over it without clearing, Gio's OpenGL
  renderer saves and restores the GL state around its frames, and on Windows it can prefer
  OpenGL ES (ANGLE) over Direct3D. `renderer = "opengl"` under `[app]` adds the `nometal`
  tag to macOS and iOS builds, for Gio's OpenGL renderer there.

And two Android fixes Gio hasn't released:

- **Full refresh rate.** Gio asks Android's Choreographer for the next frame only after
  presenting the current one. Presenting waits for a vsync, so the request misses the next
  one, and a 120Hz display gets 60 frames a second while anything animates. The patched
  copy asks before drawing, and drops a frame callback that arrives after the animation
  stopped instead of drawing a frame for nothing.
- **Held keys repeat.** After every edit Gio restarts the soft keyboard's input connection
  if the editor's snippet of text differs from the keyboard's - which it does after any
  edit, since the editor's snippet window moves with the text. The restart cancels the
  keyboard's key repeat, so holding backspace deleted one character. The patched copy
  restarts only when the text changed under the keyboard (an app clearing or replacing the
  field), not when the editor just moved its window over text the keyboard already has.

Your `go.mod` isn't touched. The copy is made once per Gio version in the user cache
directory (`giode/gio/<version>-giode<n>`, made again when the patches change), and the build gets a temporary `go.mod` with
`replace gioui.org => <copy>`, passed with `-modfile` to the `go build` and `go list`
commands run in your module (`GOFLAGS` is left alone). The patch is skipped, with a note,
when go isn't in module mode, or your module already replaces `gioui.org`, vendors its
dependencies or is in a Go workspace, or when your Gio version doesn't have the code the patch expects. Set
`gio_patches = false` under `[app]` to build against plain Gio (or under `[android]` for
Android builds only).

