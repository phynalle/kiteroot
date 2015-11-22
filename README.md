# kiteroot
kiteroot is a convenient tool to play with html tag for Go.

The project is inspired by [BeautifulSoup](http://www.crummy.com/software/BeautifulSoup/) written for Python.

To select tags after parse, you can use these functions: Find, FindAll, FindWithAttrs, FindAllWithAttrs.

    var r io.Reader
    ...
    
    doc, _ := kiteroot.Parse(r) 
    
    // you can get urls of all links in the document you parse like this:
    links := doc.FindAll("a")
    for _, link := links {
      fmt.Println(link.Attribute("href"))
    }
    
    // you can also find a tag containing specific attributes.
    title := doc.FindWithAttrs("span", kiteroot.MakeAttrs("class", "title"))
    // And it is same with: title := doc.Find("span", "class", "title")
    
    
    
