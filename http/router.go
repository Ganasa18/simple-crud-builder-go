package server

func (h *HttpServe) setupRouter() {
	v1 := h.router.Group("/api/v1")

	// Model module
	v1.POST("/models", h.modelDynamicController.CreateModel)

	// Table module
	v1.POST("/tables/:table", h.tableDynamicController.CreateRecord)

}
