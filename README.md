# Shukuangkuang
## Description
This is a tool can display the usage of per CPU in Linux just like the TaskManager in Windows

## Preview
<img width="1701" alt="image" src="https://github.com/Yeuoly/shukuangkuang/assets/45712896/80423d18-c7d5-4309-9d20-6af7fb451a58">


## Usage
You can compile it yourself or just use my [release](https://github.com/Yeuoly/shukuangkuang/releases)

### LogicalCPUMode (default)
You can use it in logicalCpuMode, which will display all the usage of logical cpus
```bash
shukuangkuang
```

### SingleCPUMode
You can also display the average usage off all cpus
```bash
shukuangkuang -mode=false
```

## Compile
```bash
./build/linux_amd64.sh
```

Or you can just compile it through your own environment
```bash
go build cmd/main/main.go
```

## Contributions
All contributions including issue/pull request/discussion is welcome

## License
MIT
