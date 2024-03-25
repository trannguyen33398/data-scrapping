package scrapdata

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/pdfcrowd/pdfcrowd-go"
	"golang.org/x/net/html"
)

func (r *controller) Scrap(c *gin.Context, url string) error {
    
	res, err := http.Get(url)
	if err != nil {
		r.logger.Fatal(err,"Failed to connect to the target page")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		r.logger.Fatalf(err,"HTTP Error %d: %s", res.StatusCode, res.Status)
	}

	// parse the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		r.logger.Fatal(err,"Failed to parse the HTML document")
	}

	// scraping logic...
	//find html with table element and has class wikitable
	tabletHTMLElement := doc.Find("table.wikitable")
	var f func(*html.Node, int, int, map[string]int, map[int][]string)

	//function to handle node selection
	f = func(n *html.Node, idx int, headerIdx int, keyMap map[string]int, valueMap map[int][]string) {

		if n.Type == html.TextNode {
			// Keep newlines and spaces, like jQuery
			trimmedData := strings.TrimSpace(n.Data)

			if len(trimmedData) > 0 {

				/*get the colum name and save as a map with value is the index of it*/
				if idx == 0 {
					keyMap[trimmedData] = headerIdx
				} 
				/* find text node contain only digit or decimal*/
				matched, _ := regexp.MatchString(`^[0-9]+(\.[0-9]+)?$`, trimmedData)
				/* if match condition store all the value along with the header index in a map*/
				if matched {	
					val, _ := valueMap[headerIdx]
					valueMap[headerIdx] = append(val, trimmedData)
				}

			}

		}
		if n.FirstChild != nil {
            /* loop through the node element and it's sibling to handle the data inside, need to use recursive technical to access it's sibling */
			for c := n.FirstChild; c != nil; c = c.NextSibling {

				f(c, idx, headerIdx, keyMap, valueMap)
				headerIdx += 1
			}
		}
	}
	var wg sync.WaitGroup
	wg.Add(tabletHTMLElement.Length())
	client := pdfcrowd.NewHtmlToImageClient("demo", "ce544b6ea52a5621fb9d55f8b542d14d")

	tabletHTMLElement.Each(func(i int, s *goquery.Selection) {
		go func (){
			defer wg.Done()
			
			header := s.Find("tr")

			keyMap := make(map[string]int)
			valueMap := make(map[int][]string)
			
			for idx, n := range header.Nodes {
				headerIdx := 0
				
				f(n, idx, headerIdx, keyMap, valueMap)
	
			}
			
			bar := charts.NewLine()
	
			// set  global option to enable save as image
	
			bar.SetGlobalOptions(charts.WithToolboxOpts(opts.Toolbox{
				Show: true,
				Feature: &opts.ToolBoxFeature{
					SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
						Show:  true,
						Type:  "png",
						Title: fmt.Sprint("table", i, ".html"),
					},
				},
			},
			))
			bar.SetXAxis([]string{}).SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
			invalidTable := 0
	
			/* Mapping the column name and it's data */
			for header, pos := range keyMap {
				data, _ := valueMap[pos]
				items := make([]opts.LineData, 0)
				if len(data) > 0 {
					for _, d := range data {
						items = append(items, opts.LineData{Value: d})
					}
					bar.AddSeries(header, items)
				} else {
					invalidTable++
				}
	
			}
	
		
			// Create html file
			if invalidTable < len(keyMap) {
				html, _ := os.Create(fmt.Sprint("table", i, ".html"))
				bar.Render(html)
	
				// Create image from html string using pdfcrowd-go 
				var buffer bytes.Buffer
	
				// Render the chart to the buffer
				bar.Render(&buffer)
	
				// configure the conversion
				client.SetOutputFormat("png")
			
				// run the conversion and write the result to a file
				err := client.ConvertStringToFile(buffer.String(), fmt.Sprint("image", i, ".png"))
				if err != nil {
					why, ok := err.(pdfcrowd.Error)
					if ok {
						os.Stderr.WriteString(fmt.Sprintf("Pdfcrowd Error: %s\n", why))
					} else {
						os.Stderr.WriteString(fmt.Sprintf("Generic Error: %s\n", err))
					}
	
					panic(err.Error())
				}
			}
		}()
		

	})
	wg.Wait()

	return nil
}
