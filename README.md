# yenc
--
    import "github.com/andrewstuart/yenc"

Package yenc implements readers writers for the YENC encoding format.

## Usage

```go
var (
	ErrBadCRC    = fmt.Errorf("CRC check error")
	ErrWrongSize = fmt.Errorf("size check error")
)
```
Error constants

#### type Reader

```go
type Reader struct {
	Length          int
	CRC             hash.Hash32
	Headers, Footer *YENCHeader
}
```

Reader implements the io.Reader methods for an underlying YENC document/stream.
It additionally exposes some of the metadata that may be useful for consumers.

#### func  NewReader

```go
func NewReader(r io.Reader) *Reader
```
NewReader returns a reader from an input reader.

#### func (*Reader) Read

```go
func (d *Reader) Read(p []byte) (bytesRead int, err error)
```

#### type Writer

```go
type Writer struct {
	CRC            hash.Hash32
	Length, Line   int
	Name           string
	Header, Footer *YENCHeader
}
```


#### func  NewWriter

```go
func NewWriter(w io.Writer) *Writer
```

#### func (*Writer) Close

```go
func (w *Writer) Close() error
```

#### func (*Writer) Write

```go
func (w *Writer) Write(p []byte) (written int, err error)
```

#### type YENCHeader

```go
type YENCHeader map[string]string
```


#### func  ReadYENCHeader

```go
func ReadYENCHeader(br *bufio.Reader) (*YENCHeader, error)
```

#### func (*YENCHeader) Add

```go
func (y *YENCHeader) Add(k, v string)
```

#### func (*YENCHeader) Get

```go
func (y *YENCHeader) Get(k string) string
```

#### func (*YENCHeader) String

```go
func (y *YENCHeader) String() string
```
