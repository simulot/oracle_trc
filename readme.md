# oracle_trc

**WORK IN PROGRESS**

I wrote this program to extract sql queries from trc files generated on client side.

## trc_dump
Dump packets from trc files.

See cmd/trc_dump/readme.md for details

## queries
Dump SQL queries found in trc file

Exemple 

``` sql
client_2548.trc(275031),05-NOV-2020 06:54:51:729, cliebt.exe(2548), Socket(1284), nsbasic_bsd:
select 'NO DUPLICATES' from DUAL where 0 = nvl(( 
select 
  sum( case when FL.ACTION_INDEX = 3 then 0 else 1 end) as C 
from 
  (select :1 DOC_ID, :2 SUPPLIER_NUM, :3 INVOICE_NUM, :4 INVOICE_DATE, :5 COMP_NO from DUAL )REF join  
  DOCS D on REF.DOC_ID <> D.DOC_ID and REF.INVOICE_NUM = D.INVOICE_NUM and REF.COMP_NO = D.COMP_NO and REF.SUPPLIER_NUM = D.SUPPLIER_NUM and D.STATUS_INDEX <> 4 
    and D.INVOICE_DATE between REF.INVOICE_DATE-180 and REF.INVOICE_DATE+180  
  left outer join (  
    select FL.DOC_ID, FL.ACTION_INDEX  from FLOW_LOG FL where FL.ACTION_INDEX = 3  
      and FL.SENDED_TO_TIMESTAMP = (select max(FL2.SENDED_TO_TIMESTAMP) from FLOW_LOG FL2 where FL.DOC_ID = FL2.DOC_ID ) 
    )FL on FL.DOC_ID = D.DOC_ID 
),0) 

  :1 = '50FEEB8B33844E69BF8C'
  :2 = '992_4189'
  :3 = 'TEST1234'
  :4 = 2020-11-04T00:00:00Z
  :5 = '992'
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
- [X] Write tests independent from trace files (confidentiality) 
- [X] Sort outputs from several trace file in time order
- [ ] Understand binary format of packets (help wanted)
- [X] Determine bind parameters value (help wanted)
- [ ] Decode responses (help wanted)



# Some information

Findings on TNS packets
- https://blog.pythian.com/repost-oracle-protocol/
- https://flylib.com/books/en/2.680.1/the_oracle_network_architecture.html
- https://github.com/wireshark/wireshark/blob/master/epan/dissectors/packet-tns.c
- https://github.com/sijms/go-ora
