# netgo
A tool like netcat but may be better.

## Functions

- Listen local port
- Connect to ip:port
- Transform command to remote machine
- "Fake" port forwarding

## How to install it

Install go first.

After you install go and config it well:
```
go get github.com/grt1st/netgo
go build github.com/grt1st/netgo
./netgo
```

## How to use it

1. Basic mode:
- Server mode: `./netgo -l -p 6666`
- Client mode: `./netgo -a localhost -p 6666`

You can use the modes just the way you use in `netcat`.

2. You can use it to transform file also:

```
./netgo -l -p 6666 > 2.txt
cat 1.txt | ./netgo -a localhost -p 6666
```

![](http://view.grt1st.cn/img/netgo0.png)

3. Also use it to transform shell:

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

4. What's more, try it:

```
./netgo -a www.baidu.com -p 80 -html
```
![](http://view.grt1st.cn/img/netgo4.png)

5. Port forwarding

The first: Listen two port at the same time, and two ports are connected:

```
./netgo -a localhost -p 6666 -p 6667
```

Use it when you don't have public ip address:

```
./netgo -a xx.xx.xx.xx -p 6666 -e /bin/bash
./netgo -a xx.xx.xx.xx -p 6667
```

The second: Forward your local port to the remote machine which has a public ip address.

```
./netgo -a localhost -p 6666 -p 6667
```

And bind you local port forward by:

```
./netgo -a localhost -p 80 -rhost xx.xx.xx.xx:6666
```

Then others can visit your port by:

```
./netgo -a xx.xx.xx.xx -p 6667
```

> Note: The two port forwarding way has a disadvantage: the port can only get 1 client, so I say it's fake :)

## To do list

- Rebuild Code
- Specification error message
- Limit client number more effective
- Encrypted message

## How to contribute

Fork it and push.
