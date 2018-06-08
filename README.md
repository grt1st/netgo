# netgo
A tool like netcat but may be better.

Still in development...

## How to use it

Install go first.

After you install go and config it well:
```
go build github.com/grt1st/netgo
./netgo
```

## Function

Basic mode:
- Server mode: `./netgo -l -p 6666`
- Client mode: `./netgo -a localhost -p 6666`

You can use the modes just the way you use in `netcat`.

You can use it to transform file also:

```
./netgo -l -p 6666 > 2.txt
cat 1.txt | ./netgo -a localhost -p 6666
```

![](http://view.grt1st.cn/img/netgo0.png)

Also use it to transform shell:

Forward shell：
```
./netgo -l -p 6666
./netgo -a localhost -p 6666 -e /bin/bash
```
![](http://view.grt1st.cn/img/netgo1.png)
Reverse shell：
```
./netgo -l -p 6666 -e /bin/bash
./netgo -a localhost -p 6666
```
![](http://view.grt1st.cn/img/netgo2.png)

Another way:

```
bash -c "bash -i &>/dev/tcp/localhost/6666 0>&1"
./netgo -l -p 6666
```
![](http://view.grt1st.cn/img/netgo3.png)



## To do list

- err info
- limit client number
- add more functions

