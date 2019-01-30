<!DOCTYPE html>
<html>
  <head>
   <meta charset="utf-8">
   <meta name="viewport" content="width=device-width, initial-scale=1">
   <title>{{.Title}}</title>
{{ if not .Raw }}
   <link rel="stylesheet" media="screen, projection" href="/static/paste.css">
   <link rel="stylesheet" media="screen, projection" href="/static/prism.css">
{{ end }}
  </head>
  <body>
{{ if not .Raw }}
    <script src="/static/prism.js"></script>
    <div class="title">
      Title: {{.Title}}
      Language: {{.Lang}}
    </div>
    <div class="date">
      Date: {{.Date}}
    </div>

<pre class="language-{{.Lang}} line-numbers"><code>
{{ else }}
<pre><code>
{{end}}
      {{printf "%s" .Content}}
</code></pre>
  </body>
</html>
