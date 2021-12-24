package response

type Response interface {
	ToContentByte() []byte //能转换成 []byte

	ToErrorByte() []byte
}
