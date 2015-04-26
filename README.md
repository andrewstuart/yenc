# yenc
--
    import "git.astuart.co/andrew/yenc"

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
	Headers, Footer textproto.MIMEHeader
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
func (d *Reader) Read(p []byte) (int, error)
```
