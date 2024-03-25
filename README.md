To start the project:

- Run go run main.go

For testing:

- Run go test -v -cover ./pkg/test/

The idea to start this project:

- I go to the wiki page and check the work to see what it response: a html string
- I check where is the table html code and I find the common thing of those table: class wikitable. I know that I need to hand a strong query selection to work on this so I decided to go with golang and "github.com/PuerkitoBio/goquery"
- I noted that some page has multiple table and on each table there can be more than one colum are digit or decimal, so there must be a loop stuff to workaround to get data of that. I go deeply the document, source code of github.com/PuerkitoBio/goquery to understand how the data are store ( html.Node and sibling)
- Now I just need to map the data correct with it's colum name then use a lib to help me generate as a picture
- "github.com/go-echarts/go-echarts/v2/charts" help me to generate a html file and option to download ( not automatically on the server side) the chart as img. Then I also use "github.com/pdfcrowd/pdfcrowd-go" to save the chart as picture.

