# fortios-release

The `fortios-release` is script for creating a better release notes.
You can give the version(s) and it will remove duplicates and change the status
of the bug.

## Instalation from release

- `MacOS`

```bash
wget https://github.com/y0rune/fortios-release/releases/latest/download/fortios-release-darwin-amd64 -O fortios-release
chmod +x fortios-release
./fortios-release
```

- `amd64`

```bash
wget https://github.com/y0rune/fortios-release/releases/latest/download/fortios-release-linux-amd64 -O fortios-release
chmod +x fortios-release
./fortios-release
```

- `arm64`

```bash
wget https://github.com/y0rune/fortios-release/releases/latest/download/fortios-release-linux-arm64 -O fortios-release
chmod +x fortios-release
./fortios-release
```

- `arm`

```bash
wget https://github.com/y0rune/fortios-release/releases/latest/download/fortios-release-linux-arm -O fortios-release
chmod +x fortios-release
./fortios-release
```

## Instalation from source

```bash
git clone https://github.com/y0rune/fortios-release.git
go get
make build
```

## Arguments

```
  -recordsFile string
        Name of the unsorted records from versions (default "records.csv")
  -sorted
        Get a sorted release notes
  -sortedFile string
        Name of the sorted output file (default "final.csv")
  -version value
        Version(s) of the FortiOS
```

## Example of usage

Getting the known and resolved issues in version `7.0.0` and `7.0.1`.

```bash
$ ./fortios-release -version 7.0.0 -version 7.0.1 -recordsFile issues-from-7.0.0-to-7.0.1.csv -sortedFile final.csv -sorted

2022/12/10 13:09:19 Starting gathering links for version 7.0.0
2022/12/10 13:09:21 The knownIssuesUrl is https://docs.fortinet.com/document/fortigate/7.0.0/fortios-release-notes/236526/known-issues
2022/12/10 13:09:21 The resolvedIssuesUrl is https://docs.fortinet.com/document/fortigate/7.0.0/fortios-release-notes/289806/resolved-issues
2022/12/10 13:09:21 Getting the resolved issue for version 7.0.0
2022/12/10 13:09:21 Starting parsing https://docs.fortinet.com/document/fortigate/7.0.0/fortios-release-notes/289806/resolved-issues data to table
2022/12/10 13:09:22 Starting writing the to file issues-from-7.0.0-to-7.0.1.csv
2022/12/10 13:09:22 Getting the known issue for version 7.0.0
2022/12/10 13:09:22 Starting parsing https://docs.fortinet.com/document/fortigate/7.0.0/fortios-release-notes/236526/known-issues data to table
2022/12/10 13:09:22 Starting writing the to file issues-from-7.0.0-to-7.0.1.csv
2022/12/10 13:09:22 Starting gathering links for version 7.0.1
2022/12/10 13:09:23 The knownIssuesUrl is https://docs.fortinet.com/document/fortigate/7.0.1/fortios-release-notes/236526/known-issues
2022/12/10 13:09:23 The resolvedIssuesUrl is https://docs.fortinet.com/document/fortigate/7.0.1/fortios-release-notes/289806/resolved-issues
2022/12/10 13:09:23 Getting the resolved issue for version 7.0.1
2022/12/10 13:09:23 Starting parsing https://docs.fortinet.com/document/fortigate/7.0.1/fortios-release-notes/289806/resolved-issues data to table
2022/12/10 13:09:24 Starting writing the to file issues-from-7.0.0-to-7.0.1.csv
2022/12/10 13:09:24 Getting the known issue for version 7.0.1
2022/12/10 13:09:24 Starting parsing https://docs.fortinet.com/document/fortigate/7.0.1/fortios-release-notes/236526/known-issues data to table
2022/12/10 13:09:25 Starting writing the to file issues-from-7.0.0-to-7.0.1.csv
2022/12/10 13:09:25 Starting writing the to file final.csv
```

All known and resolved issues in version `7.0.0` and `7.0.1`.

```
$ cat issues-from-7.0.0-to-7.0.1.csv | head -n3
BugID,Description,Status,Version
650160,"When using email filter profile, emails are being queued due to IMAP proxy being in stuck state.",resolved,7.0.0
524571,Quarantined files cannot be fetched in the AV log page if the file was already quarantined under another protocol.,resolved,7.0.0
```

Final output

```
$ cat final.csv | head -n3
BugID,Description,Status,Version
650160,"When using email filter profile, emails are being queued due to IMAP proxy being in stuck state.",resolved,7.0.0
524571,Quarantined files cannot be fetched in the AV log page if the file was already quarantined under another protocol.,resolved,7.0.0
```
