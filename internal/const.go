package internal

const (
	Space = " "
	Tcp   = "tcp"
	Slash = "/"
	Comma = ","

	HttpStatusHeader  = "http-code"
	HeaderContentType = "Content-Type"
	GrpcHeaderValue   = "application/grpc"
	RawHeaderValue    = "application/raw"

	GrpcStatusHeader      = "Grpc-Metadata-X-Http-Code"
	GrpcMetadataFormatter = "Grpc-Metadata-%s"
	GrpcGatewayUri        = "grpcgateway-uri"
	GrpcGatewayMethod     = "grpcgateway-method"
	GrpcGatewayProto      = "grpcgateway-proto"

	HeaderXRealIp       = "X-REAL-IP"
	HeaderXForwardedFor = "X-FORWARDED-FOR"
	HeaderAllowOrigin   = "Access-Control-Allow-Origin"
	HeaderAllowMethods  = "Access-Control-Allow-Methods"
	HeaderAllowHeaders  = "Access-Control-Allow-Headers"

	MethodOptions = "OPTIONS"
)
