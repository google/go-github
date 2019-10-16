[![GoDoc](https://godoc.org/github.com/google/go-github/scrape?status.svg)](https://godoc.org/github.com/google/go-github/scrape)

The scrape package provides an experimental client for accessing additional
GitHub data via screen scraping.  It is designed to be a client of last resort
for data that cannot be retrieved via the REST or GraphQL APIs.

# What should be added here

**Add only what you need.**  Whereas the main go-github library attempts to
implement the entire GitHub REST API, there is little point in trying to do that
here.  Certainly, feel free to contribution patches to get data you actually
need, but I'd rather not try and provide exhaustive coverage of all GitHub data
here.

**Add only what can't be accessed elsewhere.**  If the data can be retrieved
through the REST or GraphQL API, use the appropriate libraries for that.

**Prefer read-only access.**  For now, I'm only focusing on reading data. It
might be that writing works fine as well, but it is of course much riskier.

# How to add methods

See [apps.go](apps.go) for examples of methods that access data.  Basically,
fetch the contents of the page using `client.get`, and then use goquery to dig
into the markup on the page.  Prefer selectors that grab semantic ID or class
names, as they are more likely to be stable.
