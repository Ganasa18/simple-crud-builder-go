package server

func (h *HttpServe) setupRouter() {
	v1 := h.router.Group("/api/v1")

	// Model module
	v1.POST("/models", h.modelDynamicController.CreateModel)
	v1.GET("/models", h.modelDynamicController.ListModel)

	// Table module
	v1.POST("/tables/:table", h.tableDynamicController.CreateRecord)
	v1.GET("/tables/:table", h.tableDynamicController.ListRecord)

}
