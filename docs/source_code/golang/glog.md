# [Glog](https://github.com/golang/glog)

## Index

- [const usage](#a-usage-of-const)
- [atomic value](#atomic-value-operation)
- [convert string](#convert-string-to-certain-type)

---

#### A usage of const

```go
const (
	infoLog int32 = iota
	warningLog
	errorLog
	fatalLog
	numSeverity = 4
)
```

---

#### Atomic value operation

```go
atomic.LoadInt32(a)
atomic.StoreInt32(a,storeVal)
```

> Support `int32`,`int64`,`uint32`,`uint64`,`uintptr`,`unsafe.Pointer`

---

#### Convert string to certain type

```go
result := strconv.FormatInt(int64(a),10)
a,err := strconv.ParseInt(str,10,64)
```

> Support `FormatInt`,`FormatUint`,`FormatBool`,`FormatFloat`
>
> In contrast, `ParseInt`,`ParseUint`,`ParseBool`,`ParseFloat`

---

#### Query for time

```go
year,month,day := time.Now().Date()
hour,minute,second := time.Now().Clock()
```







