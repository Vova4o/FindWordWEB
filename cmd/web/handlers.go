package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
)

func (app *application) ShowHome(c *gin.Context) {
	app.render(c, "home.page.gohtml", nil)
}

func (app *application) ShowPage(c *gin.Context) {
	page := c.Param("page")
	app.render(c, page+".page.gohtml", nil)
}

type Filter struct {
	CurPage      int    `json:"currentPage"`
	Letters      string `json:"filter"`
	WordsPerPage int    `json:"wordsPerPage"`
}

func (app *application) FilterdList(c *gin.Context) {
	var Filter Filter
	// get from json request filter
	if err := c.BindJSON(&Filter); err != nil {
		log.Println("Error building JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "need a list of letters"})
		return
	}

	filter := Filter.Letters
	if filter == "" {
		c.JSON(http.StatusOK, gin.H{
			"words": showPerPage(app.cfg.Nouns, Filter.CurPage, Filter.WordsPerPage),
			"total": len(app.cfg.Nouns),
		})
		return
	}

	words := app.cfg.Nouns
	SortedList := make([]string, len(words)) // Replace Type with the actual type of elements in Nouns
	copy(SortedList, words)
	// list {1,2,3,4,5}
	for _, l := range filter {
		if unicode.IsLetter(l) {
			// log.Println("Letter:", string(l))
			l = unicode.ToLower(l)
			if l > -'а' && l <= 'я' {
				SortedList = listFilter(l, SortedList)
			}
		} else if unicode.IsNumber(l) {
			// log.Println("Number:", string(l))
			SortedList = listFilterByLen(l, SortedList)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"words": showPerPage(SortedList, Filter.CurPage, Filter.WordsPerPage),
		"total": len(SortedList),
	})
}

func listFilter(l rune, words []string) []string {
	var SortedList []string
	for _, word := range words {
		// check for leter and number and the split the case.
		if strings.Contains(word, string(l)) {
			SortedList = append(SortedList, word)
		}
	}
	return SortedList
}

func listFilterByLen(l rune, words []string) []string {
	var SortedList []string
	lnum, err := strconv.Atoi(string(l))
	if err != nil {
		return SortedList
	}
	for _, word := range words {
		if utf8.RuneCountInString(word) <= lnum {
			SortedList = append(SortedList, word)
		}
	}

	// log.Printf("Number of words found: %d", len(SortedList))
	return SortedList
}

func showPerPage(sort []string, currentPage int, wordsPerPage int) []string {
	if currentPage < 1 {
		currentPage = 1
	}
	// log.Println("Current page:", len(sort), currentPage, wordsPerPage)
	start := (currentPage - 1) * wordsPerPage
	end := currentPage * wordsPerPage
	if start > len(sort) || end < 0 {
		return []string{}
	}
	if end > len(sort) {
		end = len(sort)
	}
	return sort[start:end]
}
