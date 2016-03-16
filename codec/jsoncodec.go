package codec

import(
	"io"
	"github.com/funny/link"
	"github.com/funny/binary"
)

type JsonIo struct{
}

func (this JsonIo) Read(r *binary.Reader) []byte{
	b := r.ReadUint32LE()
	if b == 0 || b > (65535) {
		return nil
	}
	buf := make([]byte, b+4+4)//包长度+id +body
	binary.PutUint32LE(buf,b)
	r.ReadFull(buf[4:])

	return buf
}
func (this JsonIo) Write(w *binary.Writer,buf []byte){
	w.WriteBytes(buf)
}

func GetJsonIoCodec() JsonCodecType{
	return NewJsonCodec(JsonIo{})
}

func NewJsonCodec(Spliter binary.Spliter) JsonCodecType{
	return JsonCodecType{Spliter}
}

type JsonCodecType struct {
	Spliter binary.Spliter  //接口  Read Write
}

func (this JsonCodecType) NewEncoder(w io.Writer) link.Encoder{
	return jsonEncoder{this.Spliter,binary.NewWriter(w)}
}

func (this JsonCodecType) NewDecoder(r io.Reader) link.Decoder{
	return jsonDecoder{this.Spliter,binary.NewReader(r)}
}

type jsonEncoder struct {
	Spliter binary.Spliter  //实现write 接口
	Writer  *binary.Writer
}
func (this jsonEncoder) Encode(msg interface{}) error{
	this.Writer.WritePacket(msg.([]byte),this.Spliter)
	return this.Writer.Flush()
}

type jsonDecoder struct {
	Spliter binary.Spliter //实现read 接口
	Reader *binary.Reader
}

func (this jsonDecoder) Decode(msg interface{}) error{
	msg=this.Reader.ReadPacket(this.Spliter)
	return this.Reader.Error()
}

