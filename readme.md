# oracle_trc

**WORK IN PROGRESS**

I wrote this program to extract sql queries from trc files generated on client side.



## Usage
```
oracle_trc trace [, trace...]
  -after string
        Filter queries executed after this date. In same format as tsFormat parameter.
  -tsFormat string
        Timestamp format, oracle's way. (default "DD-MON-YYYY HH:MI:SS:FF3")
```

The output looks like:

```
C:\>oracle_trc d:\logs\oracle\*.trc

c:\log\cli_6632.trc(416937) Adapter.exe(6632) 26-FEB-2019 11:49:36:409 SELECT PARAMETER, VALUE FROM SYS.NLS_DATABASE_PARAMETERS WHERE PARAMETER IN ('NLS_CHARACTERSET', 'NLS_NCHAR_CHARACTERSET')
c:\log\cli_6632.trc(417021) Adapter.exe(6632) 26-FEB-2019 11:49:36:411 delete from IP_LOCKS where LOC_DOC_ID = '3E6244199E57400FA65A563F94FC7EC2'
...
```

## Enabling trace files
Add following lines to SQLNET.ORA file

```
adr_base=off
TRACE_LEVEL_CLIENT = 16
TRACE_FILE_CLIENT = CLIENT
TRACE_DIRECTORY_CLIENT = d:\logs\oracle
```

## Do do

- [X] Display queries after a given date
- [X] Display executable associated with pid
- [_] Write tests independent from trace files (confidentiality) 
- [_] Sort outputs from several trace file in time order
- [_] Understand binary format of packets (help wanted)
- [_] Determine bind parameters value (help wanted)
- [_] Decode responses (help wanted)

