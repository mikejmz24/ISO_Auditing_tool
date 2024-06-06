// package api
//
// //TODO: Create logic to handler HTML templ responses
// //TODO: Refactor code to avoid repeatig code on both controllers
//
// import (
//
//	"ISO_Auditing_Tool/cmd/api/repositories"
//	"ISO_Auditing_Tool/cmd/web"
//	// "net/http"
//
//	"github.com/a-h/templ"
//	"github.com/gin-gonic/gin"
//
// )
//
//	type ClauseController struct {
//		Repo repositories.Repository
//	}
//
//	func NewClauseController(repo repositories.Repository) *ClauseController {
//		return &ClauseController{
//			Repo: repo,
//		}
//	}
//
//	func (cc *ClauseController) GetAllClauses(c *gin.Context) {
//		clauses, err := cc.Repo.GetAllClauses()
//		if err != nil {
//			templ.Handler(web.Clauses(clauses)).ServeHTTP(c.Writer, c.Request)
//		}
//	}
package html

import (
	"ISO_Auditing_Tool/cmd/api/repositories"
	"ISO_Auditing_Tool/cmd/web"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HtmlClauseController struct {
	Repo repositories.Repository
}

func NewHtmlClauseController(repo repositories.Repository) *HtmlClauseController {
	return &HtmlClauseController{
		Repo: repo,
	}
}

func (cc *HtmlClauseController) GetAllClauses(c *gin.Context) {
	clauses, err := cc.Repo.GetAllClauses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	templ.Handler(web.Clauses(clauses)).ServeHTTP(c.Writer, c.Request)
}
