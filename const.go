package grpc

const (
	space = " "

	httpStatusHeader  = "http-code"
	headerContentType = "Content-Type"
	grpcHeaderValue   = "application/grpc"

	grpcStatusHeader      = "Grpc-Metadata-X-Http-Code"
	grpcMetadataFormatter = "Grpc-Metadata-%s"
	grpcGatewayUri        = "grpcgateway-uri"
	grpcGatewayMethod     = "grpcgateway-method"
	grpcGatewayProto      = "grpcgateway-proto"
)
