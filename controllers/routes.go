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

		//Genre routes
		v1.POST("/genre", s.CreateGenre)
		v1.GET("/genres", s.GetGenres)
		v1.GET("/genre/:id", s.GetGenre)
		v1.PUT("/genre/:id", s.UpdateGenre)
		v1.DELETE("/genre/:id", s.DeleteGenre)

		//CastMember routes
		v1.POST("/cast_member", s.CreateCastMember)
		v1.GET("/cast_members", s.GetCastMembers)
		v1.GET("/cast_member/:id", s.GetCastMember)
		v1.PUT("/cast_member/:id", s.UpdateCastMember)
		v1.DELETE("/cast_member/:id", s.DeleteCastMember)
	}
}
