package main

import (
	"fmt"
	"github.com/DavidBelicza/TextRank"
	"github.com/DavidBelicza/TextRank/parse"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type RequestQuery struct {
	Query string `json:"query"`
}

type Response struct {
	Summary  string `json:"Summary"`
	Keywords string `json:"Keywords"`
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func GetSummaryKeywords(rawText string, numSentences int, numKeywords int) (string, string) {
	// TextRank object
	tr := textrank.NewTextRank()
	// Default Rule for parsing.
	rule := textrank.NewDefaultRule()
	// Default Language for filtering stop words.
	language := textrank.NewDefaultLanguage()
	// Default algorithm for ranking text.
	algorithmDef := textrank.NewDefaultAlgorithm()

	// Add text.
	tr.Populate(rawText, language, rule)
	// Run the ranking.
	tr.Ranking(algorithmDef)

	// Get all phrases order by weight.
	// rankedPhrases := textrank.FindPhrases(tr)
	// Most important phrase.
	// fmt.Println(rankedPhrases[0])

	var topKeywords string
	// Get all words order by weight.
	words := textrank.FindSingleWords(tr)
	for _, k := range words[:numKeywords] {
		topKeywords += strings.Title(k.Word) + ", "
	}
	// Most important keywords
	fmt.Println(words[:numKeywords])

	// Get the most important X sentences. Importance by phrase weights.
	sentences := textrank.FindSentencesByRelationWeight(tr, numSentences)

	var topSentences string
	for _, s := range sentences {
		// sentences[i].Value = strings.TrimSpace(s.Value)
		topSentences += strings.TrimSpace(s.Value) + " "
		// s.Value = strings.TrimSpace(s.Value)
		// fmt.Println(s.Value)
	}

	// Found sentences
	fmt.Println(topSentences)

	return topSentences, topKeywords

	// // // Get the most important 10 sentences. Importance by word occurrence.
	// sentences = textrank.FindSentencesByWordQtyWeight(tr, 3)
	// // Found sentences
	// fmt.Println(sentences)

	// // Get the first 10 sentences, start from 5th sentence.
	// sentences = textrank.FindSentencesFrom(tr, 5, 10)
	// // Found sentences
	// fmt.Println(sentences)

	// // Get sentences by phrase/word chains order by position in text.
	// sentencesPh := textrank.FindSentencesByPhraseChain(tr, []string{"gnome", "shell", "extension"})
	// // Found sentence.
	// fmt.Println(sentencesPh[0])
}

func main() {
	router := gin.Default()
	router.Use(CORSMiddleware())
	// Example for binding JSON ({"user": "manu", "password": "123"})
	router.POST("/summarize", func(c *gin.Context) {
		var json RequestQuery
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		rule := textrank.NewDefaultRule()
		text := parse.TokenizeText(json.Query, rule)

		fmt.Println(len(text.GetSentences()))
		if len(text.GetSentences()) < 3 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Not enough sentences to generate a summary"})
			return
		}

		if len(text.GetSentences()) >= 3 {
			sentences, keywords := GetSummaryKeywords(json.Query, 3, 5)
			c.JSON(http.StatusOK, gin.H{"results": Response{Summary: sentences, Keywords: keywords}})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// router.POST("/post", func(c *gin.Context) {

	// 	id := c.Query("id")
	// 	page := c.DefaultQuery("page", "0")
	// 	name := c.PostForm("name")
	// 	message := c.PostForm("message")

	// 	fmt.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
	// })

	router.Run(":8080")
}
