# NES emulator in Golang

This project started as an NES emulator which was going to be built in Go. Unfortunately I had to scrap that idea because I could not test the code I have written in Go since many programs which were supposed to help me in testing my 6502 processor are not supported for Go (support is for C and C++). Due to this reason, this project will be paused here for a while but I did learn a lot from basic computer architecture to many concepts around how a processor actually works.

### 6502 processor implementation
The implementation of the processor is kept faithful to the one used in the original device. The tests will be written soon, most probably Blargg's instr_test-v5 will be used.

> Not all operations are written in this code. Only the ones which are included in instructions of NES are implemented.
