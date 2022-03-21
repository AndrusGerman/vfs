vfs for golang  [![GoDoc](https://godoc.org/github.com/AndrusGerman/vfs?status.png)](https://godoc.org/github.com/AndrusGerman/vfs) 
======

vfs is library to support virtual filesystems. It provides basic abstractions of filesystems and implementations, like `OS` accessing the file system of the underlying OS and `memfs` a full filesystem in-memory.

Basic Usage
-----
```bash
$ go get github.com/AndrusGerman/vfs@master
```
Note: Always vendor your dependencies or fix on a specific version tag.

```go
import github.com/AndrusGerman/vfs
```

```go
// Create a vfs accessing the filesystem of the underlying OS
var osfs vfs.Filesystem = vfs.OS()
osfs.Mkdir("/tmp", 0777)

// Make the filesystem read-only:
osfs = vfs.ReadOnly(osfs) // Simply wrap filesystems to change its behaviour

// os.O_CREATE will fail and return vfs.ErrReadOnly
// os.O_RDWR is supported but Write(..) on the file is disabled
f, _ := osfs.OpenFile("/tmp/example.txt", os.O_RDWR, 0)

// Return vfs.ErrReadOnly
_, err := f.Write([]byte("Write on readonly fs?"))
if err != nil {
    fmt.Errorf("Filesystem is read only!\n")
}

// Create a fully writable filesystem in memory
mfs := memfs.Create()
mfs.Mkdir("/root", 0777)

// Create a vfs supporting mounts
// The root fs is accessing the filesystem of the underlying OS
fs := mountfs.Create(osfs)

// Mount a memfs inside /memfs
// /memfs may not exist
fs.Mount(mfs, "/memfs")

// This will create /testdir inside the memfs
fs.Mkdir("/memfs/testdir", 0777)

// This would create /tmp/testdir inside your OS fs
// But the rootfs `osfs` is read-only
fs.Mkdir("/tmp/testdir", 0777)
```

Replication:
-----
```go
// Create multiple fully writable filesystem in memory
inMemoryA := memfs.Create()
inMemoryB := memfs.Create()

// Join and replication events vfs
replivfs := replicationfs.NewReplication(inMemoryA, inMemoryB)

// Create a vfs accessing the filesystem of the underlying OS
primaryDisk := prefixfs.Create(vfs.OS(), "/myosfolder")
// Sync primaryDisk to => replivfs
err := replicationfs.Sync(nil, primaryDisk, replivfs)
if err != nil {
	panic(err)
}
```
Dump Data
----

```go
// Save Data
// Create a vfs accessing the filesystem of the underlying OS
primaryDisk := prefixfs.Create(vfs.OS(), "/myosfolder")
file, err := os.Create("mybackup.vfs")
if err != nil {
	panic(err)
}
// Save data vfs
err = dumpfs.NewDumpfs(primaryDisk, file) // save data in file
if err != nil {
	panic(err)
}
file.Close()

// Read save Data
// Create a fully writable filesystem in memory
mvfs := memfs.Create()
file, err = os.Open("mybackup.vfs")
if err != nil {
	panic(err)
}
err = dumpfs.GetDumpfs(file, mvfs)// Set save data in new fs
if err != nil {
	panic(err)
}
```
----
Check detailed examples below. Also check the [GoDocs](http://godoc.org/github.com/AndrusGerman/vfs).

Why should I use this lib?
-----

- Only Stdlib
- Easy to create your own filesystem
- Mock a full filesystem for testing (or use included `memfs`)
- Compose/Wrap Filesystems `ReadOnly(OS())` and write simple Wrappers
- Many features, see [GoDocs](http://godoc.org/github.com/AndrusGerman/vfs) and examples below

Features and Examples
-----

- [OS Filesystem support](http://godoc.org/github.com/AndrusGerman/vfs#example-OsFS)
- [ReadOnly Wrapper](http://godoc.org/github.com/AndrusGerman/vfs#example-RoFS)
- [DummyFS for quick mocking](http://godoc.org/github.com/AndrusGerman/vfs#example-DummyFS)
- [MemFS - full in-memory filesystem](http://godoc.org/github.com/AndrusGerman/vfs/memfs#example-MemFS)
- [MountFS - support mounts across filesystems](http://godoc.org/github.com/AndrusGerman/vfs/mountfs#example-MountFS)

Current state: MEGA ALPHA
-----

While the functionality is quite stable and heavily tested, interfaces are subject to change. 

    You need more/less abstraction? Let me know by creating a Issue, thank you.

Motivation
-----

I simply couldn't find any lib supporting this wide range of variation and adaptability.

Contribution
-----

Feel free to make a pull request. For bigger changes create a issue first to discuss about it.

thanks to [Benedikt Lang](https://github.com/blang/vfs), for starting this incredible project

disclaimer
-----
Much of the code needs to be refactored and improved, and there may be some things that don't work.
I appreciate any contribution and I will gladly revise it to include it.
Little things are also appreciated, including translations and code homogenization.

License
-----

See [LICENSE](LICENSE) file.
