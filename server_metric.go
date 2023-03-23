package grpc

func (s *Server) metric(register register) (err error) {
	if !s.config.metricEnabled() {
		return
	}

	if handler, me := register.Metric(); nil != me {
		err = me
	} else {
		s.mux.Handle(s.config.Metric.Pattern, handler)
	}

	return
}
