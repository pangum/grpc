package grpc

const (
	space = " "
	tcp   = "tcp"
	slash = "/"

	httpStatusHeader  = "http-code"
	headerContentType = "Content-Type"
	grpcHeaderValue   = "application/grpc"
	rawHeaderValue    = "application/raw"

	grpcStatusHeader      = "Grpc-Metadata-X-Http-Code"
	grpcMetadataFormatter = "Grpc-Metadata-%s"
	grpcGatewayUri        = "grpcgateway-uri"
	grpcGatewayMethod     = "grpcgateway-method"
	grpcGatewayProto      = "grpcgateway-proto"
)
