# kiteroot
kiteroot is a easy-to-use html parser for Go.

The project is inspired by [BeautifulSoup](http://www.crummy.com/software/BeautifulSoup/) written for Python.

To find tags after parse, you can use find functions: Find, FindAll, FindWithAttrs, FindAllWithAttrs.

    var r io.Reader
    ...
    
    doc, _ := kiteroot.Parse(r) 
    
    // you can get urls of all links in the document you parse like this:
    links := doc.FindAll("a")
    for _, link := links {
      fmt.Println(link.Attribute("href"))
    }
    
    // you can also find a tag containing specific attributes.
    title := doc.FindWithAttrs("span", kiteroot.MakeAttrs("class", "title")) // And it is same as: doc.Find("span", "class", "title")
    
    
    
