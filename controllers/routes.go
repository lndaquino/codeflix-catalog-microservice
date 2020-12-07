package controllers

func (s *Server) initializeRoutes() {
	v1 := s.Router.Group("/api/v1")
	{
		//Category routes
		v1.POST("/category", s.CreateCategory)
		v1.GET("/categories", s.GetCategories)
		v1.GET("/category/:id", s.GetCategory)
		v1.PUT("/category/:id", s.UpdateCategory)
		v1.DELETE("/category/:id", s.DeleteCategory)
	}
}
