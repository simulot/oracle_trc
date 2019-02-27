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

(2808) 13-FEB-2019 15:45:55:479 select sysdate from dual
(2808) 13-FEB-2019 15:45:55:713 SELECT PARAMETER, VALUE FROM SYS.NLS_DATABASE_PARAMETERS WHERE PARAMETER IN ('NLS_CHARACTERSET', 'NLS_NCHAR_CHARACTERSET')
(2808) 13-FEB-2019 15:45:56:166 SELECT K,V from EXT_FLAGS
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
- [_] Display executable associated with pid
- [_] Write tests independent from trace files (confidentiality) 
- [_] Sort outputs from several trace file in time order
- [_] Understand binary format of packets (help wanted)
- [_] Determine bind parameters value (help wanted)
- [_] Decode responses (help wanted)

