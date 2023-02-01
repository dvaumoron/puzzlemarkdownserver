# puzzlemarkdownserver

A simple implementation of a [puzzlemarkdownservice](https://github.com/dvaumoron/puzzlemarkdownservice) server.

Use [goldmark](https://github.com/yuin/goldmark) (wich is [CommonMark](https://spec.commonmark.org) compliant) as the markdown implementation, with a custom extension for wikilink (targeting a  [WebComponent](https://www.webcomponents.org/)) :
- "[[ pageName ]]" became \<wiki-link title="pageName">pageName</wiki-link>
- "[[ pageName | linkName ]]" became \<wiki-link title="pageName">linkName</wiki-link>
- "[[ langTag/pageName ]]" became \<wiki-link lang="langTag" title="pageName">pageName</wiki-link>
- "[[ path/to/wiki#pageName ]]" became \<wiki-link wiki="path/to/wiki" title="pageName">pageName</wiki-link>
- "[[ path/to/wiki#langTag/pageName ]]" became \<wiki-link wiki="path/to/wiki" lang="langTag" title="pageName">pageName</wiki-link>

And so on...