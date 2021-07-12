package page

import (
	"encoding/json"
	"fmt"
	"strings"
)

var jsTempl = `
<!doctype html>
<html>

<head>
    <meta charset="utf-8" />
    <title>%s</title>

    <style>
    pre {outline: 1px solid #ccc; padding: 5px; margin: 5px; }
    .string { color: green; }
    .number { color: darkorange; }
    .boolean { color: blue; }
    .null { color: magenta; }
    .key { color: red; }
    </style>
</head>

<body>
    <pre id="content"></pre>
    <script>
        function output(inp) {
            document.getElementById('content').innerHTML = inp;
        }

        function syntaxHighlight(json) {
            json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
            return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function(match) {
                var cls = 'number';
                if (/^"/.test(match)) {
                    if (/:$/.test(match)) {
                        cls = 'key';
                    } else {
                        cls = 'string';
                    }
                } else if (/true|false/.test(match)) {
                    cls = 'boolean';
                } else if (/null/.test(match)) {
                    cls = 'null';
                }
                return '<span class="' + cls + '">' + match + '</span>';
            });
        }

        var str = '%s';
        var str = JSON.stringify(JSON.parse(str), undefined, 4);
        output(syntaxHighlight(str));
    </script>
</body>
</html>
`

func JSON(title string, obj interface{}) ([]byte, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	if title == "" {
		title = "json"
	}

	formated := string(bytes)
	formated = strings.Replace(formated, "\r", "\\r", -1)
	formated = strings.Replace(formated, "\n", "\\n", -1)
	return []byte(fmt.Sprintf(jsTempl, title, formated)), nil
}
