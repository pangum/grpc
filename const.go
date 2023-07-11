package grpc

const (
	space = " "
	tcp   = "tcp"
	slash = "/"
	comma = ","

	httpStatusHeader  = "http-code"
	headerContentType = "Content-Type"
	grpcHeaderValue   = "application/grpc"
	rawHeaderValue    = "application/raw"

	grpcStatusHeader      = "Grpc-Metadata-X-Http-Code"
	grpcMetadataFormatter = "Grpc-Metadata-%s"
	grpcGatewayUri        = "grpcgateway-uri"
	grpcGatewayMethod     = "grpcgateway-method"
	grpcGatewayProto      = "grpcgateway-proto"

	headerXRealIp       = "X-REAL-IP"
	headerXForwardedFor = "X-FORWARDED-FOR"
	headerAllowOrigin   = "Access-Control-Allow-Origin"
	headerAllowMethods  = "Access-Control-Allow-Methods"
	headerAllowHeaders  = "Access-Control-Allow-Headers"

	methodOptions = "OPTIONS"
)
