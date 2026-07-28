// Package snapshot implements copy-on-write workspace snapshotting
// (APFS on macOS, Btrfs/ZFS on Linux, virtual disk overlays on Windows)
// used to fork a working directory into a sandbox and extract git diffs
// back out on completion. Not yet implemented.
package snapshot
