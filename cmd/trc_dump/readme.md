# trc_dump

This program open client side Oracle traces, take packets dumps, and output them in a format like hex -C could provide.


```
Usage of trc_dump:
Display all packets contained into given files in hexadecimal format like hex -C would do.
  -after string
        Filter packets exchanged after this date. In same format as tsFormat parameter.
  -tsFormat string
        Timestamp format, oracle's way. (default "DD-MON-YYYY HH:MI:SS:FF3")
```

Output sample:
```
client_5928.trc(2247),12-FEB-2019 17:30:13:267,client.exe(5928),nsbasic_bsd:
00000000  00 1e 00 00 06 00 00 00  00 00 11 6b 04 02 03 68  |...........k...h|
00000010  02 05 d1 01 01 03 3b 05  01 02 03 00 01 01 00 00  |......;.........|
00000020  00                                                |.|


client_5928.trc(2259),12-FEB-2019 17:30:13:267,client.exe(5928),nsbasic_brc:
00000000  00 64 00 00 06 00 00 00  00 00 08 01 4c 4c 4f 72  |.d..........LLOr|
00000010  61 63 6c 65 20 44 61 74  61 62 61 73 65 20 31 31  |acle Database 11|
00000020  67 20 45 6e 74 65 72 70  72 69 73 65 20 45 64 69  |g Enterprise Edi|
00000030  74 69 6f 6e 20 52 65 6c  65 61 73 65 20 31 31 2e  |tion Release 11.|
00000040  32 2e 30 2e 34 2e 30 20  2d 20 36 34 62 69 74 20  |2.0.4.0 - 64bit |
00000050  50 72 6f 64 75 63 74 69  6f 6e 04 0b 20 04 00 09  |Production.. ...|
00000060  01 01 01 03 00 00 00 00                           |........|
```

On this sample, 
- client_5928.trc(2259) is the trace file name, and the line of the packed
- 12-FEB-2019 17:30:13:267 is the packet time stamp as written int the trace file
- client.exe(5928) is the name of the client application and its PID
- nsbasic_bsd is the type of packet as written in trace file




# Documentation


http://2014.zeronights.org/assets/files/slides/oracle-database-communication-protocol.pdf
https://medium.com/@iphelix/hacking-oracle-tns-listener-c74070bde8e4

