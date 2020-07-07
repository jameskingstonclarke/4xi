package src


import "github.com/golang/snappy"

func Compress(src []byte) (encoded []byte){
	encoded = snappy.Encode(nil, src)
	return
}

func Uncompress(src []byte) (decoded []byte){
	decoded, err := snappy.Decode(nil, src)
	if err != nil{
		CLogErr(err)
	}
	return
}