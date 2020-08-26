#canlisten
The tool issues a "CAN Channel Open" command on startup. A baudrate of 1 MBit/s is selected.

```
$: canlisten --dev=/dev/ttyACM0 --filter="f.Id == 0x100 && f.Data[3] == 0x40"
```
