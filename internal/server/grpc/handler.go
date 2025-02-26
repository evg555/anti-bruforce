package internalgrpc

type Handler struct {
	app    Application
	logger Logger

	//pb.UnsafeEventServiceServer
}

//func renderErrorResponse(err error) *pb.Response {
//	return &pb.Response{
//		Error:   true,
//		Message: err.Error(),
//	}
//}
