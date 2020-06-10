package page

var Markdown = `
<!doctype html>
<html>
<head>
  <meta charset="utf-8"/>
  <title>Markdown</title>
  <style>
    table { padding: 0px; word-break: initial; }
    table { border-collapse: collapse; border-spacing: 0px; width: 100%%; overflow: auto; break-inside: auto; text-align: left; }
    table tr { border-top: 1px solid rgb(223, 226, 229); margin: 0px; padding: 0px; }
    table tr:nth-child(2n), thead { background-color: rgb(248, 248, 248); }
    table tr th { font-weight: bold; border-width: 1px 1px 0px; border-top-style: solid; border-right-style: solid; border-left-style: solid; border-top-color: rgb(223, 226, 229); border-right-color: rgb(223, 226, 229); border-left-color: rgb(223, 226, 229); border-image: initial; border-bottom-style: initial; border-bottom-color: initial; margin: 0px; padding: 6px 13px; }
    table tr td { border: 1px solid rgb(223, 226, 229); margin: 0px; padding: 6px 13px; }
    table tr th:first-child, table tr td:first-child { margin-top: 0px; }
    table tr th:last-child, table tr td:last-child { margin-bottom: 0px; }
  </style>
</head>
<body>
  <div id="content"></div>
  <script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></script>
  <script>
    document.getElementById('content').innerHTML =
      marked('%s');
  </script>
</body>
</html>
`
