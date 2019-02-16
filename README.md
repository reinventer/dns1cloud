# 1cloud DNS API

This unofficial implementation of API of 1cloud DNS hosting

## Example
```
cli := dns1cloud.New("aaaaaaaaaaaaaaa")
record, err := cli.AddRecord(
	context.Background(),
	1234, // id of domain
	dns1cloud.Record{
		TypeRecord: dns1cloud.RecordTypeCNAME,
		HostName: "domain.com",
		MnemonicName: "www",
		TTL: 600,
	},
)
```
