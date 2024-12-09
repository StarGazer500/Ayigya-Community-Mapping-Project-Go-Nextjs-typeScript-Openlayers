package routers

import (
	"github.com/StarGazer500/ayigya/controllers"

	"github.com/gin-gonic/gin"
)

func MapRoutes(route *gin.RouterGroup) {
	route.GET("/map-display", controllers.MapPageDisplay)
	route.GET("/featurelayers", controllers.FeatureLayers)
	route.POST("/featureattributes", controllers.FeatreAttributes)
	route.POST("/featureoperatures", controllers.SelectOperator)
	route.POST("/makeqquery", controllers.MakeQuery)
	route.POST("/searchallfeaturelayersdata", controllers.SearchAllFeaturesData)
	route.POST("/searchbyfeaturelayer", controllers.SearchByFeatureLayer)
	route.POST("/searchbycolumn", controllers.SearchByColumn)
	route.POST("/simplesearch", controllers.SimpleSearch)
	
}
